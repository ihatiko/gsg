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
}

type DatabaseConnection struct {
	Name       string
	Connection *postgres.Config
}
