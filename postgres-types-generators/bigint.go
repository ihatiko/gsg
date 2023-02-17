package postgres_types_generators

import "github.com/brianvoe/gofakeit/v6"

func BigIntGenerator() any {
	return gofakeit.IntRange(-1000, 1000)
}
