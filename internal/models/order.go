package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	Buyer          string              `json:"buyer" gorm:"column:buyer;index"`
	Date           *time.Time          `json:"date" gorm:"column:date"`
	Discount       int                 `json:"discount" gorm:"column:discount;default:0" `
	Payed          bool                `json:"payed" gorm:"column:payed;default:false"`
	Delivered      bool                `json:"delivered" gorm:"column:delivered;default:false"`
	Price          int                 `json:"price" gorm:"-"`
	OrderItems     []OrderItem         `json:"-" gorm:"foreignKey:OrderID"`
	OrderItemInput []OrderItemResponse `json:"itemList" gorm:"-"`
}
type OrderItem struct {
	gorm.Model
	OrderID  uint  `json:"order_id" gorm:"column:order_id"`
	ItemID   uint  `json:"item_id" gorm:"column:item_id"`
	Quantity uint  `json:"quantity" gorm:"column:quantity"`
	Item     Item  `json:"item" gorm:"foreignKey:ItemID"`
	Order    Order `json:"order" gorm:"foreignKey:OrderID;OnDelete:CASCADE"`
}
type OrderItemInput struct {
	ItemID   uint `json:"item_id"`
	Quantity uint `json:"quantity"`
}
type OrderItemResponse struct {
	ItemID   uint `json:"item_id"`
	Quantity uint `json:"quantity"`
	Price    int  `json:"price"`
}
type OrderSearch struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	Buyer     string    `json:"buyer" gorm:"column:buyer"`
	Date      time.Time `json:"date" gorm:"column:created_at"`
	Payed     bool      `json:"payed" gorm:"column:payed"`
	Delivered bool      `json:"delivered" gorm:"column:delivered"`
}
