package personnummer

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"
)

type TestListItem struct {
	Integer         int    `json:"integer"`
	LongFormat      string `json:"long_format"`
	ShortFormat     string `json:"short_format"`
	SeparatedFormat string `json:"separated_format"`
	SeparatedLong   string `json:"separated_long"`
	Valid           bool   `json:"valid"`
	Type            string `json:"type"`
	IsMale          bool   `json:"isMale"`
	IsFemale        bool   `json:"isFemale"`
}

func (t *TestListItem) Get(key string) string {
	switch key {
	case "integer":
		return fmt.Sprintf("%d", t.Integer)
	case "long_format":
		return t.LongFormat
	case "short_format":
		return t.ShortFormat
	case "separated_format":
		return t.SeparatedFormat
	case "separated_long":
		return t.SeparatedLong
	default:
		break
	}
	return ""
}

var availableListFormats = []string{
	"integer",
	"long_format",
	"short_format",
	"separated_format",
	"separated_long",
}

func getJSON(url string, v any) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("GET %s: %s", url, res.Status)
	}
	return json.NewDecoder(res.Body).Decode(v)
}

func assertEqual[T comparable](t *testing.T, expected, actual T) {
	t.Helper()
	if expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

var testList []*TestListItem
var interimList []*TestListItem

func TestMain(m *testing.M) {
	if err := getJSON("https://raw.githubusercontent.com/personnummer/meta/HEAD/testdata/list.json", &testList); err != nil {
		log.Fatal(err)
	}

	if err := getJSON("https://raw.githubusercontent.com/personnummer/meta/HEAD/testdata/interim.json", &interimList); err != nil {
		log.Fatal(err)
	}

	code := m.Run()
	os.Exit(code)
}

func TestPersonnummerList(t *testing.T) {
	for _, item := range testList {
		for _, format := range availableListFormats {
			assertEqual(t, item.Valid, Valid(item.Get(format)))
		}
	}
}

func TestPersonnummerFormat(t *testing.T) {
	for _, item := range testList {
		if !item.Valid {
			continue
		}

		for _, format := range availableListFormats {
			if format == "short_format" {
				continue
			}

			p, _ := New(item.Get(format))
			v1, _ := p.Format()
			assertEqual(t, item.SeparatedFormat, v1)

			v2, _ := p.Format(true)
			assertEqual(t, item.LongFormat, v2)
		}
	}
}

func TestPersonnummerError(t *testing.T) {
	for _, item := range testList {
		if item.Valid {
			continue
		}

		for _, format := range availableListFormats {
			_, err := Parse(item.Get(format))
			if err == nil {
				t.Errorf("expected error for %s", item.Get(format))
			}
		}
	}
}

func TestPersonnummerSex(t *testing.T) {
	for _, item := range testList {
		if !item.Valid {
			continue
		}

		for _, format := range availableListFormats {
			p, _ := Parse(item.Get(format))
			assertEqual(t, item.IsMale, p.IsMale())
			assertEqual(t, item.IsFemale, p.IsFemale())
		}
	}
}

func TestPersonnummerDate(t *testing.T) {
	for _, item := range testList {
		if !item.Valid {
			continue
		}

		year := item.SeparatedLong[0:4]
		month := item.SeparatedLong[4:6]
		day := item.SeparatedLong[6:8]

		if item.Type == "con" {
			nDay, _ := strconv.Atoi(day)
			nDay = nDay - 60
			day = fmt.Sprintf("%02d", nDay)
			p, _ := Parse(item.SeparatedLong)
			assertEqual(t, true, p.IsCoordinationNumber())
		}

		tt, _ := time.Parse("2006-01-02", fmt.Sprintf("%s-%s-%s", year, month, day))

		for _, format := range availableListFormats {
			if format == "short_format" {
				continue
			}

			p, _ := Parse(item.Get(format))
			if got := p.GetDate(); !got.Equal(tt) {
				t.Errorf("expected %v, got %v", tt, got)
			}
		}
	}
}

func TestPersonnummerAge(t *testing.T) {
	for _, item := range testList {
		if !item.Valid {
			continue
		}

		year := item.SeparatedLong[0:4]
		month := item.SeparatedLong[4:6]
		day := item.SeparatedLong[6:8]

		if item.Type == "con" {
			nDay, _ := strconv.Atoi(day)
			nDay = nDay - 60
			day = fmt.Sprintf("%02d", nDay)
			p, _ := Parse(item.SeparatedLong)
			assertEqual(t, true, p.IsCoordinationNumber())
		}

		tt, _ := time.Parse("2006-01-02", fmt.Sprintf("%s-%s-%s", year, month, day))
		a := int(math.Floor(float64(now().Sub(tt)/1e6) / 3.15576e+10))

		for _, format := range availableListFormats {
			if format == "short_format" {
				continue
			}

			p, _ := Parse(item.Get(format))
			assertEqual(t, a, p.GetAge())
		}
	}
}

func TestInterimNumbers(t *testing.T) {
	for _, item := range interimList {
		if !item.Valid {
			continue
		}

		for _, format := range availableListFormats {
			if format == "integer" {
				continue
			}

			p, _ := New(item.Get(format), &Options{AllowInterimNumber: true})
			v1, _ := p.Format()
			assertEqual(t, item.SeparatedFormat, v1)

			v2, _ := p.Format(true)
			assertEqual(t, item.LongFormat, v2)
		}
	}
}

func TestInterimNumbersInvalid(t *testing.T) {
	for _, item := range interimList {
		if item.Valid {
			continue
		}

		for _, format := range availableListFormats {
			if format == "integer" {
				continue
			}

			_, err := New(item.Get(format), &Options{AllowInterimNumber: true})
			if err == nil {
				t.Errorf("expected error for %s", item.Get(format))
			}
		}
	}
}

func TestLuhn(t *testing.T) {
	assertEqual(t, true, luhn([]byte("1212121212")))
	assertEqual(t, false, luhn([]byte("12120111X3")))
}

func TestInvalidLengths(t *testing.T) {
	numbers := []string{"", "1", "12", "123", "1234", "12345", "123456", "12345678", "123456789", "1234567891", "12345678911", "123456789111", "1234567891111"}
	for _, n := range numbers {
		assertEqual(t, false, Valid(n))
	}
}

func BenchmarkValid(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Valid(testList[0].LongFormat)
	}
}
