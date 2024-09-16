package database

import (
	"fiber-project/models"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	err := godotenv.Load("../.env") // Load .env file
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dsn := os.Getenv("DB_DSN") // Get DSN from environment variables
	if dsn == "" {
		log.Fatal("DB_DSN is not set in .env file")
	}

	connection, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	DB = connection

	connection.AutoMigrate(&models.User{}, &models.PasswordReset{})
	fmt.Println("Database connection successful!")
}
