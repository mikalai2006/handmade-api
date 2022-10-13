package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/handmade/docs"
	"github.com/mikalai2006/handmade/internal/config"
	"github.com/mikalai2006/handmade/internal/middleware"
	"github.com/mikalai2006/handmade/internal/service"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
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
	router := gin.New()
	router.Use(
		gin.Recovery(),
		gin.Logger(),
		middleware.MiddlewareCors,
		middleware.JSONAppErrorReporter(),
	)
	// add swagger route
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port)
	if cfg.Environment != config.EnvLocal {
		docs.SwaggerInfo.Host = cfg.HTTP.Host
	}
	if cfg.Environment != config.Prod {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// create session
	// store := cookie.NewStore([]byte(os.Getenv("secret")))
	// router.Use(sessions.Sessions("mysession", store))

	auth := router.Group("/auth")
	h.registerAuth(auth)

	api := router.Group("/api")
	h.registerShop(api)
	h.RegisterUser(api)

	oauth := router.Group("/oauth")
	h.registerVkOAuth(oauth)
	h.registerGoogleOAuth(oauth)

	router.NoRoute(func(c *gin.Context) {
		c.AbortWithError(http.StatusNotFound, errors.New("page not found"))
		// .SetMeta(gin.H{
		// 	"code": http.StatusNotFound,
		// 	"status": "error",
		// 	"message": "hello",
	 	// })

	})

	return router
}