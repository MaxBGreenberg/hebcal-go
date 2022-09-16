package dafyomi

import (
	"fmt"
	"testing"
	"time"

	"github.com/hebcal/hebcal-go/hdate"
	"github.com/stretchr/testify/assert"
)

func TestMishnaYomi(t *testing.T) {
	idx := MakeMishnaYomiIndex()
	hd := hdate.FromGregorian(1995, time.December, 17)
	mishna, err := idx.Lookup(hd)
	assert.Equal(t, nil, err)
	assert.Equal(t, mishna, MishnaPair{
		Mishna{Tractate: "Bava Kamma", Chap: 5, Verse: 7},
		Mishna{Tractate: "Bava Kamma", Chap: 6, Verse: 1},
	})
	mishna, _ = idx.Lookup(hdate.FromGregorian(2024, time.April, 5))
	assert.Equal(t, mishna, MishnaPair{
		Mishna{Tractate: "Nedarim", Chap: 11, Verse: 12},
		Mishna{Tractate: "Nazir", Chap: 1, Verse: 1},
	})
	assert.Equal(t, "Nedarim 11:12-Nazir 1:1", mishna.String())
}

func ExampleMishnaYomiIndex_Lookup() {
	idx := MakeMishnaYomiIndex()
	mishna, _ := idx.Lookup(hdate.FromGregorian(1947, time.May, 20))
	fmt.Println(mishna)
	mishna, _ = idx.Lookup(hdate.FromGregorian(1995, time.December, 17))
	fmt.Println(mishna)
	mishna, _ = idx.Lookup(hdate.FromGregorian(2022, time.August, 1))
	fmt.Println(mishna)
	// Output:
	// Berakhot 1:1-2
	// Bava Kamma 5:7-6:1
	// Terumot 11:3-4
}
