package personnummer

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/frozzare/go-assert"
	"github.com/frozzare/go/http2"
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

var testList []*TestListItem

func TestMain(m *testing.M) {
	if err := http2.GetJSON("https://raw.githubusercontent.com/personnummer/meta/master/testdata/list.json", &testList); err != nil {
		log.Fatal(err)
	}

	code := m.Run()
	os.Exit(code)
}

func TestPersonnummerList(t *testing.T) {
	for _, item := range testList {
		for _, format := range availableListFormats {
			assert.Equal(t, item.Valid, Valid(item.Get(format)))
		}
	}
}

func TestPersonnummerFormat(t *testing.T) {
	for _, item := range testList {
		if !item.Valid {
			continue
		}

		for _, format := range availableListFormats {
			if format == "short_format" && strings.Contains(item.SeparatedFormat, "+") {
				continue
			}

			p, _ := New(item.Get(format))
			v1, _ := p.Format()
			assert.Equal(t, item.SeparatedFormat, v1)

			v2, _ := p.Format(true)
			assert.Equal(t, item.LongFormat, v2)
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
			assert.NotNil(t, err)
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
			assert.Equal(t, item.IsMale, p.IsMale())
			assert.Equal(t, item.IsFemale, p.IsFemale())
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
			assert.Equal(t, true, p.IsCoordinationNumber())
		}

		tt, _ := time.Parse("2006-01-02", fmt.Sprintf("%s-%s-%s", year, month, day))
		a := math.Floor(float64(now().Sub(tt)/1e6) / 3.15576e+10)

		for _, format := range availableListFormats {
			if format == "short_format" && strings.Contains(item.SeparatedFormat, "+") {
				continue
			}

			p, _ := Parse(item.Get(format))
			assert.Equal(t, a, p.GetAge())
		}
	}
}

func BenchmarkValid(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Valid("198507099805")
	}
}
