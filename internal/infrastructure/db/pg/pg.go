// Package pg provides the implementation of the database layer using PostgreSQL (sqlx).
package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-migrate/migrate/v4"
	postgresql "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	ErrInvalidArgument = errors.New("invalid argument provided")
	ErrUsrDoesNotExist = errors.New("usr was not found")
)

type closer interface {
	Close() error
}

func safeClose(ctx context.Context, c closer) {
	err := c.Close()
	if err != nil {
		slog.WarnContext(ctx, "failed to close resource", "err", err)
	}
}

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

	db, err := sqlx.Open("postgres", connection)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
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
