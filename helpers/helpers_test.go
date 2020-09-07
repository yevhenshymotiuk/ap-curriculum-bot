package helpers

import (
	"testing"
	"time"

	"github.com/matryer/is"
)

func TestFormatTime(t *testing.T) {
    is := is.New(t)

	l, err := LoadLocation()
	is.NoErr(err)

	d1 := time.Date(2020, 9, 4, 0, 0, 0, 0, l)
	is.Equal(FormatTime(&d1), "04.09")

	d2 := time.Date(2020, 9, 11, 0, 0, 0, 0, l)
	is.Equal(FormatTime(&d2), "11.09")
}
