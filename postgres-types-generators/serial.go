package postgres_types_generators

import "strconv"

var sequenceMap = map[string]int{}

func Serial(key string) (any, string) {
	data, ok := sequenceMap[key]
	if ok {
		data++
		sequenceMap[key] = data
		return data, strconv.Itoa(data)
	}
	sequenceMap[key] = 1
	return 1, strconv.Itoa(data)
}
