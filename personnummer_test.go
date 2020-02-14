package personnummer

import (
	"fmt"
	"testing"
	"time"

	"github.com/frozzare/go-assert"
)

var invalidNumbers = []interface{}{
	nil,
	[]string{},
	[]int{},
	true,
	false,
	0,
	"19112233-4455",
	"20112233-4455",
	"9999999999",
	"199999999999",
	"199909193776",
	"Just a string",
}

func TestPersonnummerWithControlDigit(t *testing.T) {
	assert.True(t, Valid("8507099805"))
	assert.True(t, Valid("198507099805"))
	assert.True(t, Valid("198507099813"))
	assert.True(t, Valid("850709-9813"))
	assert.True(t, Valid("196411139808"))
}

func TestPersonnummerWithoutControlDigit(t *testing.T) {
	assert.False(t, Valid("19850709980"))
	assert.False(t, Valid("19850709981"))
	assert.False(t, Valid("19641113980"))
}

func TestPersonnummerWithWrongPersonnummerOrTypes(t *testing.T) {
	for _, n := range invalidNumbers {
		assert.False(t, Valid(fmt.Sprintf("%v", n)))
	}
}

func TestCoordinationNumbers(t *testing.T) {
	assert.True(t, Valid("198507699802"))
	assert.True(t, Valid("850769-9802"))
	assert.True(t, Valid("198507699810"))
	assert.True(t, Valid("850769-9810"))
	assert.True(t, Valid("198507699802"))
	assert.True(t, Valid("198507699810"))
}

func TestWrongCoOrdinationNumbers(t *testing.T) {
	assert.False(t, Valid("198567099805"))
}

func TestCoordinationNumbersCheck(t *testing.T) {
	p, _ := Parse("198507699802")
	assert.True(t, p.IsCoordinationNumber())
}

func TestShouldParsePersonnummer(t *testing.T) {
	p, _ := Parse("198507699802")
	assert.Equal(t, "34", p.Age)
	assert.Equal(t, "19", p.Century)
	assert.Equal(t, "1985", p.FullYear)
	assert.Equal(t, "85", p.Year)
	assert.Equal(t, "07", p.Month)
	assert.Equal(t, "69", p.Day)
	assert.Equal(t, "-", p.Sep)
	assert.Equal(t, "980", p.Num)
	assert.Equal(t, "2", p.Check)
}

func TestShouldThrowErrorForBadInputsWhenParsing(t *testing.T) {
	for _, n := range invalidNumbers {
		_, e := Parse(fmt.Sprintf("%v", n))
		assert.NotNil(t, e)
	}
}

func TestShouldFormatInputValusAsPersonnummer(t *testing.T) {
	shortFormat := map[string]string{
		"850709-9805": "19850709-9805",
		"850709-9813": "198507099813",
	}

	for expected, input := range shortFormat {
		p, _ := Parse(input)
		v, _ := p.Format()
		assert.Equal(t, expected, v)
	}

	longFormat := map[string]string{
		"198507099805": "19850709-9805",
		"198507099813": "198507099813",
	}

	for expected, input := range longFormat {
		p, _ := Parse(input)
		v, _ := p.Format(true)
		assert.Equal(t, expected, v)
	}
}

/*
func TestShouldNotFormatInputValueAsPersonnummer(t *testing.T) {
	for _, n := range invalidNumbers {
		_, err := Format(n)
		assert.NotNil(t, err)
	}
}
*/
func TestGetAge(t *testing.T) {
	now = func() time.Time {
		return time.Date(2019, 7, 13, 0, 0, 0, 0, time.UTC)
	}

	p, _ := Parse("198507099805")
	assert.Equal(t, "34", p.Age)

	p, _ = Parse("198507099813")
	assert.Equal(t, "34", p.Age)

	p, _ = Parse("196411139808")
	assert.Equal(t, "54", p.Age)

	p, _ = Parse("19121212+1212")
	assert.Equal(t, "106", p.Age)
}

func TestGetAgeWithCoOrdinationNumbers(t *testing.T) {
	now = func() time.Time {
		return time.Date(2019, 7, 13, 0, 0, 0, 0, time.UTC)
	}

	p, _ := Parse("198507699810")
	assert.Equal(t, "34", p.Age)

	p, _ = Parse("198507699802")
	assert.Equal(t, "34", p.Age)
}

func TestSex(t *testing.T) {
	p, _ := Parse("8507099813")
	v, _ := p.IsMale()
	assert.True(t, v)

	p, _ = Parse("198507099813")
	v, _ = p.IsFemale()
	assert.False(t, v)

	p, _ = Parse("198507099805")
	v, _ = p.IsFemale()
	assert.True(t, v)

	p, _ = Parse("198507099805")
	v, _ = p.IsMale()
	assert.False(t, v)
}

func TestSexWithCoOrdinationNumbers(t *testing.T) {
	p, _ := Parse("198507099813")
	v, _ := p.IsMale()
	assert.True(t, v)

	p, _ = Parse("198507099813")
	v, _ = p.IsFemale()
	assert.False(t, v)

	p, _ = Parse("198507699802")
	v, _ = p.IsFemale()
	assert.True(t, v)

	p, _ = Parse("198507699802")
	v, _ = p.IsMale()
	assert.False(t, v)
}

func BenchmarkValid(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Valid("198507099805")
	}
}
