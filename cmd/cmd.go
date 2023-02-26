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

var Settings *cfg.Settings

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
	for _, connection := range cfg.Settings.Connections {
		fmt.Println(connection)
	}
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
		conn, err := connection.Connection.NewConnection()
		if err != nil {
			resultError = append(resultError, fmt.Errorf("database error %s %s %v", connection.Name, connection.Name, err))
		}
		err = conn.Select(&databases, getDatabasesQuery)
		if err != nil {
			panic(err)
		}
		//TODO добавить извлечение типа
	}
	if len(resultError) == 0 {
		return nil
	}
	errorFormatted := lo.Map(resultError, func(item error, index int) string {
		return item.Error()
	})
	return errors.New(strings.Join(errorFormatted, "\n"))
}
func ProcessDatabase(connection cfg.DatabaseConnection) error {
	/*	var err error
		for _, database := range Settings.Databases {
			postgresConfig.Dbname = database.Name
			db, err := (postgresConfig).NewConnection()

			if err != nil {
				return fmt.Errorf("database error %s %v", postgresConfig.Dbname, err)
			}
			var Schemas []*generator.Schema
			err = db.Select(&Schemas, getDatabaseInfoQuery)
			if err != nil {
				return fmt.Errorf("database error %s %v", postgresConfig.Dbname, err)
			}
			generator := generator.NewGenerator(db, Settings)
			if ValidateSupportedTypes(Schemas, generator) {
				return nil
			}

			relations := ToRelations(Schemas)
			for k, v := range relations {
				InsertData(k, v, relations, database, generator)
			}
		}
		return err*/
}
func ValidateDictionary() {

}
func ValidateSupportedTypes(Schemas []*generator.Schema, g *generator.Generator) bool {
	state := false
	for _, v := range Schemas {
		g := generator.ColumnGenerator{Settings: Settings}
		_, err := g.GetValue(&generator.Column{Schema: v, GeneratedData: nil}, nil)
		if err != nil {
			fmt.Println(err)
			state = true
		}
	}
	return state
}

func InsertData(tableName string, v *generator.Table, relations map[string]*generator.Table, database *cfg.Database, gen *generator.Generator) {
	if v.Inserted {
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
			InsertData(dependencyTable.DependencyTableName, relations[dependencyTable.DependencyTableName], relations, database, gen)
		}
	}

	metaTable, _ := lo.Find[*cfg.Table](database.Tables, func(item *cfg.Table) bool {
		return item.Name == tableName
	})
	_, err := postgres.Connection.Exec(fmt.Sprintf("truncate %s cascade", tableName))
	if err != nil {
		panic(err)
	}
	dataSet := gen.GetColumnData(metaTable, v.Columns, relations)
	bulkData := pgx.CopyFromRows(dataSet.Data)
	_, err = postgres.Connection.CopyFrom([]string{tableName}, dataSet.Columns, bulkData)
	if err != nil {
		panic(err)
	}
	fmt.Println(fmt.Sprintf("Inserted in database %s table %s", database.Name, tableName))
	v.Inserted = true
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
			formattedColumns[wrappedColumn.Schema.ColumnName] = wrappedColumn
		}

		result[key] = &generator.Table{
			Columns:  formattedColumns,
			Inserted: false,
		}
	}
	return result
}

func health(connections []cfg.DatabaseConnection) error {
	var resultError []error
	for _, connection := range connections {
		conn, err := connection.Connection.NewConnection()
		if err != nil {
			resultError = append(resultError, fmt.Errorf("database error %s %s %v", connection.Name, connection.Name, err))
		}
		conn.Close()
		postgres.Connection.Close()
	}
	if len(resultError) == 0 {
		return nil
	}
	errorFormatted := lo.Map(resultError, func(item error, index int) string {
		return item.Error()
	})
	return errors.New(strings.Join(errorFormatted, "\n"))
}
