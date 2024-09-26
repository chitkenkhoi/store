package handlers

import (
	"net/http"
	"server/internal/models"
	"server/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ItemHandler interface {
	Create(c *gin.Context)
	CreateList(c *gin.Context)
	GetByID(c *gin.Context)
}
type itemHandler struct {
	itemService services.ItemService
}

func NewItemHandler(itemService services.ItemService) ItemHandler {
	return &itemHandler{
		itemService: itemService,
	}
}

func (h *itemHandler) Create(c *gin.Context) {
	var item models.Item
	if err := c.ShouldBind(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.itemService.Create(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"message": "item created successfully",
		})
	}
}
func (h *itemHandler) CreateList(c *gin.Context) {
	var items []models.Item
	if err := c.ShouldBind(&items); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if err := h.itemService.CreateList(&items); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "item list created successfully",
	})
}
func (h *itemHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if item, er := h.itemService.GetByID(uint(id)); er != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": er.Error(),
		})
		return
	} else {
		c.JSON(200, gin.H{
			"data": item,
		})
	}
}
