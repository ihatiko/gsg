package config

type Settings struct {
	DefaultSet  int
	Types       Types
	Databases   []*Database
	Connections []DatabaseConnection
}
type Types struct {
	VarCharDefaultLength int
}

type Table struct {
	Name string
	Set  int
}

type DatabaseConnection struct {
	Name       string
	Connection string
}

type Database struct {
	Name   string
	Tables []*Table
}

type Config struct {
	Settings *Settings
}
