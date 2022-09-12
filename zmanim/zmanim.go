package zmanim

// Hebcal - A Jewish Calendar Generator
// Copyright (c) 2022 Michael J. Radwin
// Derived from original JavaScript version, Copyright (C) 2014 Eyal Schachter
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

import (
	"time"

	"github.com/nathan-osman/go-sunrise"
)

// Tzais (nightfall) based on the opinion of the Geonim calculated at
// the sun's position at 8.5° below the western horizon.
// https://kosherjava.com/zmanim/docs/api/com/kosherjava/zmanim/ComplexZmanimCalendar.html#getTzaisGeonim8Point5Degrees()
const Tzeit3SmallStars = 8.5

// Tzais (nightfall) based on the opinion of the
// Geonim calculated as 30 minutes after sunset during the equinox
// (on March 16, about 4 days before the astronomical equinox, the day that
// a solar hour is 60 minutes) in Yerushalayim.
// https://kosherjava.com/zmanim/docs/api/com/kosherjava/zmanim/ComplexZmanimCalendar.html#getTzaisGeonim7Point083Degrees()
const Tzeit3MediumStars = 7.083

// Zmanim are used to calculate halachic times
type Zmanim struct {
	Latitude  float64    // In the ragnge [-90,90]
	Longitude float64    // In the range [-180,180]
	Year      int        // Gregorian year
	Month     time.Month // Gregorian month
	Day       int        // Gregorian day
	loc       *time.Location
}

// New makes an instance used for calculating various halachic times during this day.
//
// tzid should be a timezone identifier such as "America/Los_Angeles" or "Asia/Jerusalem".
//
// This function panics if the latitude or longitude are out of range, or if
// the timezone cannot be loaded.
func New(latitude, longitude float64, date time.Time, tzid string) Zmanim {
	if latitude < -90 || latitude > 90 {
		panic("Latitude out of range [-90,90]")
	}
	if longitude < -180 || longitude > 180 {
		panic("Longitude out of range [-180,180]")
	}
	year, month, day := date.Date()
	loc, err := time.LoadLocation(tzid)
	if err != nil {
		panic(err)
	}
	return Zmanim{Latitude: latitude, Longitude: longitude, Year: year, Month: month, Day: day, loc: loc}
}

var nilTime = time.Time{}

func (z *Zmanim) inLoc(dt time.Time) time.Time {
	if dt == nilTime {
		return dt
	}
	return dt.In(z.loc)
}

// Sunset ("shkiah") calculates when the sun will set on the given day
// at the specified location.
//
// Sunset is defined as when the upper edge of the Sun disappears below
// the horizon (0.833° below horizon)
//
// Returns time.Time{} if there sun does not rise or set
func (z *Zmanim) Sunset() time.Time {
	_, set := sunrise.SunriseSunset(z.Latitude, z.Longitude, z.Year, z.Month, z.Day)
	return z.inLoc(set)
}

// Sunrise ("neitz haChama") is defined as when the upper edge of the
// Sun appears over the eastern horizon in the morning
// (0.833° above horizon).
func (z *Zmanim) Sunrise() time.Time {
	rise, _ := sunrise.SunriseSunset(z.Latitude, z.Longitude, z.Year, z.Month, z.Day)
	return z.inLoc(rise)
}

func (z *Zmanim) timeAtAngle(angle float64, rising bool) time.Time {
	morning, evening := sunrise.TimeOfElevation(z.Latitude, z.Longitude, -angle, z.Year, z.Month, z.Day)
	if rising {
		return z.inLoc(morning)
	} else {
		return z.inLoc(evening)
	}
}

// Civil dawn; Sun is 6° below the horizon in the morning
func (z *Zmanim) Dawn() time.Time {
	return z.timeAtAngle(6.0, true)
}

// Civil dusk; Sun is 6° below the horizon in the evening
func (z *Zmanim) Dusk() time.Time {
	return z.timeAtAngle(6.0, false)
}

// ms in hour
func (z *Zmanim) hour() int {
	rise, set := sunrise.SunriseSunset(z.Latitude, z.Longitude, z.Year, z.Month, z.Day)
	millis := set.UnixMilli() - rise.UnixMilli()
	return int(millis / 12)
}

// hour in ms / (1000 ms in s * 60 s in m) = mins in halachic hour
func (z *Zmanim) hourMins() int {
	return z.hour() / (1000 * 60)
}

func (z *Zmanim) GregEve() time.Time {
	prev := time.Date(z.Year, z.Month, z.Day-1, 0, 0, 0, 0, z.loc)
	year, month, day := prev.Date()
	zman := Zmanim{
		Latitude:  z.Latitude,
		Longitude: z.Longitude,
		Year:      year,
		Month:     month,
		Day:       day,
		loc:       z.loc,
	}
	return zman.Sunset()
}

// ms in hour
func (z *Zmanim) nightHour() int {
	set := z.GregEve()
	rise := z.Sunrise()
	millis := rise.UnixMilli() - set.UnixMilli()
	return int(millis / 12)
}

// hour in ms / (1000 ms in s * 60 s in m) = mins in halachic hour
func (z *Zmanim) nightHourMins() int {
	return z.nightHour() / (1000 * 60)
}

// sunrise plus N halachic hours
func (z *Zmanim) hourOffset(hours float64) time.Time {
	rise := z.Sunrise()
	millis := rise.UnixMilli() + int64(float64(z.hour())*hours)
	return time.UnixMilli(millis).In(z.loc)
}

// Midday – Chatzot; Sunrise plus 6 halachic hours
func (z *Zmanim) Chatzot() time.Time {
	return z.hourOffset(6)
}

// Midnight – Chatzot; Sunset plus 6 halachic hours
func (z *Zmanim) ChatzotNight() time.Time {
	rise := z.Sunrise()
	millis := rise.UnixMilli() - int64(z.nightHour()*6)
	return time.UnixMilli(millis).In(z.loc)
}

// Dawn – Alot haShachar; Sun is 16.1° below the horizon in the morning
func (z *Zmanim) AlotHaShachar() time.Time {
	return z.timeAtAngle(16.1, true)
}

// Earliest talis & tefillin – Misheyakir; Sun is 11.5° below the horizon in the morning
func (z *Zmanim) Misheyakir() time.Time {
	return z.timeAtAngle(11.5, true)
}

// Earliest talis & tefillin – Misheyakir Machmir; Sun is 10.2° below the horizon in the morning
func (z *Zmanim) MisheyakirMachmir() time.Time {
	return z.timeAtAngle(10.2, true)
}

// Latest Shema (Gra); Sunrise plus 3 halachic hours, according to the Gra
func (z *Zmanim) SofZmanShma() time.Time {
	return z.hourOffset(3)
}

// Latest Shacharit (Gra); Sunrise plus 4 halachic hours, according to the Gra
func (z *Zmanim) SofZmanTfilla() time.Time {
	return z.hourOffset(4)
}

func (z *Zmanim) sofZmanMGA(hours int) time.Time {
	alot72 := z.SunriseOffset(-72, false)
	tzeit72 := z.SunsetOffset(72, false)
	alot72ms := alot72.UnixMilli()
	temporalHour := (tzeit72.UnixMilli() - alot72ms) / 12 // ms in hour
	millis := alot72ms + (int64(hours) * temporalHour)
	return time.UnixMilli(millis).In(z.loc)
}

// Latest Shema (MGA); Sunrise plus 3 halachic hours, according to Magen Avraham
func (z *Zmanim) SofZmanShmaMGA() time.Time {
	return z.sofZmanMGA(3)
}

// Latest Shacharit (MGA); Sunrise plus 4 halachic hours, according to Magen Avraham
func (z *Zmanim) SofZmanTfillaMGA() time.Time {
	return z.sofZmanMGA(4)
}

// Earliest Mincha – Mincha Gedola; Sunrise plus 6.5 halachic hours
func (z *Zmanim) MinchaGedola() time.Time {
	return z.hourOffset(6.5)
}

// Preferable earliest time to recite Minchah – Mincha Ketana; Sunrise plus 9.5 halachic hours
func (z *Zmanim) MinchaKetana() time.Time {
	return z.hourOffset(9.5)
}

// Plag haMincha; Sunrise plus 10.75 halachic hours
func (z *Zmanim) PlagHaMincha() time.Time {
	return z.hourOffset(10.75)
}

// Tzeit is defined as nightfall, when 3 stars are observable in the night sky with the naked eye.
//
// For 3 small stars use 8.5°
//
// For 3 medium stars use 7.083°
func (z *Zmanim) Tzeit(angle float64) time.Time {
	if angle == 0 {
		angle = Tzeit3SmallStars
	}
	return z.timeAtAngle(angle, false)
}

func (z *Zmanim) riseSetOffset(t time.Time, offset int, roundTime bool) time.Time {
	if t == nilTime {
		return t
	}
	year, month, day := t.Date()
	hour, min, sec := t.Clock()
	if roundTime {
		// For positive offsets only, round up to next minute if needed
		if offset > 0 && sec >= 30 {
			offset++
		}
		sec = 0
	}
	return time.Date(year, month, day, hour, min+offset, sec, 0, z.loc)
}

// Returns sunrise + offset minutes (either positive or negative).
//
// If roundTime is true, rounds to the nearest minute (setting seconds to zero).
func (z *Zmanim) SunriseOffset(offset int, roundTime bool) time.Time {
	return z.riseSetOffset(z.Sunrise(), offset, roundTime)
}

// Returns sunset + offset minutes (either positive or negative).
//
// This function is used with a negative offset to calculate candle-lighting times,
// typically -18 minutes before sundown (or -40 in Jerusalem).
//
// This function can be used with a positive offset to calculate Tzeit (nightfall).
//
// For Havdalah according to Rabbeinu Tam, use 72, which approximates
// when 3 small stars are observable in the night sky with the naked eye.
// Other typical values include 50 minutes (3 small stars) or 42 minutes
// (3 medium stars).
//
// If roundTime is true, rounds to the nearest minute (setting seconds to zero).
func (z *Zmanim) SunsetOffset(offset int, roundTime bool) time.Time {
	return z.riseSetOffset(z.Sunset(), offset, roundTime)
}
