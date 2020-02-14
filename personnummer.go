package personnummer

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

const (
	lengthWithoutCentury = 10
	lengthWithCentury    = 12
)

var (
	errInvalidSecurityNumber = errors.New("Invalid swedish social security number")
	monthDays                = map[int]int{
		1:  31,
		3:  31,
		4:  30,
		5:  31,
		6:  30,
		7:  31,
		8:  31,
		9:  30,
		10: 31,
		11: 30,
		12: 31,
	}
	now   = time.Now
	rule3 = [...]int{0, 2, 4, 6, 8, 1, 3, 5, 7, 9}
)

// charsToDigit converts char bytes to a digit
// example: ['1', '1'] => 11
func charsToDigit(chars []byte) int {
	l := len(chars)
	r := 0
	for i, c := range chars {
		p := int((c - '0'))
		for j := 0; j < l-i-1; j++ {
			p *= 10
		}
		r += p
	}
	return r
}

// getCleanNumber will return clean numbers.
func getCleanNumber(in string) []byte {
	cleanNumber := make([]byte, 0, len(in))

	for _, c := range in {
		if c == '+' {
			continue
		}
		if c == '-' {
			continue
		}

		if c > '9' {
			return nil
		}
		if c < '0' {
			return nil
		}

		cleanNumber = append(cleanNumber, byte(c))
	}

	return cleanNumber
}

// getCoOrdinationDay will return co-ordination day.
func getCoOrdinationDay(day []byte) []byte {
	d := charsToDigit(day)
	if d < 60 {
		return day
	}

	d -= 60

	if d < 10 {
		return []byte{'0', byte(d) + '0'}
	}

	return []byte{
		byte(d)/10 + '0',
		byte(d)%10 + '0',
	}
}

// luhn will test if the given string is a valid luhn string.
func luhn(s []byte) bool {
	odd := len(s) & 1

	var sum int

	for i, c := range s {
		if i&1 == odd {
			sum += rule3[c-'0']
		} else {
			sum += int(c - '0')
		}
	}

	return sum%10 == 0
}

// toString converts int to string.
func toString(in interface{}) string {
	switch v := in.(type) {
	case int, int32, int64, uint, uint32, uint64:
		return fmt.Sprint(v)
	case string:
		return v
	default:
		return ""
	}
}

// input time without centry.
func validateTime(time []byte) bool {
	length := len(time)

	date := charsToDigit(time[length-2 : length])
	month := charsToDigit(time[length-4 : length-2])

	if month != 2 {
		days, ok := monthDays[month]
		if !ok {
			return false
		}
		return date <= days
	}

	year := charsToDigit(time[:length-4])

	leapYear := year%4 == 0 && year%100 != 0 || year%400 == 0

	if leapYear {
		return date <= 29
	}
	return date <= 28
}

// Personnummer represents the personnummer struct.
type Personnummer struct {
	Age                string
	Century            string
	FullYear           string
	Year               string
	Month              string
	Day                string
	Sep                string
	Num                string
	Check              string
	leapYear           bool
	coordinationNumber bool
}

// Options represents the personnummer options.
type Options struct {
}

// New parse a Swedish social security numbers and returns a new struct or a error.
func New(ssn string, options ...*Options) (*Personnummer, error) {
	p := &Personnummer{}

	if err := p.parse(ssn); err != nil {
		return nil, err
	}

	return p, nil
}

// parse Swedish social security numbers and set struct properpties or return a error.
func (p *Personnummer) parse(ssn string) error {
	var num, check string

	if ssn == "" {
		return errInvalidSecurityNumber
	}

	dateBytes := getCleanNumber(ssn)

	if len(dateBytes) == 0 || len(dateBytes) < 8 {
		return errInvalidSecurityNumber
	}

	plus := strings.Contains(ssn, "+")

	switch len(dateBytes) {
	case lengthWithCentury:
		num = string(dateBytes[8:11])
		check = string(dateBytes[11:])
		dateBytes = dateBytes[2:8]
		break
	case lengthWithoutCentury:
		num = string(dateBytes[6:9])
		check = string(dateBytes[9:])
		dateBytes = dateBytes[0:6]
		break
	}

	length := len(dateBytes)
	day := charsToDigit(dateBytes[length-2 : length])
	month := charsToDigit(dateBytes[length-4 : length-2])

	if month != 2 {
		if _, ok := monthDays[month]; !ok {
			return errInvalidSecurityNumber
		}
	}

	p.Check = check
	p.Num = num
	p.Sep = "-"

	if day < 10 {
		p.Day = fmt.Sprintf("0%d", day)
	} else {
		p.Day = toString(day)
	}

	if month < 10 {
		p.Month = fmt.Sprintf("0%d", month)
	} else {
		p.Month = toString(month)
	}

	year := 0
	fullYear := 0

	if len(dateBytes[:length-4]) == 4 {
		fullYear = charsToDigit(dateBytes[:length-4])
		p.Century = string(dateBytes[:length-4][0:2])
		year = charsToDigit(dateBytes[:length-4][2:])
	} else {
		year = charsToDigit(dateBytes[:length-4])

		var baseYear int
		if plus {
			baseYear = now().Year() - 100
			p.Sep = "+"
		} else {
			baseYear = now().Year()
		}

		centuryStr := strconv.Itoa((baseYear - ((baseYear - year) % 100)))
		century, err := strconv.Atoi(centuryStr[0:2])
		if err != nil {
			return err
		}

		p.Century = toString(century)
		fullYear = charsToDigit(append([]byte(centuryStr[0:2]), dateBytes[:length-4]...))
	}

	p.leapYear = year%4 == 0 && year%100 != 0 || year%400 == 0
	p.Year = toString(year)
	p.FullYear = toString(fullYear)

	ageDay := day
	if ageDay >= 61 && ageDay < 91 {
		p.coordinationNumber = true
		ageDay = ageDay - 60
	}

	t := time.Date(fullYear, time.Month(month), ageDay, 0, 0, 0, 0, time.UTC)
	age := math.Floor(float64(now().Sub(t).Milliseconds()) / 3.15576e+10)
	p.Age = fmt.Sprintf("%.0f", age)

	if !p.valid() {
		return errInvalidSecurityNumber
	}

	return nil
}

// Valid will validate Swedish social security numbers.
func (p *Personnummer) valid() bool {
	ssn := fmt.Sprintf("%s%s%s%s%s%s", p.Century, p.Year, p.Month, p.Day, p.Num, p.Check)

	bytes := []byte(ssn)
	if !luhn(bytes[2:]) {
		return false
	}

	var dateBytes = append(bytes[:6], getCoOrdinationDay(bytes[6:8])...)

	return validateTime(dateBytes)
}

// Format a Swedish social security number as one of the official formats,
// a long format or a short format.
func (p *Personnummer) Format(longFormat ...bool) (string, error) {
	if len(longFormat) > 0 && longFormat[0] {
		return fmt.Sprintf("%s%s%s%s%s%s", p.Century, p.Year, p.Month, p.Day, p.Num, p.Check), nil
	}

	return fmt.Sprintf("%s%s%s%s%s%s", p.Year, p.Month, p.Day, p.Sep, p.Num, p.Check), nil
}

// Check if a Swedish social security number is a coordination number or not.
// Returns true if it's a coordination number.
func (p *Personnummer) IsCoordinationNumber() bool {
	return p.coordinationNumber
}

// IsFemale checks if a Swedish social security number is for a female.
func (p *Personnummer) IsFemale() (bool, error) {
	male, err := p.IsMale()

	if err != nil {
		return false, err
	}

	return !male, err
}

// IsMale checks if a Swedish social security number is for a male.
// The second argument should be a boolean
func (p *Personnummer) IsMale() (bool, error) {
	sexDigit := int(p.Num[2])

	return sexDigit%2 == 1, nil
}

// Valid will validate Swedish social security numbers
func Valid(ssn string, options ...*Options) bool {
	_, err := Parse(ssn, options...)
	return err == nil
}

// Parse Swedish social security numbers and return a new struct.
func Parse(ssn string, options ...*Options) (*Personnummer, error) {
	return New(ssn, options...)
}
