package handlers

import (
	"server/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrderHandler interface {
	GetByID(c *gin.Context)
}
type orderHandler struct {
	orderService services.OrderService
}

func NewOrderHandler(orderService services.OrderService) OrderHandler {
	return &orderHandler{
		orderService: orderService,
	}
}
func (h *orderHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	newid, _ := strconv.Atoi(id)
	order, err := h.orderService.GetByID(uint(newid))
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"data": order,
		})
	}
}
