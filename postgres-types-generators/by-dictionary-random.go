package postgres_types_generators

import "math/rand"

func ByDictionaryStringRandom(data []string) (any, string) {
	result := data[rand.Intn(len(data))]
	return result, result
}
