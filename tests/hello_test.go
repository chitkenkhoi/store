package tests

import (
	// "bytes"
	"fmt"
	"server/internal/repositories"
	"server/internal/services"
	"server/internal/utils"
	"server/migrations"
	"testing"
)

//	type Order struct {
//	    gorm.Model
//	    Buyer      string      `json:"buyer" gorm:"column:buyer"`
//	    OrderItems []OrderItem `json:"order_items"`
//	    Date       *time.Time  `json:"date" gorm:"column:date"`
//	    Discount   int         `json:"discount" gorm:"column:discount"`
//	    Payed      bool        `json:"payed" gorm:"column:payed"`
//	    Delivered  bool        `json:"delivered" gorm:"column:delivered"`
//	}
func Test(t *testing.T) {
	db := utils.DBconnector()
	migrations.DB_migrate(db)
	itemRepo := repositories.NewItemRepository(db)
	orderRepo := repositories.NewOrderRepository(db)
	orderService := services.NewOrderService(itemRepo, orderRepo)
	fmt.Println(orderService.GetByID(2))
	// order := models.Order{
	// 	Buyer: "Chú Hải",
	// 	Payed: true,
	// }
	// orderItemInput := []models.OrderItemInput{
	// 	{
	// 		ItemID:   1,
	// 		Quantity: 2,
	// 	},
	// 	{
	// 		ItemID:   4,
	// 		Quantity: 1,
	// 	},
	// }
	// keyword, sort string, page int, isAsc, isFilteredByPayed, isFilteredByDelivered, filterByPayed, filterByDelivered bool
	// order := &models.Order{
	// 	Model: gorm.Model{
	// 		ID: 1,
	// 	},
	// 	Discount: 5,
	// }4
	order, _ := orderRepo.GetByID(2)
	fmt.Println(order.OrderItems[2].Item.DeletedAt)
	// orderTest, _ := orderRepo.GetByIDTest(2)
	// fmt.Println(orderTest.OrderItems[3])
}
