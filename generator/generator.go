package generator

import (
	"database/sql"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	cfg "gsg/config"
	"gsg/postgres-types-generators"
	"math/rand"
	"sort"
	"strings"
)

func NewGenerator(db *sqlx.DB, settings *cfg.Settings) *Generator {
	return &Generator{Db: db, Settings: settings}
}

func (g *Generator) GetColumnData(table *cfg.Table, schemas map[string]*Column, relations map[string]*Table) DataSet {
	var columns []string

	for _, column := range schemas {
		columns = append(columns, column.Schema.ColumnName)
	}
	count := 0
	if table != nil && table.Set != 0 {
		count = table.Set
	}
	if count == 0 {
		count = g.Settings.DefaultSet
	}
	return g.generate(table, count, schemas, relations)
}

type DataSet struct {
	Columns []string
	Data    [][]any
}

func (g *Generator) generate(table *cfg.Table, count int, columns map[string]*Column, relations map[string]*Table) DataSet {
	dataset := DataSet{Data: [][]any{}, Columns: nil}
	formatted := lo.Values(columns)
	sort.Slice(formatted, func(i, j int) bool {
		return !(formatted[i].Schema.ConstraintType != nil && *formatted[i].Schema.ConstraintType == "FOREIGN KEY")
	})
	for _, column := range formatted {
		dataset.Columns = append(dataset.Columns, column.Schema.ColumnName)
		if column.Constraints != nil {
			if data, ok := column.Constraints["FOREIGN KEY"]; ok {
				dependencyTable := relations[data.DependencyTableName]
				dependencyColumn := dependencyTable.Columns[data.DependencyColumnName]
				var shuffledData []any
				for {
					if count <= len(shuffledData) {
						break
					}
					for _, d := range dependencyColumn.GeneratedData {
						if column.Schema.IsNullable {
							if gofakeit.IntRange(0, 1) == 0 {
								d = nil
							}
						}
						shuffledData = append(shuffledData, d)
					}
				}
				shuffledData = lo.Shuffle(shuffledData)
				shuffledSet := shuffledData[:count]
				for i := 0; i < len(shuffledSet); i++ {
					if len(dataset.Data) <= i {
						dataset.Data = append(dataset.Data, []any{shuffledSet[i]})
					} else {
						dataset.Data[i] = append(dataset.Data[i], shuffledSet[i])
					}
				}
				column.GeneratedData = append(column.GeneratedData, shuffledSet...)
				continue
			}
		}
		var columnSet []any
		for i := 0; i < count; i++ {
			val, _ := g.GetValue(column.Schema, table)
			columnSet = append(columnSet, val)
			if len(dataset.Data) <= i {
				dataset.Data = append(dataset.Data, []any{val})
			} else {
				dataset.Data[i] = append(dataset.Data[i], val)
			}
		}
		column.GeneratedData = append(column.GeneratedData, columnSet...)
	}
	return dataset
}

func (d *Schema) GetKey() string {
	return fmt.Sprintf("%s.%s.%s", d.Database, *d.DependencyTableName, *d.DependencyColumnName)
}

func (g *Generator) GetCustomDatabaseType(columnData *Schema, table *cfg.Table) (any, error) {
	if g.Enums == nil {
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
	}
	return nil, errors.New(fmt.Sprintf("unknown type %s in table %s in database %s", columnData.DataType, columnData.TableName, columnData.Database))
}

func (g *Generator) GetValue(columnData *Schema, table *cfg.Table) (any, error) {
	var result any
	switch columnData.DataType {
	case "USER-DEFINED":
		return g.GetCustomDatabaseType(columnData, table)
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
		result = postgres_types_generators.UUIDGenerator()
	case "bit":
		result = postgres_types_generators.BitGenerator()
	case "jsonb":
		result = postgres_types_generators.JsonBGenerator()
	case "integer":
		if columnData.ColumnDefault != nil && strings.Contains(*columnData.ColumnDefault, "nextval") {
			result = postgres_types_generators.Serial(fmt.Sprintf("%s_%s_%s",
				columnData.Database,
				columnData.TableName,
				columnData.ColumnName,
			))
			break
		}
		result = gofakeit.IntRange(0, +2147483647)
	case "text":
		result = postgres_types_generators.RandStringRunes(10)
	case "smallint":
		result = postgres_types_generators.SmallIntGenerator()
	case "bigint":
		result = postgres_types_generators.BigIntGenerator()
	case "character varying":
		result = postgres_types_generators.RandStringRunes(10)
	default:
		return nil, errors.New(fmt.Sprintf("unknown type %s in table %s in database %s", columnData.DataType, columnData.TableName, columnData.Database))
	}
	if columnData.IsNullable {
		if gofakeit.IntRange(0, 1) == 0 {
			return nil, nil
		}
	}
	return result, nil
}
