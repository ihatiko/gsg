package postgres_types_generators

import (
	"math/rand"
	"strconv"
)

var boolState = []bool{
	true,
	false,
}

func BoolGenerator() (any, string) {
	data := boolState[rand.Intn(len(boolState))]
	return data, strconv.FormatBool(data)
}
