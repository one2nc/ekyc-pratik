package db

import (

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func InitiateDB() (*gorm.DB, error) {
	dbConfig := newDbConfig()
	err := dbMigrationUP(dbConfig)
	if err != nil {
		return nil, err
	}


	db, err := gorm.Open(postgres.Open(dbConfig.connString), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "ekyc_schema.", // schema name
			SingularTable: false,
		}})

	if err != nil {
		return nil, err
	}
	return db, nil
}
