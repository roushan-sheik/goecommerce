package config

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase(cfg *Config) *gorm.DB {

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  cfg.Dsn,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connecting to the database")
	}

	return db

}
