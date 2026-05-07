package db

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigration(dbURL string) error {

	m, err := migrate.New(
		"file://internal/infrastructure/db/migrations",
		dbURL,
	)

	if err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			return nil
		}
		return err
	}

	log.Println("Migration applied successfully")

	return nil

}
