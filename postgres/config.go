package postgres

type Config struct {
	Host           string
	Port           string
	User           string
	Password       string
	Dbname         string
	SSLMode        string
	PgDriver       string
	CreateDb       bool
	Schema         string
	InitMigrations bool
}
