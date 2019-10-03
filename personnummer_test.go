package personnummer

import (
	"testing"

	"github.com/frozzare/go-assert"
)

var invalidNumbers = []interface{}{
	[]string{},
	[]int{},
	true,
	false,
	1122334455,
	"112233-4455",
	"19112233-4455",
	"9999999999",
	"199999999999",
	"9913131315",
	"9911311232",
	"9902291237",
	"19990919_3766",
	"990919_3766",
	"199909193776",
	"Just a string",
	"990919+3776",
	"990919-3776",
	"9909193776",
}

func TestPersonnummerWithControlDigit(t *testing.T) {
	assert.True(t, Valid(6403273813))
	assert.True(t, Valid("510818-9167"))
	assert.True(t, Valid("19900101-0017"))
	assert.True(t, Valid("19130401+2931"))
	assert.True(t, Valid("196408233234"))
	assert.True(t, Valid("0001010107"))
	assert.True(t, Valid("000101-0107"))
	assert.True(t, Valid("200002296127"))
}

func TestPersonnummerWithoutControlDigit(t *testing.T) {
	assert.False(t, Valid(640327381))
	assert.False(t, Valid("510818-916"))
	assert.False(t, Valid("19900101-001"))
	assert.False(t, Valid("100101+001"))
}

func TestPersonnummerWithWrongPersonnummerOrTypes(t *testing.T) {
	for _, n := range invalidNumbers {
		assert.False(t, Valid(n))
	}
}

func TestLeapYear(t *testing.T) {
	assert.True(t, Valid("20000229-0005"))  // Divisible by 400
	assert.False(t, Valid("19000229-0005")) // Divisible by 100
	assert.True(t, Valid("20080229-0007"))  // Divisible by 4
	assert.False(t, Valid("20090229-0006")) // Not divisible by
}

func TestCoOrdinationNumbers(t *testing.T) {
	assert.True(t, Valid("701063-2391"))
	assert.True(t, Valid("640883-3231"))
	assert.False(t, Valid("701063-2391", &Options{
		CoordinatioNumber: false,
	}))
	assert.False(t, Valid("640883-3231", &Options{
		CoordinatioNumber: false,
	}))
}

func TestWrongCoOrdinationNumbers(t *testing.T) {
	assert.False(t, Valid("900161-0017"))
	assert.False(t, Valid("640893-3231"))
}

func TestSex(t *testing.T) {
	valid, err := IsMale(6403273813, &Options{
		CoordinatioNumber: false,
	})
	assert.True(t, valid)
	assert.Nil(t, err)

	valid, err = IsFemale(6403273813, &Options{
		CoordinatioNumber: false,
	})
	assert.False(t, valid)
	assert.Nil(t, err)

	valid, err = IsFemale("510818-9167", &Options{
		CoordinatioNumber: false,
	})
	assert.True(t, valid)
	assert.Nil(t, err)

	valid, err = IsMale("510818-9167", &Options{
		CoordinatioNumber: false,
	})
	assert.False(t, valid)
	assert.Nil(t, err)
}

func TestSexWithCoOrdinationNumbers(t *testing.T) {
	valid, err := IsMale("701063-2391")
	assert.True(t, valid)
	assert.Nil(t, err)

	valid, err = IsFemale("701063-2391")
	assert.False(t, valid)
	assert.Nil(t, err)

	valid, err = IsFemale("640883-3223")
	assert.True(t, valid)
	assert.Nil(t, err)

	valid, err = IsMale("640883-3223")
	assert.False(t, valid)
	assert.Nil(t, err)
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
		Valid("19900101-0017")
	}
}
