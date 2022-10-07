package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/handmade/internal/domain"
)

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
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Authorization.CreateAuth(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
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
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if input.Strategy == "" {
		input.Strategy = "local"
	}

	if input.Strategy == "local" {
		tokens, err := h.services.Authorization.SignIn(input)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		c.SetCookie("jwt-handmade", tokens.RefreshToken, 900, "/", "", false, true)
		c.JSON(http.StatusOK, map[string]interface{}{
			"token_access": tokens.AccessToken,
			"refresh_token": tokens.RefreshToken,
		})
	} else {
		fmt.Print("JWT auth")
	}
	// session.Set(userkey, input.Username)
	// session.Save()
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

func (h *Handler) OAuthGoogle(c *gin.Context) {
	urlReferer := c.Request.Referer()
	scope := strings.Join(h.oauth.GoogleScopes, " ")

	Url, err := url.Parse(h.oauth.GoogleAuthUri)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	parameters := url.Values{}
	parameters.Add("client_id", h.oauth.GoogleClientId)
	parameters.Add("redirect_uri", h.oauth.GoogleRedirectUri)
	parameters.Add("scope", scope)
	parameters.Add("response_type", "code")
	parameters.Add("state", urlReferer)

	Url.RawQuery = parameters.Encode()
	c.Redirect(http.StatusFound, Url.String())
}

func (h *Handler) OAuthVK(c *gin.Context)  {
	urlReferer := c.Request.Referer()
	scope := strings.Join(h.oauth.VkScopes, "+")

	Url, err := url.Parse(h.oauth.VkAuthUri)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	parameters := url.Values{}
	parameters.Add("client_id", h.oauth.VkClientId)
	parameters.Add("redirect_uri", h.oauth.VkRedirectUri)
	parameters.Add("scope", scope)
	parameters.Add("response_type", "code")
	parameters.Add("state", urlReferer)

	Url.RawQuery = parameters.Encode()
	c.Redirect(http.StatusFound, Url.String())
}
