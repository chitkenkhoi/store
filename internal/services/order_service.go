package services

import (
	"errors"
	"server/internal/models"
	"server/internal/repositories"
	"server/internal/utils"
	"strconv"
)

type OrderService interface {
	Create(order *models.Order, items []models.OrderItemInput) error
	GetByID(id uint) (*models.Order, error)
	SearchWithFilter(keyword, sort string, page int, isAsc, isFilteredByPayed, isFilteredByDelivered, filterByPayed, filterByDelivered bool) (*[]models.OrderSearch, int64, error)
	Update(order *models.Order) error
	UpdateOneField(id uint, field string, content string) error
	UpdateItemInOrder(itemId uint, orderId uint, quantity uint, function int) error
	Delete(id uint) error
}
type orderService struct {
	itemRepo  repositories.ItemRepository
	orderRepo repositories.OrderRepository
}

func NewOrderService(itemRepo repositories.ItemRepository, orderRepo repositories.OrderRepository) OrderService {
	return &orderService{
		itemRepo:  itemRepo,
		orderRepo: orderRepo,
	}
}
func (s *orderService) Create(order *models.Order, items []models.OrderItemInput) error {
	order.ID = 0
	if newBuyer, err := utils.TextValidateProcess(order.Buyer); err != nil {
		return errors.New("buyer is not valid")
	} else {
		order.Buyer = newBuyer
	}
	return s.orderRepo.Create(order, items)
}
func (s *orderService) GetByID(id uint) (*models.Order, error) {
	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	price := 0
	for _, item := range order.OrderItems {
		price += item.Item.Price * int(item.Quantity)
		order.OrderItemInput = append(order.OrderItemInput, models.OrderItemResponse{
			ItemID:   item.ItemID,
			Quantity: item.Quantity,
			Price:    item.Item.Price,
		})
	}
	order.Price = price
	order.OrderItems = []models.OrderItem{}
	return order, nil
}
func (s *orderService) SearchWithFilter(keyword, sort string, page int, isAsc, isFilteredByPayed, isFilteredByDelivered, filterByPayed, filterByDelivered bool) (*[]models.OrderSearch, int64, error) {
	if page < 0 {
		page = 0
	}
	if newKeyword, err := utils.TextValidateProcess(keyword); err != nil {
		return nil, -1, errors.New("keyword is not valid")
	} else {
		switch sort {
		case "buyer", "created_at":
			return s.orderRepo.SearchWithFilter(newKeyword, sort, page, isAsc, isFilteredByPayed, isFilteredByDelivered, filterByPayed, filterByDelivered)
		default:
			return nil, -1, errors.New("input is invalid")
		}
	}
}
func (s *orderService) Update(order *models.Order) error {
	if _, err := s.orderRepo.GetByIdWithNoItems(order.ID); err != nil {
		return err
	}
	if newBuyer, err := utils.TextValidateProcess(order.Buyer); err != nil {
		return errors.New("buyer is not valid")
	} else {
		order.Buyer = newBuyer
	}

	return s.orderRepo.Update(order)
}

func (s *orderService) UpdateOneField(id uint, field string, content string) error {
	if _, err := s.orderRepo.GetByIdWithNoItems(id); err != nil {
		return err
	}
	switch field {
	case "buyer":
		if newBuyer, err := utils.TextValidateProcess(content); err != nil {
			return errors.New("buyer is not valid")
		} else {
			return s.orderRepo.UpdateBuyer(id, newBuyer)
		}
	case "discount":
		if discount, err := strconv.Atoi(content); err != nil {
			return errors.New("discount is not valid")
		} else {
			return s.orderRepo.UpdateDiscount(id, discount)
		}
	case "delivery":
		if status, err := strconv.ParseBool(content); err != nil {
			return errors.New("delivery is not valid")
		} else {
			return s.orderRepo.UpdateDeliveryStatus(id, status)
		}
	case "payment":
		if status, err := strconv.ParseBool(content); err != nil {
			return errors.New("payment is not valid")
		} else {
			return s.orderRepo.UpdatePaymentStatus(id, status)
		}
	default:
		return errors.New("input is not valid")
	}

}
func (s *orderService) UpdateItemInOrder(itemId uint, orderId uint, quantity uint, function int) error {
	if err := s.isOrderItemIn(itemId, orderId); err != nil {
		return err
	}
	if function == 0 { // Update new quantity
		return s.orderRepo.ModifyQuantityOfItem(orderId, itemId, quantity)
	} else if function == 1 { //Add an item to an order
		return s.orderRepo.AddItemToOrder(orderId, itemId, quantity)
	} else if function == -1 { //Remove an item out of an order
		return s.orderRepo.RemoveItemFromOrder(orderId, itemId)
	} else {
		return errors.New("function is not valid")
	}
}
func (s *orderService) isOrderItemIn(itemId, orderId uint) error {
	if _, er := s.itemRepo.GetByID(itemId); er != nil {
		return errors.New("item not found")
	}
	if _, er := s.orderRepo.GetByIdWithNoItems(orderId); er != nil {
		return errors.New("order not found")
	}
	return nil
}
func (s *orderService) Delete(id uint) error {
	return s.orderRepo.Delete(id)
}
