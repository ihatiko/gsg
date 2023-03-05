package config

import "gsg/postgres"

type Config struct {
	Settings *Settings
}

type Settings struct {
	DefaultSet          int
	DefaultTypeSettings Types
	Databases           []*Database
	Connections         []DatabaseConnection
	BlackListPath       []string
	Separator           string
}

type DatabaseConnection struct {
	Name       string
	Connection postgres.Config
}
