package generators

import uuid "github.com/satori/go.uuid"

func UUIDGenerator() any {
	u2 := uuid.NewV4()
	data, _ := u2.MarshalBinary()
	return data
}
