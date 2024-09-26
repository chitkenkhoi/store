package utils

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DBconnector() *gorm.DB {
	dsn := "host=localhost user=myuser password=mypassword dbname=mydatabase port=5431 sslmode=disable TimeZone=UTC"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Check the connection
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get database:", err)
	}

	err = sqlDB.Ping()
	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	fmt.Println("Successfully connected to the database!")
	return db
}
