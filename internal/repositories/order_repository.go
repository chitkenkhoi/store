package repositories

import (
	"errors"
	"fmt"
	"server/internal/models"
	"strings"

	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(order *models.Order, items []models.OrderItemInput) error // C
	GetByID(id uint) (*models.Order, error)                          //R
	GetByIdWithNoItems(id uint) (*models.Order, error)
	SearchWithFilter(keyword, sort string, page int, isAsc, isFilteredByPayed, isFilteredByDelivered, filterByPayed, filterByDelivered bool) (*[]models.OrderSearch, int64, error)
	Update(order *models.Order) error // U
	UpdateBuyer(id uint, buyer string) error
	UpdateDiscount(id uint, discount int) error
	UpdateDeliveryStatus(id uint, status bool) error
	UpdatePaymentStatus(id uint, status bool) error
	AddItemToOrder(orderID, itemID uint, quantity uint) error
	ModifyQuantityOfItem(orderID, itemID uint, newQuantity uint) error
	RemoveItemFromOrder(orderID, itemID uint) error
	Delete(id uint) error //D
	TestFunc() error
}
type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}
func (r *orderRepository) TestFunc() error {
	var test models.Test
	r.db.Select("name; drop table test;").First(&test)
	fmt.Println(test)
	return nil
}
func (r *orderRepository) Create(order *models.Order, items []models.OrderItemInput) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()
		return err
	}
	for _, item := range items {
		orderItem := models.OrderItem{
			OrderID:  order.ID,
			ItemID:   item.ItemID,
			Quantity: item.Quantity,
		}
		if err := tx.Create(&orderItem).Error; err != nil {
			tx.Rollback()
			return errors.New("item not found")
		}
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
func (r *orderRepository) GetByIdWithNoItems(id uint) (*models.Order, error) {
	var order models.Order
	if err := r.db.First(&order, id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}
func (r *orderRepository) GetByID(id uint) (*models.Order, error) {
	var order models.Order

	// First, query the order without preloading, respecting soft delete
	err := r.db.First(&order, id).Error
	if err != nil {
		return nil, err
	}

	// Then, preload the OrderItems and Items, ignoring soft delete for Items
	err = r.db.Model(&order).
		Preload("OrderItems", func(db *gorm.DB) *gorm.DB {
			return db // This unscopes the join table if it has soft delete
		}).
		Preload("OrderItems.Item", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped() // This unscopes the Items
		}).
		First(&order, id).Error

	if err != nil {
		return nil, err
	}

	return &order, nil
}
func (r *orderRepository) SearchWithFilter(keyword, sort string, page int, isAsc, isFilteredByPayed, isFilteredByDelivered, filterByPayed, filterByDelivered bool) (*[]models.OrderSearch, int64, error) {
	var count int64 = -1
	var orderSearchs []models.OrderSearch
	offsetValue := 7 * (page - 1)
	var result *gorm.DB
	var builder strings.Builder
	// builder.WriteString("%")
	// builder.WriteString(keyword)
	// builder.WriteString("%")

	builder.WriteString(sort)
	if !isAsc {
		builder.WriteString(" desc")
	}
	orderString := builder.String()
	builder.Reset()
	if keyword == "" {
		if isFilteredByPayed && isFilteredByDelivered {
			result = r.db.Table("orders").Where("payed = ? AND delivered = ?", filterByPayed, filterByDelivered).Count(&count)
		} else if isFilteredByPayed && !isFilteredByDelivered {
			result = r.db.Table("orders").Where("payed = ?", filterByPayed).Count(&count)
		} else if !isFilteredByPayed && isFilteredByDelivered {
			result = r.db.Table("orders").Where("delivered = ?", filterByDelivered).Count(&count)
		} else {
			result = r.db.Table("orders").Count(&count)
		}
		if page == 0 {
			offsetValue = int(count - count%7)
		}
		result = result.Order(orderString).Limit(7).Offset(offsetValue).Find(&orderSearchs)

	} else {
		builder.WriteString("%")
		builder.WriteString(keyword)
		builder.WriteString("%")
		if isFilteredByPayed && isFilteredByDelivered {
			result = r.db.Table("orders").Where("payed = ? AND delivered = ?", filterByPayed, filterByDelivered).Where("buyer LIKE ?", builder.String()).Count(&count)
		} else if isFilteredByPayed && !isFilteredByDelivered {
			result = r.db.Table("orders").Where("payed = ?", filterByPayed).Where("buyer LIKE ?", builder.String()).Count(&count)
		} else if !isFilteredByPayed && isFilteredByDelivered {
			result = r.db.Table("orders").Where("delivered = ?", filterByDelivered).Where("buyer LIKE ?", builder.String()).Count(&count)
		} else {
			result = r.db.Table("orders").Where("buyer LIKE ?", builder.String()).Count(&count)
		}
		if page == 0 {
			offsetValue = int(count - count%7)
		}
		result = result.Order(orderString).Limit(7).Offset(offsetValue).Find(&orderSearchs)
	}
	if err := result.Error; err != nil {
		return nil, -1, err
	}
	return &orderSearchs, count, nil
}
func (r *orderRepository) Update(order *models.Order) error {
	return r.db.Save(order).Error
}
func (r *orderRepository) UpdateBuyer(id uint, buyer string) error {
	return r.db.Model(&models.Order{Model: gorm.Model{ID: id}}).Update("buyer", buyer).Error
}
func (r *orderRepository) UpdateDiscount(id uint, discount int) error {
	return r.db.Model(&models.Order{Model: gorm.Model{ID: id}}).Update("discount", discount).Error
}
func (r *orderRepository) UpdatePaymentStatus(id uint, status bool) error {
	return r.db.Model(&models.Order{Model: gorm.Model{ID: id}}).Update("payed", status).Error
}
func (r *orderRepository) UpdateDeliveryStatus(id uint, status bool) error {
	return r.db.Model(&models.Order{Model: gorm.Model{ID: id}}).Update("delivered", status).Error
}
func (r *orderRepository) AddItemToOrder(orderID, itemID uint, quantity uint) error {
	var existingOrderItem models.OrderItem
	result := r.db.Where("order_id = ? AND item_id = ?", orderID, itemID).First(&existingOrderItem)
	if result.Error == nil {
		// Item already exists in the order
		return errors.New("item already in order")
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// An error occurred that wasn't "record not found"
		return result.Error
	}
	newOrderItem := models.OrderItem{
		OrderID:  orderID,
		ItemID:   itemID,
		Quantity: quantity,
	}
	return r.db.Create(&newOrderItem).Error
}
func (r *orderRepository) ModifyQuantityOfItem(orderID, itemID uint, newQuantity uint) error {
	result := r.db.Model(&models.OrderItem{}).
		Where("order_id = ? AND item_id = ?", orderID, itemID).
		Update("quantity", newQuantity)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("item not found in the order")
	}
	return nil
}
func (r *orderRepository) RemoveItemFromOrder(orderID, itemID uint) error {
	result := r.db.Where("order_id = ? AND item_id = ?", orderID, itemID).Delete(&models.OrderItem{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("item not found in the order")
	}
	return nil
}

//	type Order struct {
//		gorm.Model
//		Buyer     string     `json:"buyer" gorm:"buyer"`
//		Items     []Item     `json:"items" gorm:"many2many:order_items;"`
//		Date      *time.Time `json:"date" gorm:"date"`
//		Discount  int        `json:"discount" gorm:"discount"`
//		Payed     bool       `json:"payed" gorm:"payed"`
//		Delivered bool       `json:"delivered" gorm:"delivered"`
//	}
func (r *orderRepository) Delete(id uint) error {
	return r.db.Delete(&models.Order{}, id).Error
}
