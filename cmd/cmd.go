package cmd

import (
	"fmt"
	"github.com/ihatiko/config"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	cfg "gsg/config"
	"gsg/generator"
	"strings"
)

const (
	configPath = "./config/config"
)

type SchemasWrapper struct {
	Schemas []*generator.Schema
	Db      *sqlx.DB
}

var Settings *cfg.Settings
var ColumnBuffer = map[string]*SchemasWrapper{}

func Run() {
	cfg, err := config.GetConfig[cfg.Config](configPath)
	if err != nil {
		panic(err)
	}
	Settings = cfg.Settings
	err = health(Settings.Connections)
	if err != nil {
		panic(errors.Wrap(err, "Could not connect to databases"))
	}
	err = ValidateSupportedDatabaseTypes(Settings.Connections)
	if err != nil {
		panic(errors.Wrap(err, "Does not support types"))
	}
	PrepareDataSets()

	/*	for _, connection := range cfg.Settings.Connections {
		err = health(pgCfg)
		if err != nil {
			panic(err)
		}
		err = scanDatabase(pgCfg)
		if err != nil {
			panic(err)
		}
	}*/
}
func ValidateSupportedDatabaseTypes(connections []cfg.DatabaseConnection) error {
	var resultError []error
	for _, connection := range connections {
		var databases []string
		conn, err := connection.Connection.GetConnection()
		if err != nil {
			resultError = append(resultError, fmt.Errorf("database error %s %s \n%v", connection.Name, connection.Name, err))
		}
		err = conn.Select(&databases, getDatabasesQuery)
		if err != nil {
			panic(err)
		}
		for _, dbName := range databases {
			conn, err = connection.Connection.ChangeDb(dbName).GetConnection()
			if err != nil {
				panic(err)
			}
			var Schemas []*generator.Schema
			err = conn.Select(&Schemas, getDatabaseInfoQuery)
			if err != nil {
				resultError = append(resultError, fmt.Errorf("database error %s %v", dbName, err))
				continue
			}
			//TODO ToRelations
			if err := ValidateSupportedTypes(Schemas, conn); err != nil {
				resultError = append(resultError, fmt.Errorf("ValidateSupportedTypes error %s %s \n%v", connection.Name, connection.Name, err))
			}
			ColumnBuffer[conn] = &SchemasWrapper{
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
	for db, columns := range ColumnBuffer {
		relations := ToRelations(columns.Schemas)
		for _, v := range relations {
			for _, column := range v.Columns {
				generator.SetGenerator(Settings, column, db)
			}

			FillDataSet(relations)
		}
	}
}

func FillDataSet(relations map[string]*generator.Table) {
	for _, relation := range relations {
		for _, columns := range relation.Columns {

		}
	}
}

func ValidateSupportedTypes(Schemas []*generator.Schema, db *sqlx.DB) error {
	var resultError []error
	for _, v := range Schemas {
		g := generator.GetGenerator(Settings, &generator.Column{Schema: v}, db)
		_, err := g.GetValue()
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

func InsertData(relations map[string]*generator.Table) {
	/*	if v.Inserted {
			return
		}
		columns := lo.Values(v.Columns)
		dependencyConstraints := lo.Reduce[*generator.Column, []generator.Constraint](columns, func(agg []generator.Constraint, item *generator.Column, index int) []generator.Constraint {
			for k, constraint := range item.Constraints {
				if k == "FOREIGN KEY" && tableName != constraint.DependencyTableName {
					agg = append(agg, constraint)
				}
			}
			return agg
		}, []generator.Constraint{})

		for _, dependencyTable := range dependencyConstraints {
			dpd := relations[dependencyTable.DependencyTableName]
			if !dpd.Inserted {
				fmt.Println(fmt.Sprintf("Redirect create %s -> %s", tableName, dependencyTable.DependencyTableName))
				InsertData(dependencyTable.DependencyTableName, relations[dependencyTable.DependencyTableName], relations)
			}
		}

		metaTable, _ := lo.Find[*cfg.Table](database.Tables, func(item *cfg.Table) bool {
			return item.Name == tableName
		})
		_, err := postgres.Connection.Exec(fmt.Sprintf("truncate %s cascade", tableName))
		if err != nil {
			panic(err)
		}
		dataSet := generator.GetColumnData(metaTable, v.Columns, relations)
		bulkData := pgx.CopyFromRows(dataSet.Data)
		_, err = postgres.Connection.CopyFrom([]string{tableName}, dataSet.Columns, bulkData)
		if err != nil {
			panic(err)
		}
		fmt.Println(fmt.Sprintf("Inserted in database %s table %s", database.Name, tableName))
		v.Inserted = true*/
}

func ToRelations(schemas []*generator.Schema) map[string]*generator.Table {
	result := map[string]*generator.Table{}
	tableGroup := lo.GroupBy[*generator.Schema, string](schemas, func(item *generator.Schema) string {
		return item.TableName
	})
	for key, tables := range tableGroup {
		columns := lo.Map[*generator.Schema, *generator.Column](tables, func(item *generator.Schema, index int) *generator.Column {
			return &generator.Column{Schema: item, GeneratedData: nil}
		})
		formattedColumns := map[string]*generator.Column{}
		wrappedColumns := lo.GroupBy[*generator.Column, string](columns, func(item *generator.Column) string {
			return item.Schema.ColumnName
		})

		for _, data := range wrappedColumns {
			wrappedColumn := data[0]
			wrappedColumn.Constraints = map[string]generator.Constraint{}
			for _, constraints := range data {
				if constraints.Schema.ConstraintType != nil {
					wrappedColumn.Constraints[*constraints.Schema.ConstraintType] = generator.Constraint{
						DependencyColumnName: *constraints.Schema.DependencyColumnName,
						DependencyTableName:  *constraints.Schema.DependencyTableName,
					}
				}
			}
			wrappedColumn.Schema.DependencyColumnName = nil
			wrappedColumn.Schema.DependencyTableName = nil
			wrappedColumn.Schema.ConstraintType = nil
			formattedColumns[wrappedColumn.Schema.ColumnName] = wrappedColumn
		}

		result[key] = &generator.Table{
			Name:     key,
			Columns:  formattedColumns,
			Inserted: false,
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
