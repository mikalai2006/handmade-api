package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/handmade/internal/domain"
	"github.com/mikalai2006/handmade/internal/middleware"
	"github.com/mikalai2006/handmade/internal/utils"
)

func (h *Handler) RegisterUser(router *gin.RouterGroup) {
		user := router.Group("/user")
		{
			user.POST("/", middleware.SetUserIdentity, h.CreateUser)
			user.DELETE("/:id", middleware.SetUserIdentity, h.DeleteUser)
			user.PATCH("/:id", middleware.SetUserIdentity, h.UpdateUser)
			user.GET("/:id", h.GetUser)
			user.GET("/find", h.FindUser)
		}
}

// @Summary Get user by Id
// @Tags user
// @Description get user info
// @ModuleID user
// @Accept  json
// @Produce  json
// @Param id path string true "user id"
// @Success 200 {object} domain.User
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/user/{id} [get]
func (h *Handler) GetUser(c *gin.Context) {
	id := c.Param("id")

	user, err := h.services.User.GetUser(id)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, user)
}


type InputUser struct {
	domain.RequestParams
	domain.User
}

// @Summary Find few users
// @Security ApiKeyAuth
// @Tags user
// @Description Input params for search users
// @ModuleID user
// @Accept  json
// @Produce  json
// @Param input query InputUser true "params for search users"
// @Success 200 {object} []domain.Response[domain.User]
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/user [get]
func (h *Handler) FindUser(c *gin.Context) {
	params, err := utils.GetParamsFromRequest(c, domain.User{})
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	users, err := h.services.User.FindUser(params)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *Handler) CreateUser(c *gin.Context) {
	userId, err := middleware.GetUserId(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	var input domain.User
	if err := c.BindJSON(&input); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	user, err := h.services.User.CreateUser(userId, input)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		// utils.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary Delete user
// @Security ApiKeyAuth
// @Tags user
// @Description Delete user
// @ModuleID user
// @Accept  json
// @Produce  json
// @Param id path string true "user id"
// @Success 200 {object} []domain.Response
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/user/{id} [delete]
func (h *Handler) DeleteUser(c *gin.Context) {

	id := c.Param("id")

	// var input domain.User
	// if err := c.BindJSON(&input); err != nil {
	// 	newErrorResponse(c, http.StatusInternalServerError, err.Error())

	// 	return
	// }

	user, err := h.services.User.DeleteUser(id) // , input
	if err != nil {
		// tils.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		c.AbortWithError(http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusOK, user)
}


// @Summary Update user
// @Security ApiKeyAuth
// @Tags user
// @Description Update user
// @ModuleID user
// @Accept  json
// @Produce  json
// @Param id path string true "user id"
// @Param input body domain.User true "body for update user"
// @Success 200 {object} []domain.Response
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/user/{id} [put]
func (h *Handler) UpdateUser(c *gin.Context)  {

	id := c.Param("id")

	var input domain.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithError(http.StatusBadRequest, err) // .SetMeta(gin.H{"hello": "World"})

		return
	}

	user, err := h.services.User.UpdateUser(id, input)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusOK, user)
}