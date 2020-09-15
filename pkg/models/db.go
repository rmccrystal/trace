package models

import "gorm.io/gorm"

// DatabaseCofig configures the connection to a database
type DatabaseConfig struct {

}

// A Database manages all of the models
type Database struct {
	gorm.DB
}

// NewDatabase connects to the database using a DatabaseConfig
func NewDatabase(config DatabaseConfig) *Database, error {
	db, err := gorm.Open("test.db", &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &Database{DB: db}
}