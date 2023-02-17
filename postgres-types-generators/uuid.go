package postgres_types_generators

import uuid "github.com/satori/go.uuid"

func UUIDGenerator() any {
	u2 := uuid.NewV4()
	data, err := u2.MarshalBinary()
	if err != nil {
		panic(err)
	}
	return data
}
