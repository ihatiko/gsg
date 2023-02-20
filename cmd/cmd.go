package cmd

import (
	"fmt"
	"github.com/ihatiko/config"
	"github.com/jackc/pgx"
	"github.com/samber/lo"
	cfg "gsg/config"
	"gsg/generator"
	"gsg/postgres"
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
	for _, connections := cfg {
		Settings = cfg.Settings
		pgCfg := &postgres.Config{
			Host:     "localhost",
			Password: "postgres",
			PgDriver: "pgx",
			Port:     "5432",
			User:     "postgres",
			SSLMode:  "disable",
			Schema:   "public",
		}
		err = health(pgCfg)
		if err != nil {
			panic(err)
		}
		err = scanDatabase(pgCfg)
		if err != nil {
			panic(err)
		}
	}
}

func scanDatabase(postgresConfig *postgres.Config) error {
	var err error
	for _, database := range Settings.Databases {
		postgresConfig.Dbname = database.Name
		db, err := (postgresConfig).NewConnection()

		if err != nil {
			return fmt.Errorf("database error %s %v", postgresConfig.Dbname, err)
		}
		var Schemas []*generator.Schema
		err = db.Select(&Schemas, getDatabaseInfo)
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
	return err
}

func ValidateSupportedTypes(Schemas []*generator.Schema, g *generator.Generator) bool {
	state := false
	for _, v := range Schemas {
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

func health(cfg *postgres.Config) error {
	var err error
	for _, database := range Settings.Databases {
		cfg.Dbname = database.Name
		conn, err := (cfg).NewConnection()
		if err != nil {
			err = fmt.Errorf("database error %s %v", cfg.Dbname, err)
		}
		conn.Close()
		postgres.Connection.Close()
	}
	return err
}
