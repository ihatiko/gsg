package database_utils

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	cfg "gsg/config"
	"gsg/generators"
	"strings"
)

type Schema struct {
	Database             string  `db:"table_catalog"`
	TableName            string  `db:"table_name"`
	ColumnName           string  `db:"column_name"`
	ColumnDefault        *string `db:"column_default"`
	DataType             string  `db:"data_type"`
	Length               *string `db:"character_maximum_length"`
	IsNullable           string  `db:"is_nullable"`
	DependencyTableName  *string `db:"dependency_table_name"`
	DependencyColumnName *string `db:"dependency_column_name"`
}

func (d *Schema) GetKey() string {
	return fmt.Sprintf("%s.%s.%s", d.Database, *d.DependencyTableName, *d.DependencyColumnName)
}
func GetValue(columnData *Schema, table *cfg.Table) any {
	var result any
	switch columnData.DataType {
	case "date":
		return generators.DateGenerator()
	case "timestamp without time zone":
		return generators.TimeStampGenerator()
	case "boolean":
		return generators.BoolGenerator()
	case "numeric":
		return generators.NumericGenerator()
	case "uuid":
		return generators.UUIDGenerator()
	case "bit":
		return generators.BitGenerator()
	case "integer":
		if columnData.ColumnDefault != nil && strings.Contains(*columnData.ColumnDefault, "nextval") {
			result = generators.Serial(fmt.Sprintf("%s_%s_%s",
				columnData.Database,
				columnData.TableName,
				columnData.ColumnName,
			))
			break
		}
		result = gofakeit.IntRange(0, +2147483647)
	case "character varying":
		result = generators.RandStringRunes(10)
	default:
		panic(fmt.Sprintf("unknown type %s", columnData.DataType))
		return "unknown type"
	}

	return result
}
