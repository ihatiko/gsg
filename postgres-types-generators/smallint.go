package postgres_types_generators

import (
	"github.com/brianvoe/gofakeit/v6"
	"strconv"
)

func SmallIntGenerator() (any, string) {
	data := gofakeit.IntRange(-32768, 32767)
	return data, strconv.Itoa(data)
}
