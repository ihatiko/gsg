package generator

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
	Name    string
	Columns map[string]*Column
}
type Constraint struct {
	DependencyTableName    string
	DependencyColumnName   string
	DependencyDatabaseName string
}
type Column struct {
	GeneratedData []any
	Schema        *Schema
	Constraints   map[string]Constraint
}

type Enum struct {
	Value string `db:"value"`
}
