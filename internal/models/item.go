package models

import "gorm.io/gorm"

type Item struct {
	gorm.Model
	Price       int         `json:"price" gorm:"column:price"`
	Name        string      `json:"name" gorm:"column:name"`
	Unit        string      `json:"unit" gorm:"column:unit"`
	Description *string     `json:"description" gorm:"column:description"`
	OrderItems  []OrderItem `json:"-" gorm:"foreignKey:ItemID"`
}
type ItemSearch struct {
	ID    uint   `json:"id" gorm:"primarykey"`
	Name  string `json:"name" gorm:"name"`
	Unit  string `json:"unit" gorm:"unit"`
	Price int    `json:"price" gorm:"price"`
}
type Test struct {
	Test string
}
