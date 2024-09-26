package services

import (
	"errors"
	"server/internal/models"
	"server/internal/repositories"
	"server/internal/utils"
	"strconv"
)

type ItemService interface {
	Create(item *models.Item) error
	CreateList(items *[]models.Item) error
	GetByID(id uint) (*models.Item, error)
	SearchPagination(keyword string, page int, isName, isAsc bool) (*[]models.Item, int64, error)
	ChooseItem(keyword string) (*[]models.ItemSearch, error)
	Update(item *models.Item) error
	UpdateOneField(id uint, field string, content string) error
	Delete(id uint) error
}
type itemService struct {
	itemRepo repositories.ItemRepository
}

func NewItemService(itemRepo repositories.ItemRepository) ItemService {
	return &itemService{
		itemRepo: itemRepo,
	}
}
func itemValidate(item *models.Item) error {
	item.ID = 0
	if item.Price <= 0 {
		return errors.New("item price can not be negative")
	}
	if newName, err := utils.TextValidateProcess(item.Name); err != nil {
		return errors.New("name is not valid")
	} else {
		item.Name = newName
	}
	if newUnit, err := utils.TextValidateProcess(item.Unit); err != nil {
		return errors.New("unit is not valid")
	} else {
		item.Unit = newUnit
	}
	return nil
}
func (s *itemService) Create(item *models.Item) error {
	if err := itemValidate(item); err != nil {
		return err
	}

	return s.itemRepo.Create(item)
}
func (s *itemService) CreateList(items *[]models.Item) error {
	for index := range *items {
		if err := itemValidate(&((*items)[index])); err != nil {
			return errors.New(err.Error() + strconv.Itoa(index))
		}
	}
	return s.itemRepo.CreateList(items)
}
func (s *itemService) GetByID(id uint) (*models.Item, error) {
	return s.itemRepo.GetByID(id)
}
func (s *itemService) SearchPagination(keyword string, page int, isName, isAsc bool) (*[]models.Item, int64, error) {
	if page < 0 {
		page = 0
	}
	if keyword == "" {
		return s.itemRepo.GetPagination(page, isName, isAsc)
	} else {
		newKeyword, err := utils.TextValidateProcess(keyword)
		if err != nil {
			return nil, -1, err
		}
		return s.itemRepo.SearchPagination(newKeyword, page, isName, isAsc)
	}
}
func (s *itemService) ChooseItem(keyword string) (*[]models.ItemSearch, error) {
	if newKeyword, err := utils.TextValidateProcess(keyword); err != nil {
		return nil, err
	} else {
		return s.itemRepo.ChooseItem(newKeyword)
	}
}
func (s *itemService) Update(item *models.Item) error {
	if _, err := s.itemRepo.GetByID(item.ID); err != nil {
		return err
	}
	id := item.ID
	if er := itemValidate(item); er != nil {
		return er
	}
	item.ID = id
	return s.itemRepo.Update(item)
}
func (s *itemService) UpdateOneField(id uint, field string, content string) error {
	switch field {
	case "name":
		if newName, err := utils.TextValidateProcess(content); err != nil {
			return err
		} else {
			return s.itemRepo.UpdateName(id, newName)
		}
	case "unit":
		if newUnit, err := utils.TextValidateProcess(content); err != nil {
			return err
		} else {
			return s.itemRepo.UpdateUnit(id, newUnit)
		}
	case "price":
		if price, err := strconv.Atoi(content); err != nil {
			return errors.New("price is not integer")
		} else {
			return s.itemRepo.UpdatePrice(id, price)
		}
	case "description":
		return s.itemRepo.UpdateDescription(id, content)
	default:
		return errors.New("field is invalid")
	}
}
func (s *itemService) Delete(id uint) error {
	return s.itemRepo.Delete(id)
}
