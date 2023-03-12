package postgres_types_generators

import (
	"github.com/jackc/pgx/pgtype"
	"math/rand"
)

var various = []any{
	&pgtype.Bit{Bytes: []byte{1}, Len: 1, Status: pgtype.Present},
}

func BitGenerator() (any, string) {
	data := various[rand.Intn(len(various))]
	return data, "&pgtype.Bit"
}
