package postgres_types_generators

import "github.com/brianvoe/gofakeit/v6"

func DateGenerator() (any, string) {
	//TODO date-range-rule
	data := gofakeit.Date()
	return data, data.String()
}
