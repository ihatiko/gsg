package generator

import (
	"database/sql"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	cfg "gsg/config"
	postgres_types_generators "gsg/postgres-types-generators"
	"strings"
)

var Generators = map[string]*ColumnGenerator{}

type ColumnGenerator struct {
	Settings *cfg.Settings
	Table    *Table
	Column   *Column
	Enum     []string
	Db       *sqlx.DB
}

func NewColumnGenerator(settings *cfg.Settings, column *Column, db *sqlx.DB) *ColumnGenerator {
	return &ColumnGenerator{Settings: settings, Db: db, Column: column}
}

func SetGenerator(settings *cfg.Settings, column *Column, db *sqlx.DB) *ColumnGenerator {
	gen := NewColumnGenerator(settings, column, db)
	Generators[column.Schema.GetKey()] = gen
	return gen
}

func GetGenerator(settings *cfg.Settings, column *Column, db *sqlx.DB) *ColumnGenerator {
	if data, ok := Generators[column.Schema.GetKey()]; ok {
		return data
	}
	gen := NewColumnGenerator(settings, column, db)
	Generators[column.Schema.GetKey()] = gen
	return gen
}

func (g *ColumnGenerator) GetValue() (any, error) {
	var result any
	_, unique := g.Column.Constraints["UNIQUE"]
	switch g.Column.Schema.DataType {
	case "USER-DEFINED":
		return g.GetCustomDatabaseType(g.Column.Schema)
	case "date":
		result = postgres_types_generators.DateGenerator()
	case "timestamp with time zone":
		result = postgres_types_generators.TimeStampGenerator()
	case "timestamp without time zone":
		result = postgres_types_generators.TimeStampGenerator()
	case "boolean":
		result = postgres_types_generators.BoolGenerator()
	case "numeric":
		result = postgres_types_generators.NumericGenerator()
	case "uuid":
		result = postgres_types_generators.UUIDGenerator(unique)
	case "bit":
		result = postgres_types_generators.BitGenerator()
	case "jsonb":
		result = postgres_types_generators.JsonBGenerator()
	case "integer":
		if g.Column.Schema.ColumnDefault != nil && strings.Contains(*g.Column.Schema.ColumnDefault, "nextval") {
			result = postgres_types_generators.Serial(fmt.Sprintf("%s_%s_%s",
				g.Column.Schema.Database,
				g.Column.Schema.TableName,
				g.Column.Schema.ColumnName,
			))
			break
		}
		result = gofakeit.IntRange(0, +2147483647)
	case "text":
		result = postgres_types_generators.RandStringRunes(g.Settings.DefaultTypeSettings.VarCharLength, unique)
	case "smallint":
		result = postgres_types_generators.SmallIntGenerator()
	case "bigint":
		result = postgres_types_generators.BigIntGenerator()
	case "character varying":
		result = postgres_types_generators.RandStringRunes(g.Settings.DefaultTypeSettings.VarCharLength, unique)
	default:
		return nil, errors.New(fmt.Sprintf("unknown type %s in table %s in database %s",
			g.Column.Schema.DataType,
			g.Column.Schema.TableName,
			g.Column.Schema.Database,
		))
	}
	if g.Column.Schema.IsNullable {
		if gofakeit.IntRange(0, 1) == 0 {
			return nil, nil
		}
	}
	return result, nil
}

func (g *ColumnGenerator) GetCustomDatabaseType(columnData *Schema) (any, error) {
	enumKey := strings.Split(*columnData.ColumnDefault, "::")
	if len(enumKey) < 2 {
		return nil, errors.New(fmt.Sprintf("unknown type %s in table %s in database %s", columnData.DataType, columnData.TableName, columnData.Database))
	}
	if g.Enum == nil {
		var enum []string
		err := g.Db.Select(&enum, getEnumData, enumKey[1])
		if err != nil && err != sql.ErrNoRows {
			panic(err)
		}
		if len(enum) == 0 {
			return nil, errors.New(fmt.Sprintf("unknown type %s in table %s in database %s", columnData.DataType, columnData.TableName, columnData.Database))
		}
		g.Enum = enum
	}
	return postgres_types_generators.ByDictionaryStringRandom(g.Enum), nil
}
