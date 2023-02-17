package postgres_types_generators

import uuid "github.com/satori/go.uuid"

var UUIDUNIQUE = map[string]struct{}{}

func UUIDGenerator(unique bool) any {
	u2 := getUniqueData(unique)
	data, err := u2.MarshalBinary()
	if err != nil {
		panic(err)
	}
	return data
}

func getUniqueData(unique bool) uuid.UUID {
	u2 := uuid.NewV4()
	if unique {
		if _, ok := UUIDUNIQUE[u2.String()]; !ok {
			return getUniqueData(unique)
		}
		UUIDUNIQUE[u2.String()] = struct{}{}
	}
	return u2
}
