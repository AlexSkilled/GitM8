package model

import (
	"strings"
	"time"

	"github.com/araddon/dateparse"
)

type Datetime time.Time

func (d *Datetime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)
	t, err := dateparse.ParseAny(s)
	*d = Datetime(t)
	return
}

func (d Datetime) MarshalJSON() ([]byte, error) {
	t := time.Time(d)
	return t.MarshalJSON()
}
