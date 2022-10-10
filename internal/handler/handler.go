package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/handmade/docs"
	"github.com/mikalai2006/handmade/internal/config"
	"github.com/mikalai2006/handmade/internal/service"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware

	_ "github.com/mikalai2006/handmade/docs"
)

type Handler struct {
	services *service.Services
	oauth config.OauthConfig
}

func NewHandler(services *service.Services, oauth config.OauthConfig) *Handler  {
	return &Handler{
		services: services,
		oauth: oauth,
	}
}

func (h *Handler) InitRoutes(cfg config.Config) *gin.Engine {
	router := gin.Default()
	router.Use(
		gin.Recovery(),
		gin.Logger(),
		middlewareCors,
	)
	// add swagger route
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port)
	if cfg.Environment != config.EnvLocal {
		docs.SwaggerInfo.Host = cfg.HTTP.Host
	}
	if cfg.Environment != config.Prod {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
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
			shops.GET("/",  h.Find)
			shops.POST("/", h.userIdentity, h.CreateShop)
		}
	}
	return router
}