package generator

import (
	"github.com/jmoiron/sqlx"
	cfg "gsg/config"
)

type Schema struct {
	Database             string  `db:"table_catalog"`
	TableName            string  `db:"table_name"`
	ColumnName           string  `db:"column_name"`
	ColumnDefault        *string `db:"column_default"`
	DataType             string  `db:"data_type"`
	Length               *string `db:"character_maximum_length"`
	IsNullable           bool    `db:"is_nullable"`
	DependencyTableName  *string `db:"dependency_table_name"`
	DependencyColumnName *string `db:"dependency_column_name"`
	ConstraintType       *string `db:"constraint_type"`
}

type Table struct {
	Inserted bool
	Columns  map[string]*Column
}
type Constraint struct {
	DependencyTableName  string `db:"dependency_table_name"`
	DependencyColumnName string `db:"dependency_column_name"`
}
type Column struct {
	GeneratedData []any
	Schema        *Schema
	Constraints   map[string]Constraint
}

type Generator struct {
	Db       *sqlx.DB
	Settings *cfg.Settings
	Enums    map[string]map[string][]string
}

type Enum struct {
	Type  string `db:"type"`
	Value string `db:"value"`
}
