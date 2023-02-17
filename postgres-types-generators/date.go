package postgres_types_generators

import "github.com/brianvoe/gofakeit/v6"

func DateGenerator() any {
	//TODO date-range-rule
	return gofakeit.Date()
}
