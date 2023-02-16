package cmd

import (
	"fmt"
	"github.com/ihatiko/config"
	"github.com/jackc/pgx"
	"github.com/samber/lo"
	cfg "gsg/config"
	databaseUtils "gsg/database-utils"
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

func scanDatabase(postgresConfig *postgres.Config) error {
	var err error
	for _, database := range Settings.Databases {
		postgresConfig.Dbname = database.Name
		db, err := (postgresConfig).NewConnection()

		if err != nil {
			return fmt.Errorf("database error %s %v", postgresConfig.Dbname, err)
		}
		var Schemas []*databaseUtils.Schema
		err = db.Select(&Schemas, getDatabaseInfo)

		if err != nil {
			return fmt.Errorf("database error %s %v", postgresConfig.Dbname, err)
		}
		relations := ToRelations(Schemas)
		for k, v := range relations {
			InsertData(k, v, relations, database)
		}
	}
	return err
}

func InsertData(tableName string, v *Table, relations map[string]*Table, database *cfg.Database) {
	if v.Inserted {
		return
	}
	for _, column := range v.Columns {
		if column.Schema.DependencyColumnName != nil {
			table := relations[*column.Schema.DependencyTableName]
			InsertData(*column.Schema.DependencyTableName, table, relations, database)
		}
	}
	metaTable, _ := lo.Find[*cfg.Table](database.Tables, func(item *cfg.Table) bool {
		return item.Name == tableName
	})

	columns, values := getColumnData(metaTable, v.Columns, relations)
	for i := 0; i < 10; i++ {
		bulkData := pgx.CopyFromRows(values)
		_, err := postgres.Connection.CopyFrom([]string{tableName}, columns, bulkData)
		fmt.Println(fmt.Sprintf("Inserted in database %s table %s data: %d", database.Name, tableName, len(values)))
		if err == nil {
			break
		}
		if i == 10 {
			panic(err)
		}
	}
	v.Inserted = true
}

type Table struct {
	Inserted bool
	Columns  map[string]*Column
}

type Column struct {
	GeneratedData []any
	Schema        *databaseUtils.Schema
}

func ToRelations(schemas []*databaseUtils.Schema) map[string]*Table {
	result := map[string]*Table{}
	tableGroup := lo.GroupBy[*databaseUtils.Schema, string](schemas, func(item *databaseUtils.Schema) string {
		return item.TableName
	})
	for key, tables := range tableGroup {
		columns := lo.Map[*databaseUtils.Schema, *Column](tables, func(item *databaseUtils.Schema, index int) *Column {
			return &Column{Schema: item, GeneratedData: nil}
		})
		formattedColumns := map[string]*Column{}
		for _, data := range columns {
			formattedColumns[data.Schema.ColumnName] = data
		}
		result[key] = &Table{
			Columns:  formattedColumns,
			Inserted: false,
		}
	}
	return result
}

// TODO pk shuffle
func getColumnData(table *cfg.Table, schemas map[string]*Column, relations map[string]*Table) ([]string, [][]any) {
	var columns []string

	for _, column := range schemas {
		columns = append(columns, column.Schema.ColumnName)
	}
	count := 0
	if table != nil && table.Set != 0 {
		count = table.Set
	}
	if count == 0 {
		count = Settings.DefaultSet
	}
	return columns, Generate(table, count, schemas, relations)
}

func Generate(table *cfg.Table, count int, columns map[string]*Column, relations map[string]*Table) [][]any {
	var values [][]any
	for _, column := range columns {
		if column.Schema.DependencyColumnName != nil {
			dependencyTable := relations[*column.Schema.DependencyTableName]
			dependencyColumn := dependencyTable.Columns[*column.Schema.DependencyColumnName]
			shuffledData := dependencyColumn.GeneratedData
			for {
				if count <= len(shuffledData) {
					break
				}
				shuffledData = append(shuffledData, dependencyColumn.GeneratedData...)
			}
			shuffledData = lo.Shuffle(shuffledData)
			shuffledSet := shuffledData[:count]
			for i := 0; i < len(shuffledSet); i++ {
				if len(values) <= i {
					values = append(values, []any{shuffledSet[i]})
				} else {
					values[i] = append(values[i], shuffledSet[i])
				}
			}
			column.GeneratedData = append(column.GeneratedData, shuffledSet...)
			continue
		}
		var columnSet []any
		for i := 0; i < count; i++ {
			val := databaseUtils.GetValue(column.Schema, table)
			columnSet = append(columnSet, val)
			if len(values) <= i {
				values = append(values, []any{val})
			} else {
				values[i] = append(values[i], val)
			}
		}
		column.GeneratedData = append(column.GeneratedData, columnSet...)
	}
	return values
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
