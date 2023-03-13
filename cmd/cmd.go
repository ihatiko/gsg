package cmd

import (
	"fmt"
	"github.com/ihatiko/config"
	"github.com/jackc/pgx"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	cfg "gsg/config"
	"gsg/generator"
	"gsg/postgres"
	"strings"
)

const (
	configPath = "./config/config"
)

type SchemasWrapper struct {
	Schemas []*generator.Schema
	Db      *postgres.ConnectionSet
}

var Settings *cfg.Settings
var ColumnBuffer = map[string]*SchemasWrapper{}
var BlackList = map[string]struct{}{}

func Run() {
	cfg, err := config.GetConfig[cfg.Config](configPath)
	if err != nil {
		panic(err)
	}
	Settings = cfg.Settings
	for _, blackListItem := range Settings.BlackListPath {
		data := strings.Split(blackListItem, Settings.Separator)
		if len(data) > 1 {
			table := ""
			if len(data) >= 2 {
				table = data[2]
			}
			BlackList[fmt.Sprintf("%s%s%s", data[1], Settings.Separator, table)] = struct{}{}
		}

	}
	err = health(Settings.Connections)
	if err != nil {
		panic(errors.Wrap(err, "Could not connect to databases"))
	}
	err = ValidateSupportedDatabaseTypes(Settings.Connections)
	if err != nil {
		panic(errors.Wrap(err, "Does not support types"))
	}
	PrepareDataSets()
	InsertData()
}

var Inserted = map[string]struct{}{}
var GroupedTables = map[string][]*generator.ColumnGenerator{}

func InsertData() {
	GroupedTables = lo.GroupBy[*generator.ColumnGenerator, string](lo.Values(generator.Generators), func(item *generator.ColumnGenerator) string {
		return fmt.Sprintf("%s%s%s", item.Column.Schema.Database, Settings.Separator, item.Column.Schema.TableName)
	})
	for _, gen := range GroupedTables {
		_, err := gen[0].Db.PgxConn.Exec(fmt.Sprintf("truncate %s cascade", gen[0].Column.Schema.TableName))
		if err != nil {
			panic(err)
		}
	}
	for k, gen := range GroupedTables {
		for _, g := range gen {
			if d, ok := g.Column.Constraints["FOREIGN KEY"]; ok {
				sd := generator.GeneratedSets[d.GetDependencyKey()].ToStringGeneratedData[0]
				g2 := generator.GeneratedSets[g.Column.Schema.GetKey()].ToStringGeneratedData[0]
				if sd != g2 && g2 != "null" {
					for k, std := range generator.GeneratedSets {
						fmt.Println(k, std.ToStringGeneratedData)
					}
					panic("TEST")
				}
			}
		}
		InsertSingleTable(k, gen)
	}
}

func InsertSingleTable(k string, gen []*generator.ColumnGenerator) {
	if _, ok := Inserted[k]; ok {
		return
	}
	gData := generator.GeneratedSets[gen[0].Column.Schema.GetKey()]
	dataSet := generator.DataSet{
		Columns: make([]string, len(gen)),
		Data:    make([][]any, len(gData.GeneratedData))}
	dataSetString := generator.DataSet{
		Columns: make([]string, len(gen)),
		Data:    make([][]any, len(gData.GeneratedData))}

	for j, tableColumns := range gen {
		for kConstraint, constraint := range tableColumns.Column.Constraints {
			if kConstraint == "FOREIGN KEY" && constraint.DependencyTableName != tableColumns.Column.Schema.TableName {
				key := fmt.Sprintf("%s%s%s", constraint.DependencyDatabaseName, Settings.Separator, constraint.DependencyTableName)
				if _, ok := Inserted[key]; !ok {
					InsertSingleTable(
						key,
						GroupedTables[key],
					)
				}
			}
		}
		genData := generator.GeneratedSets[tableColumns.Column.Schema.GetKey()]
		dataSet.Columns[j] = tableColumns.Column.Schema.ColumnName
		for i, columnValue := range genData.GeneratedData {
			if len(dataSet.Data[i]) == 0 {
				dataSet.Data[i] = make([]any, len(gen))
				dataSetString.Data[i] = make([]any, len(gen))
			}
			dataSet.Data[i][j] = columnValue
			dataSetString.Data[i][j] = genData.ToStringGeneratedData[i]
		}

	}
	fmt.Println("Dataset ------------------------", gen[0].Column.Schema.Database, gen[0].Column.Schema.TableName)
	fmt.Println(dataSet.Columns)
	for _, dataSetValue := range dataSetString.Data {
		fmt.Println(dataSetValue)
	}
	fmt.Println("------------------------ \n")
	bulkData := pgx.CopyFromRows(dataSet.Data)
	_, err := gen[0].Db.PgxConn.CopyFrom([]string{gen[0].Column.Schema.TableName}, dataSet.Columns, bulkData)
	if err != nil {
		detailedError := err.(pgx.PgError)
		panic(fmt.Sprintf("%s \n %v", detailedError.Detail, detailedError.Error()))
	}
	fmt.Println(fmt.Sprintf("Inserted in database %s table %s \n", gen[0].Column.Schema.Database, gen[0].Column.Schema.TableName))
	Inserted[k] = struct{}{}
}

func ValidateSupportedDatabaseTypes(connections []cfg.DatabaseConnection) error {
	var resultError []error
	for _, connection := range connections {
		var databases []string
		conn, err := connection.Connection.GetConnection()
		if err != nil {
			resultError = append(resultError, fmt.Errorf("database error %s %s \n%v", connection.Name, connection.Name, err))
		}
		err = conn.SqlxConn.Select(&databases, getDatabasesQuery)
		if err != nil {
			panic(err)
		}
		for _, dbName := range databases {
			conn, err = connection.Connection.ChangeDb(dbName).GetConnection()
			if err != nil {
				panic(err)
			}
			var Schemas []*generator.Schema
			err = conn.SqlxConn.Select(&Schemas, getDatabaseInfoQuery)
			if err != nil {
				resultError = append(resultError, fmt.Errorf("database error %s %v", dbName, err))
				continue
			}
			//TODO ToRelations
			if err := ValidateSupportedTypes(Schemas, conn); err != nil {
				resultError = append(resultError, fmt.Errorf("ValidateSupportedTypes error %s %s \n%v", connection.Name, connection.Name, err))
			}
			ColumnBuffer[dbName] = &SchemasWrapper{
				Db:      conn,
				Schemas: Schemas,
			}
		}
	}
	if len(resultError) == 0 {
		return nil
	}
	errorFormatted := lo.Map(resultError, func(item error, index int) string {
		return item.Error()
	})
	return errors.New(strings.Join(errorFormatted, "\n"))
}

func PrepareDataSets() {
	for _, columns := range ColumnBuffer {
		relations := ToRelations(columns.Schemas)
		for _, v := range relations {
			for _, column := range v.Columns {
				//TODO превратить заранее в мапу для оптимизации
				database, _ := lo.Find[*cfg.Database](Settings.Databases, func(item *cfg.Database) bool {
					return item.Name == column.Schema.Database
				})
				if database == nil {
					generator.SetGenerator(Settings, column, columns.Db, nil, v)
					continue
				}
				//TODO превратить заранее в мапу для оптимизации
				table, _ := lo.Find[*cfg.Table](database.Tables, func(item *cfg.Table) bool {
					return item.Name == column.Schema.TableName
				})
				generator.SetGenerator(Settings, column, columns.Db, table, v)
			}

			FillDataSet(relations)
		}
	}
}

func FillDataSet(relations map[string]*generator.Table) {
	for _, relation := range relations {
		for _, column := range relation.Columns {
			generator.GetGenerator(column).FillData()
		}
	}
}

func ValidateSupportedTypes(Schemas []*generator.Schema, db *postgres.ConnectionSet) error {
	var resultError []error
	for _, v := range Schemas {
		if _, ok := BlackList[fmt.Sprintf("%s%s%s", v.Database, Settings.Separator, v.TableName)]; ok {
			continue
		}
		g := generator.SetGenerator(Settings, &generator.Column{Schema: v}, db, nil, nil)
		_, _, err := g.GetValue()
		if err != nil {
			resultError = append(resultError, err)
		}
	}
	if len(resultError) == 0 {
		return nil
	}
	errorFormatted := lo.Map(resultError, func(item error, index int) string {
		return item.Error()
	})
	return errors.New(strings.Join(errorFormatted, "\n"))
}

func ToRelations(schemas []*generator.Schema) map[string]*generator.Table {
	result := map[string]*generator.Table{}
	tableGroup := lo.GroupBy[*generator.Schema, string](schemas, func(item *generator.Schema) string {
		return item.TableName
	})
	for key, tables := range tableGroup {
		columns := lo.Map[*generator.Schema, *generator.Column](tables, func(item *generator.Schema, index int) *generator.Column {
			return &generator.Column{Schema: item}
		})
		formattedColumns := map[string]*generator.Column{}
		wrappedColumns := lo.GroupBy[*generator.Column, string](columns, func(item *generator.Column) string {
			return item.Schema.ColumnName
		})

		for _, data := range wrappedColumns {
			wrappedColumn := data[0]
			if _, ok := BlackList[fmt.Sprintf("%s%s%s", wrappedColumn.Schema.Database, Settings.Separator, wrappedColumn.Schema.TableName)]; ok {
				break
			}
			wrappedColumn.Constraints = map[string]generator.Constraint{}
			for _, constraints := range data {
				if constraints.Schema.ConstraintType != nil {
					wrappedColumn.Constraints[*constraints.Schema.ConstraintType] = generator.Constraint{
						DependencyColumnName:   *constraints.Schema.DependencyColumnName,
						DependencyTableName:    *constraints.Schema.DependencyTableName,
						DependencyDatabaseName: constraints.Schema.Database,
					}
				}
			}
			wrappedColumn.Schema.DependencyColumnName = nil
			wrappedColumn.Schema.DependencyTableName = nil
			wrappedColumn.Schema.ConstraintType = nil
			formattedColumns[wrappedColumn.Schema.ColumnName] = wrappedColumn
		}

		result[key] = &generator.Table{
			Name:    key,
			Columns: formattedColumns,
		}
	}
	return result
}

func health(connections []cfg.DatabaseConnection) error {
	var resultError []error
	for _, connection := range connections {
		_, err := connection.Connection.GetConnection()
		if err != nil {
			resultError = append(resultError, fmt.Errorf("database error %s %s %v", connection.Name, connection.Name, err))
		}
	}
	if len(resultError) == 0 {
		return nil
	}
	errorFormatted := lo.Map(resultError, func(item error, index int) string {
		return item.Error()
	})
	return errors.New(strings.Join(errorFormatted, "\n"))
}
