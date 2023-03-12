package postgres_types_generators

import uuid "github.com/satori/go.uuid"

func UUIDGenerator() (any, string) {
	u2 := uuid.NewV4()
	data, _ := u2.MarshalBinary()

	return data, u2.String()
}
