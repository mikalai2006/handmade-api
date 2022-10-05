package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/handmade/internal/domain"
)

func (h *Handler) CreateShop(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	var input domain.Shop
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	shop, err := h.services.Shop.CreateShop(userId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, shop)
}


// @Summary Shop Get all shops
// @Security ApiKeyAuth
// @Tags shop
// @Description get all shops
// @ModuleID shops
// @Accept  json
// @Produce  json
// @Success 200 {object} []domain.Shop
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/shops [get]
func (h *Handler) GetAllShops(c *gin.Context) {
	shops, err := h.services.Shop.GetAllShops()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, shops)
}

func (h *Handler) GetShopById(c *gin.Context) {

}

func (h *Handler) UpdateShop(c *gin.Context) {

}

func (h *Handler) DeleteShop(c *gin.Context) {

}