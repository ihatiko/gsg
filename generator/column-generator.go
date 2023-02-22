package generator

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/pkg/errors"
	cfg "gsg/config"
	postgres_types_generators "gsg/postgres-types-generators"
	"strings"
)

type ColumnGenerator struct {
	Settings *cfg.Settings
}

func (g *ColumnGenerator) GetValue(columnData *Column, table *cfg.Table) (any, error) {
	var result any
	_, unique := columnData.Constraints["UNIQUE"]
	switch columnData.Schema.DataType {
	case "USER-DEFINED":
		return g.GetCustomDatabaseType(columnData.Schema, table)
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
		if columnData.Schema.ColumnDefault != nil && strings.Contains(*columnData.Schema.ColumnDefault, "nextval") {
			result = postgres_types_generators.Serial(fmt.Sprintf("%s_%s_%s",
				columnData.Schema.Database,
				columnData.Schema.TableName,
				columnData.Schema.ColumnName,
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
			columnData.Schema.DataType,
			columnData.Schema.TableName,
			columnData.Schema.Database,
		))
	}
	if columnData.Schema.IsNullable {
		if gofakeit.IntRange(0, 1) == 0 {
			return nil, nil
		}
	}
	return result, nil
}

func (g *ColumnGenerator) GetCustomDatabaseType(columnData *Schema, table *cfg.Table) (any, error) {
	/*	if g.Enums == nil {
			var enum []*Enum
			g.Enums = map[string]map[string][]string{}
			err := g.Db.Select(&enum, getEnumData)
			if err != nil && err != sql.ErrNoRows {
				panic(err)
			}
			if enum != nil && len(enum) > 0 {
				parsedEnum := map[string][]string{}
				for _, v := range enum {
					parsedEnum[v.Type] = append(parsedEnum[v.Type], v.Value)
				}
				g.Enums[columnData.Database] = parsedEnum
			}
		}
		if v, ok := g.Enums[columnData.Database]; ok && columnData.ColumnDefault != nil {
			enumKey := strings.Split(*columnData.ColumnDefault, "::")
			if len(enumKey) < 2 {
				return nil, errors.New(fmt.Sprintf("unknown type %s in table %s in database %s", columnData.DataType, columnData.TableName, columnData.Database))
			}
			enum := v[enumKey[1]]
			return enum[rand.Intn(len(enum))], nil
		}*/
	return nil, errors.New(fmt.Sprintf("unknown type %s in table %s in database %s", columnData.DataType, columnData.TableName, columnData.Database))
}
