package generator

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
	cfg "gsg/config"
	"sort"
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
		/*		var columnSet []any
				for i := 0; i < count; i++ {
					val, _ := g.GetValue(column, table)
					columnSet = append(columnSet, val)
					if len(dataset.Data) <= i {
						dataset.Data = append(dataset.Data, []any{val})
					} else {
						dataset.Data[i] = append(dataset.Data[i], val)
					}
				}
				column.GeneratedData = append(column.GeneratedData, columnSet...)*/
	}
	return dataset
}

func (d *Schema) GetKey() string {
	return fmt.Sprintf("%s.%s.%s", d.Database, *d.DependencyTableName, *d.DependencyColumnName)
}
