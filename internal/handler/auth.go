package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/handmade/internal/domain"
	"github.com/mikalai2006/handmade/internal/middleware"
)

func (h *Handler) registerAuth(router *gin.RouterGroup) {
		router.POST("/sign-up", h.SignUp)
		router.POST("/sign-in", h.SignIn)
		router.POST("/logout", h.Logout)
		router.POST("/refresh", h.tokenRefresh)
		router.GET("/refresh", h.tokenRefresh)
		router.GET("/verification/:code", middleware.SetUserIdentity, h.VerificationAuth)
}

// @Summary SignUp
// @Tags auth
// @Description Create account
// @ID create-account
// @Accept json
// @Produce json
// @Param input body domain.Auth true "account info"
// @Success 200 {integer} 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-up [post]
func (h *Handler) SignUp(c *gin.Context) {
	var input  domain.SignInInput

	if err := c.BindJSON(&input); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	id, err := h.services.Authorization.CreateAuth(input)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}



// @Summary SignIn
// @Tags auth
// @Description Login user
// @ID signin-account
// @Accept json
// @Produce json
// @Param input body domain.SignInInput true "credentials"
// @Success 200 {integer} 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-in [post]
func (h *Handler) SignIn(c *gin.Context) {
	// jwt_cookie, _ := c.Cookie("jwt-handmade")
	// fmt.Println("+++++++++++++")
	// fmt.Printf("jwt_handmade = %s", jwt_cookie)
	// fmt.Println("+++++++++++++")
	// session := sessions.Default(c)
	var input domain.SignInInput

	if err := c.BindJSON(&input); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if input.Strategy == "" {
		input.Strategy = "local"
	}

	if input.Email == "" && input.Login == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("request must be with email or login"))
		return
	}

	if input.Strategy == "local" {
		tokens, err := h.services.Authorization.SignIn(input)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		c.SetCookie("jwt-handmade", tokens.RefreshToken, h.oauth.TimeExpireCookie, "/", c.Request.URL.Hostname(), false, true)

		c.JSON(http.StatusOK, domain.ResponseTokens{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		})
	} else {
		fmt.Print("JWT auth")
	}
	// session.Set(userkey, input.Username)
	// session.Save()
}

// @Summary User Refresh Tokens
// @Tags users-auth
// @Description user refresh tokens
// @Accept  json
// @Produce  json
// @Param input body refreshInput true "sign up info"
// @Success 200 {object} tokenResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /users/auth/refresh [post]
func (h *Handler) tokenRefresh(c *gin.Context) {
	jwt_cookie, _ := c.Cookie("jwt-handmade")
	// fmt.Println("jwt_handmade = ", jwt_cookie)
	// jwt_header := c.GetHeader("hello")
	// fmt.Println("jwt_header = ", jwt_header)
	// fmt.Println("+++++++++++++")
	// session := sessions.Default(c)
	var input domain.RefreshInput

	if jwt_cookie == "" {
		if err := c.BindJSON(&input); err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid input body"))
			return
		}
	} else {
		input.Token = jwt_cookie
	}

	res, err := h.services.Authorization.RefreshTokens(input.Token)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.SetCookie("jwt-handmade", res.RefreshToken, h.oauth.TimeExpireCookie, "/", c.Request.URL.Hostname(), false, true)

	c.JSON(http.StatusOK, domain.ResponseTokens{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	})
}


func (h *Handler) Logout(c *gin.Context)  {
	// session := sessions.Default(c)
	// session.Delete(userkey)
	// if err := session.Save(); err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{
		"message": "Succesfully logged out",
	})
}


func (h *Handler) VerificationAuth(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("code empty"))
		return
	}

	userId, err := middleware.GetUserId(c)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := h.services.Authorization.VerificationCode(userId, code); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}