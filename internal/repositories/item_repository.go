package repositories

import (
	"errors"
	"server/internal/models"
	"strings"

	"gorm.io/gorm"
)

type ItemRepository interface {
	Create(item *models.Item) error //C
	CreateList(items *[]models.Item) error
	GetByID(id uint) (*models.Item, error) //R
	GetPagination(page int, isName bool, isAsc bool) (*[]models.Item, int64, error)
	SearchPagination(keyword string, page int, isName bool, isAsc bool) (*[]models.Item, int64, error)
	ChooseItem(keyword string) (*[]models.ItemSearch, error)
	Update(item *models.Item) error //U
	UpdateName(id uint, name string) error
	UpdatePrice(id uint, price int) error
	UpdateUnit(id uint, unit string) error
	UpdateDescription(id uint, description string) error
	Delete(id uint) error //D
	TestFunc(id uint) (*models.Item, error)
}
type itemRepository struct {
	db *gorm.DB
}

func NewItemRepository(db *gorm.DB) ItemRepository {
	return &itemRepository{db: db}
}
func (r *itemRepository) Create(item *models.Item) error {
	return r.db.Create(item).Error
}
func (r *itemRepository) GetByID(id uint) (*models.Item, error) {
	var item models.Item
	result := r.db.First(&item, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &item, nil
}

func (r *itemRepository) Update(item *models.Item) error {
	return r.db.Save(item).Error
}
func (r *itemRepository) Delete(id uint) error {
	return r.db.Delete(&models.Item{}, id).Error
}
func (r *itemRepository) UpdateName(id uint, name string) error {
	item := models.Item{}
	item.ID = id
	return r.db.Model(&item).Update("name", name).Error
}
func (r *itemRepository) UpdatePrice(id uint, price int) error {
	item := models.Item{}
	item.ID = id
	return r.db.Model(&item).Update("price", price).Error
}
func (r *itemRepository) UpdateUnit(id uint, unit string) error {
	item := models.Item{}
	item.ID = id
	return r.db.Model(&item).Update("unit", unit).Error
}
func (r *itemRepository) UpdateDescription(id uint, description string) error {
	// des := description
	item := models.Item{}
	item.ID = id
	return r.db.Model(&item).Update("description", description).Error
}
func (r *itemRepository) CreateList(items *[]models.Item) error {
	return r.db.Create(items).Error
}
func (r *itemRepository) GetPagination(page int, isName bool, isAsc bool) (*[]models.Item, int64, error) {
	var items []models.Item
	var count int64 = -1
	if page <= 1 {
		r.db.Table("items").Count(&count)
		if page == 0 {
			page = int(count/7) + 1
		}
	}
	offsetValue := 7 * (page - 1)
	orderStr := ""
	if isName && isAsc {
		orderStr = "name"
	} else if isName && !isAsc {
		orderStr = "name desc"
	} else if !isName && isAsc {
		orderStr = "price"
	} else {
		orderStr = "price desc"
	}
	result := r.db.Order(orderStr).Limit(7).Offset(offsetValue).Find(&items)
	if err := result.Error; err != nil {
		return nil, count, err
	}
	return &items, count, nil
}
func (r *itemRepository) SearchPagination(keyword string, page int, isName bool, isAsc bool) (*[]models.Item, int64, error) {
	if keyword == "" {
		return nil, -1, errors.New("no keyword")
	}
	var items []models.Item
	var count int64 = -1
	offsetValue := 7 * (page - 1)
	orderStr := ""
	if isName && isAsc {
		orderStr = "name"
	} else if isName && !isAsc {
		orderStr = "name desc"
	} else if !isName && isAsc {
		orderStr = "price"
	} else {
		orderStr = "price desc"
	}
	var result *gorm.DB
	var builder strings.Builder
	builder.WriteString("%")
	builder.WriteString(keyword)
	builder.WriteString("%")
	if page <= 1 {
		result = r.db.Table("items").Where("name LIKE ? OR unit LIKE ?", builder.String(), builder.String()).Count(&count)
		if page == 0 {
			offsetValue = int(count - count%7)
		}
		result = result.Order(orderStr).Limit(7).Offset(offsetValue).Find(&items)
	} else {
		result = r.db.Table("items").Where("name LIKE ? OR unit LIKE ?", builder.String(), builder.String()).Count(&count).Order(orderStr).Limit(7).Offset(offsetValue).Find(&items)
	}
	if result.Error != nil {
		return nil, count, result.Error
	}
	return &items, count, nil
}
func (r *itemRepository) ChooseItem(keyword string) (*[]models.ItemSearch, error) {
	if keyword == "" {
		return nil, errors.New("NOKEYWORD")
	}
	var itemSearch []models.ItemSearch
	var builder strings.Builder
	builder.WriteString("%")
	builder.WriteString(keyword)
	builder.WriteString("%")
	result := r.db.Model(&models.Item{}).Select("id, name, unit, price").Where("name LIKE ? OR unit LIKE ?", builder.String(), builder.String()).Order("name").Find(&itemSearch)
	if result.Error != nil {
		return nil, result.Error
	}
	return &itemSearch, nil
}
func (r *itemRepository) TestFunc(id uint) (*models.Item, error) {
	var item models.Item
	err := r.db.Preload("OrderItems.Order").First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}
