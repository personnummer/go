package personnummer

import (
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
	assert.True(t, Valid(8507099805))
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
		assert.False(t, Valid(n))
	}
}

func TestCoOrdinationNumbers(t *testing.T) {
	assert.True(t, Valid("198507699802"))
	assert.True(t, Valid("850769-9802"))
	assert.True(t, Valid("198507699810"))
	assert.True(t, Valid("850769-9810"))
	assert.False(t, Valid("198507699802", &Options{
		CoordinatioNumber: false,
	}))
	assert.False(t, Valid("198507699810", &Options{
		CoordinatioNumber: false,
	}))
}

func TestWrongCoOrdinationNumbers(t *testing.T) {
	assert.False(t, Valid("198567099805"))
}

func TestShouldFormatInputValusAsPersonnummer(t *testing.T) {
	shortFormat := map[string]interface{}{
		"850709-9805": "19850709-9805",
		"850709-9813": "198507099813",
	}

	for expected, input := range shortFormat {
		v, _ := Format(input)
		assert.Equal(t, expected, v)
	}

	longFormat := map[string]interface{}{
		"198507099805": "19850709-9805",
		"198507099813": "198507099813",
	}

	opt := &Options{
		LongFormat: true,
	}

	for expected, input := range longFormat {
		v, _ := Format(input, opt)
		assert.Equal(t, expected, v)
	}
}

func TestShouldNotFormatInputValueAsPersonnummer(t *testing.T) {
	for _, n := range invalidNumbers {
		_, err := Format(n)
		assert.NotNil(t, err)
	}
}

func TestGetAge(t *testing.T) {
	now = func() time.Time {
		return time.Date(2019, 7, 13, 01, 01, 01, 01, time.UTC)
	}

	age, _ := GetAge("198507099805")
	assert.Equal(t, 34, age)

	age, _ = GetAge("198507099813")
	assert.Equal(t, 34, age)

	age, _ = GetAge("196411139808")
	assert.Equal(t, 54, age)

	age, _ = GetAge("19121212+1212")
	assert.Equal(t, 106, age)
}

func TestGetAgeWithCoOrdinationNumbers(t *testing.T) {
	now = func() time.Time {
		return time.Date(2019, 7, 13, 01, 01, 01, 01, time.UTC)
	}

	age, _ := GetAge("198507699810")
	assert.Equal(t, 34, age)

	age, _ = GetAge("198507699802")
	assert.Equal(t, 34, age)
}

func TestGetAgeExcludeCoOrdinationNumbers(t *testing.T) {
	now = func() time.Time {
		return time.Date(2019, 7, 13, 01, 01, 01, 01, time.UTC)
	}

	age, err := GetAge("198507699810", &Options{
		CoordinatioNumber: false,
	})
	assert.Equal(t, 0, age)
	assert.NotNil(t, err)

	age, err = GetAge("198507699802", &Options{
		CoordinatioNumber: false,
	})
	assert.Equal(t, 0, age)
	assert.NotNil(t, err)
}

func TestSex(t *testing.T) {
	valid, _ := IsMale(8507099813, &Options{
		CoordinatioNumber: false,
	})
	assert.True(t, valid)

	valid, _ = IsFemale(198507099813, &Options{
		CoordinatioNumber: false,
	})
	assert.False(t, valid)

	valid, _ = IsFemale("198507099805", &Options{
		CoordinatioNumber: false,
	})
	assert.True(t, valid)

	valid, _ = IsMale("198507099805", &Options{
		CoordinatioNumber: false,
	})
	assert.False(t, valid)
}

func TestSexWithCoOrdinationNumbers(t *testing.T) {
	valid, _ := IsMale("198507099813")
	assert.True(t, valid)

	valid, _ = IsFemale("198507099813")
	assert.False(t, valid)

	valid, _ = IsFemale("198507699802")
	assert.True(t, valid)

	valid, _ = IsMale("198507699802")
	assert.False(t, valid)
}

func TestSexInvalidNumbers(t *testing.T) {
	for _, n := range invalidNumbers {
		_, err := IsMale(n)
		assert.NotNil(t, err)

		_, err = IsFemale(n)
		assert.NotNil(t, err)
	}
}

func BenchmarkValid(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Valid("198507099805")
	}
}

func BenchmarkFormat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Format("850709-9805")
	}
}
func BenchmarkGetAge(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetAge("198507099805")
	}
}
