package postgres_types_generators

import "math/rand"

func ByDictionaryStringRandom(data []string) any {
	return data[rand.Intn(len(data))]
}
