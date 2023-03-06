package postgres

import (
	"fmt"
	"time"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

const (
	maxOpenConnections = 60
	connMaxLifetime    = 120
	maxIdleConnections = 30
	connMaxIdleTime    = 20
)

var Connections = map[string]*ConnectionSet{}

type ConnectionSet struct {
	SqlxConn *sqlx.DB
	PgxConn  *pgx.Conn
}

func (c Config) toPgConnection() string {
	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		c.Host,
		c.Port,
		c.User,
		c.Dbname,
		c.Password,
		c.SSLMode,
	)
	return dataSourceName
}
func (c Config) ChangeDb(database string) Config {
	c.Dbname = database
	return c
}
func (c Config) GetConnection() (*ConnectionSet, error) {
	if data, ok := Connections[c.toPgConnection()]; ok {
		return data, nil
	}
	config := stdlib.DriverConfig{
		ConnConfig: pgx.ConnConfig{
			PreferSimpleProtocol: true,
		},
	}
	stdlib.RegisterDriverConfig(&config)
	connectionString := config.ConnectionString(c.toPgConnection())
	db, err := sqlx.Connect("pgx", connectionString)
	if err != nil {
		return nil, errors.Wrap(err, "Database.Connect")
	}
	pgxConnectionString, err := pgx.ParseConnectionString(c.toPgConnection())
	if err != nil {
		return nil, errors.Wrap(err, "pgx.ParseConnectionString(connectionString)")
	}
	pgxConn, err := pgx.Connect(pgxConnectionString)
	if err != nil {
		return nil, errors.Wrap(err, "pgx.Connect(pgxConnectionString)")
	}
	db.SetMaxOpenConns(maxOpenConnections)
	db.SetConnMaxLifetime(connMaxLifetime * time.Second)
	db.SetMaxIdleConns(maxIdleConnections)
	db.SetConnMaxIdleTime(connMaxIdleTime * time.Second)
	if err = db.Ping(); err != nil {
		return nil, err
	}
	connSet := &ConnectionSet{
		SqlxConn: db,
		PgxConn:  pgxConn,
	}
	Connections[c.toPgConnection()] = connSet
	return connSet, err
}
