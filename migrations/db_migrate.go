package migrations

import (
	"fmt"
	"log"
	"server/internal/models"

	"gorm.io/gorm"
)

func DB_migrate(db *gorm.DB) {
	if err := db.AutoMigrate(&models.Item{}, &models.Order{}, &models.OrderItem{}); err != nil {
		log.Fatal("Can not migrate the database")
	}
	fmt.Println("Migrate the database successfully")
}
