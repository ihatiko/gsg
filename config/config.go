package config

type Settings struct {
	Bulk       bool
	DefaultSet int
	MaxBatch   int
	Databases  []*Database
}

type Table struct {
	Name string
	Set  int
}

type Database struct {
	Name   string
	Tables []*Table
}

type Config struct {
	Settings *Settings
}
