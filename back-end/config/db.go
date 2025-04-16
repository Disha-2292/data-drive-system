package config

import (
	"fmt"
	"log"
	"os"

	"github.com/Disha-2292/data-drive-system/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	// Automigrate models
	err = db.AutoMigrate(&models.User{}, &models.File{})
	if err != nil {
		log.Fatal("Failed to automigrate models: ", err)
	}

	DB = db
	fmt.Println("✅ Database connected and migrated")
}
