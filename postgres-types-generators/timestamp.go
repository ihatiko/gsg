package postgres_types_generators

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-module/carbon"
	"time"
)

func prevYear(t time.Time) time.Time {
	y, m, _ := t.Date()
	return time.Date(y-1, m, 1, 0, 0, 0, 0, time.UTC)
}

func TimeStampGenerator() (any, string) {
	//TODO range generator
	data := carbon.FromStdTime(gofakeit.DateRange(prevYear(time.Now()), time.Now())).ToStdTime()
	return data, data.String()
}
