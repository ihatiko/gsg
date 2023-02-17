package config

type Settings struct {
	DefaultSet int
	Types      Types
	Databases  []*Database
}
type Types struct {
	VarCharDefaultLength int
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
