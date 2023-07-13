package db

import (
	"fmt"
	"os"
)

type dbConfig struct {
	dbHost, dbPassword, dbUser, dbName, dbPort string
	connString                                 string
	migrationFile                              string
}

func newDbConfig() dbConfig {
	dbHost, dbPassword, dbUser, dbName, dbPort, dbMigrationFile :=
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_MIGRATION_FILE")
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName)
	migrationFile := fmt.Sprintf("file://%s",
		dbMigrationFile)
	dbConfig := dbConfig{
		dbHost:        dbHost,
		dbPassword:    dbPassword,
		dbUser:        dbUser,
		dbName:        dbName,
		dbPort:        dbPort,
		connString:    connectionString,
		migrationFile: migrationFile,
	}
	return dbConfig
}
