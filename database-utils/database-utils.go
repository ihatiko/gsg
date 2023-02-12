package database_utils

import (
	"fmt"
	"gena/generators"
	"github.com/brianvoe/gofakeit/v6"
	"strings"
)

type Schema struct {
	Database      string  `db:"table_catalog"`
	TableName     string  `db:"table_name"`
	ColumnName    string  `db:"column_name"`
	ColumnDefault *string `db:"column_default"`
	DataType      string  `db:"data_type"`
	Length        *string `db:"character_maximum_length"`
	IsNullable    string  `db:"is_nullable"`
}

func GetValue(columnData *Schema) any {
	var result any
	switch columnData.DataType {
	case "integer":
		if columnData.ColumnDefault != nil && strings.Contains(*columnData.ColumnDefault, "nextval") {
			result = generators.Sequence(fmt.Sprintf("%s_%s_%s",
				columnData.Database,
				columnData.TableName,
				columnData.ColumnName,
			))
			break
		}
		result = gofakeit.IntRange(0, +2147483647)
	case "character varying":
		/*		data := columnData.Length
				parsedLength, _ := strconv.Atoi(*data)*/
		result = generators.RandStringRunes(10)
	default:
		return "unknown type"
	}

	return result
}
