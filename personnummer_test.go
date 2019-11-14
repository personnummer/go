package personnummer

import (
	"testing"

	"github.com/frozzare/go-assert"
)

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
	assert.False(t, Valid([]string{}))
	assert.False(t, Valid([]int{}))
	assert.False(t, Valid(true))
	assert.False(t, Valid(false))
	assert.False(t, Valid(0))
	assert.False(t, Valid("19112233-4455"))
	assert.False(t, Valid("20112233-4455"))
	assert.False(t, Valid("9999999999"))
	assert.False(t, Valid("199999999999"))
	assert.False(t, Valid("199909193776"))
	assert.False(t, Valid("Just a string"))
}

func TestCoOrdinationNumbers(t *testing.T) {
	assert.True(t, Valid("198507699802"))
	assert.True(t, Valid("850769-9802"))
	assert.True(t, Valid("198507699810"))
	assert.True(t, Valid("850769-9810"))
}

func TestWrongCoOrdinationNumbers(t *testing.T) {
	assert.False(t, Valid("198567099805"))
}

func BenchmarkValid(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Valid("198507099805")
	}
}

func BenchmarkValidString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ValidString("198507099805")
	}
}
