package postgres_types_generators

import "github.com/brianvoe/gofakeit/v6"

// TODO json parser and generator
func JsonBGenerator() (any, string) {
	//TODO date-range-rule
	data, _ := gofakeit.JSON(&gofakeit.JSONOptions{
		Type: "object",
		Fields: []gofakeit.Field{
			{Name: "first_name", Function: "firstname"},
			{Name: "last_name", Function: "lastname"},
			{Name: "address", Function: "address"},
			{Name: "password", Function: "password", Params: gofakeit.MapParams{"special": {"false"}}},
		},
		Indent: true,
	})
	return data, string(data)
}
