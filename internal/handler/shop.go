package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/handmade/internal/domain"
	"github.com/mikalai2006/handmade/internal/middleware"
	"github.com/mikalai2006/handmade/internal/utils"
)


func (h *Handler) registerShop(router *gin.RouterGroup) {
		shops := router.Group("/shops")
		{
			shops.GET("/",  h.Find)
			shops.POST("/", middleware.SetUserIdentity, h.CreateShop)
		}
}


func (h *Handler) CreateShop(c *gin.Context) {
	userId, err := middleware.GetUserId(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	var input domain.Shop
	if err := c.BindJSON(&input); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	shop, err := h.services.Shop.CreateShop(userId, input)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
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
	params, err := utils.GetParamsFromRequest(c, domain.Shop{})
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)

		return
	}

	shops, err := h.services.Shop.GetAllShops(params)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusOK, shops)
}

type inputs struct {
	domain.RequestParams
	domain.Shop
}

// @Summary Find shops by params
// @Security ApiKeyAuth
// @Tags shop
// @Description Input params for search shops
// @ModuleID shops
// @Accept  json
// @Produce  json
// @Param input query inputs true "params for search shops"
// @Success 200 {object} []domain.Response
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/shops [get]
func (h *Handler) Find(c *gin.Context) {
	params, err := utils.GetParamsFromRequest(c, domain.Shop{})
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)

		return
	}

	shops, err := h.services.Shop.Find(params)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)

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