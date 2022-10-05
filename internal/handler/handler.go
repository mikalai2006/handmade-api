package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/handmade/internal/service"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware

	_ "github.com/mikalai2006/handmade/docs"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler  {
	return &Handler{
		services: services,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	// add swagger route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// store := cookie.NewStore([]byte(os.Getenv("secret")))
	// router.Use(sessions.Sessions("mysession", store))

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.SignUp)
		auth.POST("/sign-in", h.SignIn)
		auth.POST("/logout", h.Logout)
	}

	oauth := router.Group("/oauth")
	{
		oauth.GET("/vk", h.OAuthVK)
		oauth.GET("/google", h.OAuthGoogle)
	}
	router.GET("/me", h.Me)
	router.GET("/googleme", h.MeGoogle)
	router.GET("/ws", h.wsEndPoint)

	api := router.Group("/api")
	{
		lists := api.Group("/lists")
		{
			lists.POST("/", h.CreateList)
			lists.GET("/", h.GetAllLists)
			lists.GET("/:id", h.GetListById)
			lists.PUT("/:id", h.UpdateList)
			lists.DELETE("/:id", h.DeleteList)
			items := api.Group(":id/items")
			{
				items.POST("/", h.CreateItem)
				items.GET("/", h.GetAllItems)
				items.GET("/:item_id", h.GetItemById)
				items.PUT("/:item_id", h.UpdateItem)
				items.DELETE("/:item_id", h.DeleteItem)
			}
		}
		shops := api.Group("/shops")
		{
			shops.GET("/", h.userIdentity, h.GetAllShops)
			shops.POST("/", h.userIdentity, h.CreateShop)
		}
	}
	return router
}