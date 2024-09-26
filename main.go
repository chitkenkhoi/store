package main

import (
	"server/internal/handlers"
	"server/internal/repositories"
	"server/internal/services"
	"server/internal/utils"
	"server/migrations"

	"github.com/gin-gonic/gin"
)

func main() {
	db := utils.DBconnector()
	migrations.DB_migrate(db)
	itemRepo := repositories.NewItemRepository(db)
	orderRepo := repositories.NewOrderRepository(db)

	orderService := services.NewOrderService(itemRepo, orderRepo)
	itemService := services.NewItemService(itemRepo)

	itemHandler := handlers.NewItemHandler(itemService)
	orderHandler := handlers.NewOrderHandler(orderService)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/item/:id", itemHandler.GetByID)
	r.POST("/item/create", itemHandler.Create)
	r.POST("/item/createList", itemHandler.CreateList)
	r.GET("/order/:id", orderHandler.GetByID)
	r.Run(":8080")
}
