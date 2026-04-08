package database

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(dbURI string, dbName string) error {

	uri := fmt.Sprintf("%s/%s", dbURI, dbName)

	m, err := migrate.New(
		"file://migrations",
		uri,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
