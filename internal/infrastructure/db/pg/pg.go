// Package pg provides the implementation of the database layer using PostgreSQL (sqlx).
package pg

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-migrate/migrate/v4"
	postgresql "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewDB(
	dbName string,
	host string,
	port int,
	user string,
	password string,
	ssl string,
) (*sqlx.DB, error) {
	connection := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		host,
		port,
		user,
		password,
		dbName,
		ssl)

	return sqlx.Open("postgres", connection)
}

func MigrateUp(db *sql.DB) error {
	timeout := 1 * time.Second
	driver, err := postgresql.WithInstance(db, &postgresql.Config{StatementTimeout: timeout})
	if err != nil {
		return err
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	path := "file://" + filepath.Join(
		wd,
		"internal/infrastructure/db/migrations/postgres",
	)

	m, err := migrate.NewWithDatabaseInstance(path, "postgres", driver)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
