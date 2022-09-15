package hebcal

// Hebcal - A Jewish Calendar Generator
// Copyright (c) 2022 Michael J. Radwin
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
	"strconv"
	"strings"

	"github.com/hebcal/hebcal-go/hdate"
	"github.com/hebcal/hebcal-go/locales"
)

type omerEvent struct {
	Date            hdate.HDate
	OmerDay         int
	WeekNumber      int
	DaysWithinWeeks int
}

func newOmerEvent(hd hdate.HDate, omerDay int) omerEvent {
	if omerDay < 1 || omerDay > 49 {
		panic("invalid omerDay")
	}
	week := ((omerDay - 1) / 7) + 1
	days := (omerDay % 7)
	if days == 0 {
		days = 7
	}
	return omerEvent{Date: hd, OmerDay: omerDay, WeekNumber: week, DaysWithinWeeks: days}
}

func (ev omerEvent) GetDate() hdate.HDate {
	return ev.Date
}

func (ev omerEvent) Render(locale string) string {
	dayOfTheOmer, _ := locales.LookupTranslation("day of the Omer", locale)
	return strconv.Itoa(ev.OmerDay) + " " + dayOfTheOmer
}

func (ev omerEvent) GetFlags() HolidayFlags {
	return OMER_COUNT
}

func (ev omerEvent) GetEmoji() string {
	number := ev.OmerDay
	var r rune
	if number <= 20 {
		r = rune(9312 + number - 1)
	} else if number <= 35 {
		// between 21 and 35 inclusive
		r = rune(12881 + number - 21)
	} else {
		// between 36 and 49 inclusive
		r = rune(12977 + number - 36)
	}
	return string(r)
}

func (ev omerEvent) Basename() string {
	return ev.Render("en")
}

func (ev omerEvent) GetWeeks() int {
	if ev.DaysWithinWeeks == 7 {
		return ev.WeekNumber
	} else {
		return ev.WeekNumber - 1
	}
}

// adapted from pip hdate package (GPL)
// https://github.com/py-libhdate/py-libhdate/blob/master/hdate/date.py

var tens = []string{"", "עֲשָׂרָה", "עֶשְׂרִים", "שְׁלוֹשִׁים", "אַרְבָּעִים"}
var ones = []string{
	"",
	"אֶחָד",
	"שְׁנַיִם",
	"שְׁלוֹשָׁה",
	"אַרְבָּעָה",
	"חֲמִשָׁה",
	"שִׁשָׁה",
	"שִׁבְעָה",
	"שְׁמוֹנָה",
	"תִּשְׁעָה",
}

const shnei = "שְׁנֵי"
const yamim = "יָמִים"
const shneiYamim = shnei + " " + yamim
const shavuot = "שָׁבוּעוֹת"
const yom = "יוֹם"

var yomEchad = yom + " " + ones[1]

func todayIsHe(omer int) string {
	var ten = (omer / 10)
	var one = omer % 10
	var str = "הַיוֹם "
	if 10 < omer && omer < 20 {
		str += ones[one] + " עָשָׂר"
	} else if omer > 9 {
		str += ones[one]
		if one != 0 {
			str += " וְ"
		}
	}
	if omer > 2 {
		if (omer > 20) || (omer == 10) || (omer == 20) {
			str += tens[ten]
		}
		if omer < 11 {
			str += ones[one] + " " + yamim + " "
		} else {
			str += " " + yom + " "
		}
	} else if omer == 1 {
		str += yomEchad + " "
	} else { // omer == 2
		str += shneiYamim + " "
	}
	if omer > 6 {
		str = strings.TrimSpace(str) // remove trailing space before comma
		str += ", שְׁהֵם "
		var weeks = (omer / 7)
		var days = omer % 7
		if weeks > 2 {
			str += ones[weeks] + " " + shavuot + " "
		} else if weeks == 1 {
			str += "שָׁבוּעַ" + " " + ones[1] + " "
		} else { // weeks == 2
			str += shnei + " " + shavuot + " "
		}
		if days != 0 {
			str += "וְ"
			if days > 2 {
				str += ones[days] + " " + yamim + " "
			} else if days == 1 {
				str += yomEchad + " "
			} else { // days == 2
				str += shneiYamim + " "
			}
		}
	}
	str += "לָעוֹמֶר"
	return str
}

func (ev omerEvent) TodayIs(locale string) string {
	if locale == "he" {
		return todayIsHe(ev.OmerDay)
	}
	totalDaysStr := "days"
	if ev.OmerDay == 1 {
		totalDaysStr = "day"
	}
	str := "Today is " + strconv.Itoa(ev.OmerDay) + " " + totalDaysStr
	if ev.WeekNumber > 1 || ev.OmerDay == 7 {
		day7 := ev.DaysWithinWeeks == 7
		numWeeks := ev.WeekNumber - 1
		if day7 {
			numWeeks = ev.WeekNumber
		}
		weeksStr := "weeks"
		if numWeeks == 1 {
			weeksStr = "week"
		}
		str += ", which is " + strconv.Itoa(numWeeks) + " " + weeksStr
		if !day7 {
			dayStr := "days"
			if ev.DaysWithinWeeks == 1 {
				dayStr = "day"
			}
			str += " and " + strconv.Itoa(ev.DaysWithinWeeks) + " " + dayStr
		}
	}
	return str + " of the Omer"
}

var sefirot = []string{
	"",
	"Lovingkindness",
	"Might",
	"Beauty",
	"Eternity",
	"Splendor",
	"Foundation",
	"Majesty",
}

var sefirotTranslit = []string{
	"",
	"Chesed",
	"Gevurah",
	"Tiferet",
	"Netzach",
	"Hod",
	"Yesod",
	"Malkhut",
}

func (ev omerEvent) Sefira(locale string) string {
	weekStr := sefirot[ev.WeekNumber]
	dayWithinWeekStr := sefirot[ev.DaysWithinWeeks]
	weekNum2or6 := ev.WeekNumber == 2 || ev.WeekNumber == 6
	if locale == "he" {
		week, _ := locales.LookupTranslation(weekStr, locale)
		dayWithinWeek, _ := locales.LookupTranslation(dayWithinWeekStr, locale)
		prefix := "שֶׁבְּ"
		if weekNum2or6 {
			prefix = "שֶׁבִּ"
		}
		return dayWithinWeek + " " + prefix + week
	} else if locale == "translit" {
		week := sefirotTranslit[ev.WeekNumber]
		dayWithinWeek := sefirotTranslit[ev.DaysWithinWeeks]
		prefix := "sheb'"
		if weekNum2or6 {
			prefix = "shebi"
		}
		return dayWithinWeek + " " + prefix + week
	} else {
		return dayWithinWeekStr + " within " + weekStr
	}
}
