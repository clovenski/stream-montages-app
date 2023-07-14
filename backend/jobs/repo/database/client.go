package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/clovenski/stream-montages-app/backend/jobs/repo/config"
	"github.com/clovenski/stream-montages-app/backend/jobs/repo/models"
)

var DB *gorm.DB

func Connect(config *config.Config) {
	var err error
	cnxStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", config.DBHost, config.DBUserName, config.DBUserPassword, config.DBName, config.DBPort)

	DB, err = gorm.Open(postgres.Open(cnxStr), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database")
	}
	log.Println("Successfully connected to the Database")
}

func Migrate() {
	DB.AutoMigrate(&models.MontageJobDefinition{})
	log.Println("Successfully migrated database")
}
