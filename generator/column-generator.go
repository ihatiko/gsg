package generator

import (
	"database/sql"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	cfg "gsg/config"
	"gsg/postgres"
	postgres_types_generators "gsg/postgres-types-generators"
	"strings"
)

var Generators = map[string]*ColumnGenerator{}
var GeneratedSets = map[string]*GeneratedData{}

type DataSet struct {
	Columns []string
	Data    [][]any
}
type GeneratedData struct {
	GeneratedData         []any
	ToStringGeneratedData []string
}
type ColumnGenerator struct {
	Settings      *cfg.Settings
	Table         *Table
	Column        *Column
	Enum          []string
	Db            *postgres.ConnectionSet
	TableSettings *cfg.Table
}

func NewColumnGenerator(settings *cfg.Settings, column *Column, db *postgres.ConnectionSet, tableSettings *cfg.Table, table *Table) *ColumnGenerator {
	return &ColumnGenerator{Settings: settings, Db: db, Column: column, TableSettings: tableSettings, Table: table}
}

func SetGenerator(settings *cfg.Settings, column *Column, db *postgres.ConnectionSet, tableSettings *cfg.Table, table *Table) *ColumnGenerator {
	gen := NewColumnGenerator(settings, column, db, tableSettings, table)
	Generators[column.Schema.GetKey()] = gen
	return gen
}

func GetGenerator(column *Column) *ColumnGenerator {
	return Generators[column.Schema.GetKey()]
}

func (g *ColumnGenerator) FillData() {
	if _, ok := GeneratedSets[g.Column.Schema.GetKey()]; ok {
		return
	}
	count := 0
	if g.TableSettings != nil && g.TableSettings.Set != 0 {
		count = g.TableSettings.Set
	}
	if count == 0 {
		count = g.Settings.DefaultSet
	}
	if g.FillDependencyValues(count) {
		return
	}
	if _, ok := g.Column.Constraints["UNIQUE"]; ok {
		g.FillUniqueValues(count)
		return
	}
	g.FillRandomValues(count)
}

func (g *ColumnGenerator) FillUniqueValues(count int) {
	data := &GeneratedData{}
	for i := 0; i < count; i++ {
		val, valToString, _ := g.GetValue()
		data.GeneratedData = append(data.GeneratedData, val)
		data.ToStringGeneratedData = append(data.ToStringGeneratedData, valToString)
	}
	GeneratedSets[g.Column.Schema.GetKey()] = data
}
func (g *ColumnGenerator) FillRandomValues(count int) {
	data := &GeneratedData{}
	for i := 0; i < count; i++ {
		val, valToString, _ := g.GetValue()
		data.GeneratedData = append(data.GeneratedData, val)
		data.ToStringGeneratedData = append(data.ToStringGeneratedData, valToString)
	}
	GeneratedSets[g.Column.Schema.GetKey()] = data
}

type GenSet struct {
	Data         any
	ToStringData string
}

func (g *ColumnGenerator) FillDependencyValues(count int) bool {
	if len(g.Column.Constraints) > 0 {
		if data, ok := g.Column.Constraints["FOREIGN KEY"]; ok {
			key := data.GetDependencyKey()
			constrainGen := Generators[key]
			generatedSet, okGen := GeneratedSets[constrainGen.Column.Schema.GetKey()]
			if !okGen {
				constrainGen.FillData()
				generatedSet = GeneratedSets[constrainGen.Column.Schema.GetKey()]
			}

			var set []GenSet
			for {
				if count <= len(set) {
					break
				}
				for i := range generatedSet.GeneratedData {
					resultValue := generatedSet.GeneratedData[i]
					resultToValue := generatedSet.ToStringGeneratedData[i]
					if g.Column.Schema.IsNullable {
						if gofakeit.IntRange(0, 1) == 0 {
							resultValue = nil
							resultToValue = "null"
						}
					}
					set = append(set, GenSet{
						Data:         resultValue,
						ToStringData: resultToValue,
					})
				}
			}
			set = lo.Shuffle(set[:count])
			currentGeneratedSet := GeneratedData{}
			for _, setData := range set {
				currentGeneratedSet.GeneratedData = append(currentGeneratedSet.GeneratedData, setData.Data)
				currentGeneratedSet.ToStringGeneratedData = append(currentGeneratedSet.ToStringGeneratedData, setData.ToStringData)
			}

			GeneratedSets[g.Column.Schema.GetKey()] = &currentGeneratedSet
			return true
		}
	}
	return false
}

func (d *Schema) GetKey() string {
	return fmt.Sprintf("%s//%s//%s", d.Database, d.TableName, d.ColumnName)
}
func (d *Constraint) GetDependencyKey() string {
	return fmt.Sprintf("%s//%s//%s", d.DependencyDatabaseName, d.DependencyTableName, d.DependencyColumnName)
}

func (g *ColumnGenerator) GetValue() (any, string, error) {
	var result any
	var resultToString string
	switch g.Column.Schema.DataType {
	case "USER-DEFINED":
		return g.GetCustomDatabaseType(g.Column.Schema)
	case "date":
		result, resultToString = postgres_types_generators.DateGenerator()
	case "timestamp with time zone":
		result, resultToString = postgres_types_generators.TimeStampGenerator()
	case "timestamp without time zone":
		result, resultToString = postgres_types_generators.TimeStampGenerator()
	case "boolean":
		result, resultToString = postgres_types_generators.BoolGenerator()
	case "numeric":
		result, resultToString = postgres_types_generators.NumericGenerator()
	case "uuid":
		result, resultToString = postgres_types_generators.UUIDGenerator()
	case "bit":
		result, resultToString = postgres_types_generators.BitGenerator()
	case "jsonb":
		result, resultToString = postgres_types_generators.JsonBGenerator()
	case "integer":
		if g.Column.Schema.ColumnDefault != nil && strings.Contains(*g.Column.Schema.ColumnDefault, "nextval") {
			result, resultToString = postgres_types_generators.Serial(fmt.Sprintf("%s_%s_%s",
				g.Column.Schema.Database,
				g.Column.Schema.TableName,
				g.Column.Schema.ColumnName,
			))
			break
		}
		result = gofakeit.IntRange(0, +2147483647)
	case "text":
		result, resultToString = postgres_types_generators.RandStringRunes(g.Settings.DefaultTypeSettings.VarCharLength)
	case "smallint":
		result, resultToString = postgres_types_generators.SmallIntGenerator()
	case "bigint":
		result, resultToString = postgres_types_generators.BigIntGenerator()
	case "character varying":
		result, resultToString = postgres_types_generators.RandStringRunes(g.Settings.DefaultTypeSettings.VarCharLength)
	default:
		return nil, "", errors.New(fmt.Sprintf("unknown type %s in table %s in database %s",
			g.Column.Schema.DataType,
			g.Column.Schema.TableName,
			g.Column.Schema.Database,
		))
	}
	//TODO Добавить шанс срабатывания
	if g.Column.Schema.IsNullable {
		if gofakeit.IntRange(0, 1) == 0 {
			return nil, "null", nil
		}
	}
	return result, resultToString, nil
}

func (g *ColumnGenerator) GetCustomDatabaseType(columnData *Schema) (any, string, error) {
	enumKey := strings.Split(*columnData.ColumnDefault, "::")
	if len(enumKey) < 2 {
		return nil, "", errors.New(fmt.Sprintf("unknown type %s in table %s in database %s",
			columnData.DataType, columnData.TableName, columnData.Database))
	}
	if g.Enum == nil {
		var enum []string
		err := g.Db.SqlxConn.Select(&enum, getEnumData, enumKey[1])
		if err != nil && err != sql.ErrNoRows {
			panic(err)
		}
		if len(enum) == 0 {
			return nil, "", errors.New(fmt.Sprintf("unknown type %s in table %s in database %s",
				columnData.DataType, columnData.TableName, columnData.Database))
		}
		g.Enum = enum
	}
	result, stringResult := postgres_types_generators.ByDictionaryStringRandom(g.Enum)
	return result, stringResult, nil
}
