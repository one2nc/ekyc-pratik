package db

import (
	"log"

	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func dbMigrationUP(dbConfig dbConfig) error {
	// Set up database connection
	db, err := sql.Open("postgres", dbConfig.connString)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		return err
	}
	defer db.Close()

	// Set up migration driver
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Failed to set up migration driver: %v", err)
		return err
	}

	// Set up migration source
	m, err := migrate.NewWithDatabaseInstance(dbConfig.migrationFile, dbConfig.dbName, driver)
	if err != nil {
		return err
	}
	defer m.Close()

	// Apply "up" migration
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
