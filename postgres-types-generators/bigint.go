package postgres_types_generators

import (
	"github.com/brianvoe/gofakeit/v6"
	"strconv"
)

func BigIntGenerator() (any, string) {
	data := gofakeit.IntRange(-1000, 1000)
	stringData := strconv.Itoa(data)
	return data, stringData
}
