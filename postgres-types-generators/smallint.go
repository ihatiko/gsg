package postgres_types_generators

import "github.com/brianvoe/gofakeit/v6"

func SmallIntGenerator() any {
	return gofakeit.IntRange(-32768, 32767)
}
