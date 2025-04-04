package databases

import (
	"log"
	"os"

	"vnuid-identity/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {}

var DB *gorm.DB

func ConnectDB() {
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Could not open database")
	}

	db.AutoMigrate(&models.User{})
	DB = db
}
