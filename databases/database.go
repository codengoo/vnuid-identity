package databases

import (
	"log"
	"os"

	// "vnuid-identity/entities"

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

	// Auto sync at management service
	// db.AutoMigrate(&entities.Profile{}, &entities.User{}, entities.Session{}, entities.NFC{})
	DB = db
}
