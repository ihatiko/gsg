package postgres_types_generators

import (
	"math/rand"
	"strings"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
var STRINGUNIQUE = map[string]struct{}{}

func RandStringRunes(n int, unique bool) string {
	return getUniqueStringData(n, unique)
}

func getUniqueStringData(n int, unique bool) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	result := strings.ToValidUTF8(string(b), "")
	if unique {
		if _, ok := STRINGUNIQUE[result]; !ok {
			return getUniqueStringData(n, unique)
		}
		STRINGUNIQUE[result] = struct{}{}
	}
	return result
}
