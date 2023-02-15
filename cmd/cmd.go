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
		tableGroup := lo.GroupBy[*databaseUtils.Schema, string](Schemas, func(item *databaseUtils.Schema) string {
			return item.TableName
		})
		for k, v := range tableGroup {
			metaTable, _ := lo.Find[*cfg.Table](database.Tables, func(item *cfg.Table) bool {
				return item.Name == k
			})
			for i := 0; i < metaTable.Set/Settings.MaxBatch; i++ {
				columns, values := getColumnData(metaTable, v)
				bulkData := pgx.CopyFromRows(values)
				_, err := postgres.Connection.CopyFrom([]string{k}, columns, bulkData)
				if err != nil {
					panic(err)
				}
			}
		}
	}
	return err
}

type Relation struct {
	GeneratedData map[string][]any
	Schema        *databaseUtils.Schema
}

type RelationSchemaLink struct {
	Schema *databaseUtils.Schema
}
type RelationGraph struct {
	AllRelations        map[string]Relation
	RelationSchemaLinks []*RelationSchemaLink
}

func BuildRelations(data []*databaseUtils.Schema) *RelationGraph {
	graph := RelationGraph{AllRelations: map[string]Relation{}}
	for _, d := range data {
		schemaLink := &RelationSchemaLink{}
		schemaLink.Schema = d
		if d.DependencyTableName != nil {
			dependencySchema, _ := lo.Find[*databaseUtils.Schema](data, func(item *databaseUtils.Schema) bool {
				return item.TableName == *d.DependencyTableName && item.TableName == *d.DependencyColumnName
			})
			currentRelation, _ := lo.Find[*databaseUtils.Schema](graph.RelationSchemaLinks, func(item *databaseUtils.Schema) bool {
				return item.TableName == *d.DependencyTableName && item.TableName == *d.DependencyColumnName
			})
			rel := Relation{Schema: dependencySchema}
			graph.AllRelations[d.GetKey()] = rel
			schemaLink.Schema = dependencySchema
		}
		graph.RelationSchemaLinks = append(graph.RelationSchemaLinks, schemaLink)
		graph.RelationSchemaLinks = append(graph.RelationSchemaLinks)
	}
	return &graph
}

// TODO pk shuffle
func getColumnData(table *cfg.Table, schemas []*databaseUtils.Schema) ([]string, [][]any) {
	var columns []string

	var result [][]any
	for _, column := range schemas {
		columns = append(columns, column.ColumnName)
	}
	result = Generate(table, Settings.MaxBatch, schemas, result)
	return columns, result
}

func Generate(table *cfg.Table, count int, schemas []*databaseUtils.Schema, result [][]any) [][]any {
	for i := 0; i < count; i++ {
		var values []any
		for _, column := range schemas {
			values = append(values, databaseUtils.GetValue(column))
		}
		result = append(result, values)
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
