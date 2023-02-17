package postgres_types_generators

import "math/rand"

var boolState = []any{
	true,
	false,
}

func BoolGenerator() any {
	return boolState[rand.Intn(len(boolState))]
}
