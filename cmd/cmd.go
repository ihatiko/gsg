package cmd

import (
	"fmt"
	cfg "gena/config"
	database_utils "gena/database-utils"
	"gena/postgres"
	"github.com/ihatiko/config"
	"github.com/jackc/pgx"
	"github.com/samber/lo"
)

var GlobalGenerationConfig *cfg.Configuration

const (
	configPath = "./config/config"
)

func Run() {
	globalGenerationConfig, err := config.GetConfig[cfg.Configuration](configPath)
	if err != nil {
		panic(err)
	}
	fmt.Println(globalGenerationConfig.Bulk)
	databases := []string{"test_database"}
	pgCfg := &postgres.Config{
		CreateDb: false,
		Host:     "localhost",
		Password: "postgres",
		PgDriver: "pgx",
		Port:     "5432",
		User:     "postgres",
		SSLMode:  "disable",
		Schema:   "public",
	}
	err = health(databases, pgCfg)
	if err != nil {
		panic(err)
	}
	err = scanDatabase(databases, pgCfg)
	if err != nil {
		panic(err)
	}
}

func scanDatabase(databases []string, cfg *postgres.Config) error {
	var err error
	for _, database := range databases {
		cfg.Dbname = database
		db, err := (cfg).NewConnection()

		if err != nil {
			return fmt.Errorf("database error %s %v", cfg.Dbname, err)
		}
		var Schemas []*database_utils.Schema
		err = db.Select(&Schemas, `select
    		table_catalog,
     		table_name,
     		column_name,
     		column_default,
     		data_type,
     		is_nullable,
     		character_maximum_length
 		from INFORMATION_SCHEMA.COLUMNS where table_schema = 'public'`)

		if err != nil {
			return fmt.Errorf("database error %s %v", cfg.Dbname, err)
		}
		tableGroup := lo.GroupBy[*database_utils.Schema, string](Schemas, func(item *database_utils.Schema) string {
			return item.TableName
		})
		for k, v := range tableGroup {
			columns, values := getColumnData(v)
			rows := [][]any{values}
			bulkData := pgx.CopyFromRows(rows)
			_, err := postgres.Connection.CopyFrom([]string{k}, columns, bulkData)
			if err != nil {
				panic(err)
			}
		}
	}
	return err
}

func getColumnData(schemas []*database_utils.Schema) ([]string, []any) {
	var columns []string
	var values []any
	for _, column := range schemas {
		columns = append(columns, column.ColumnName)
		values = append(values, database_utils.GetValue(column))
	}
	return columns, values
}

func health(databases []string, cfg *postgres.Config) error {
	var err error
	for _, database := range databases {
		cfg.Dbname = database
		conn, err := (cfg).NewConnection()
		if err != nil {
			err = fmt.Errorf("database error %s %v", cfg.Dbname, err)
		}
		conn.Close()
		postgres.Connection.Close()
	}
	return err
}
