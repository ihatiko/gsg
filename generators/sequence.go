package generators

var sequenceMap = map[string]int{}

func Sequence(key string) any {
	data, ok := sequenceMap[key]
	if ok {
		data++
		sequenceMap[key] = data
		return data
	}
	sequenceMap[key] = 1
	return 1
}
