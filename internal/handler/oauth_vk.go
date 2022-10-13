package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/handmade/internal/domain"
)

type VKBodyResponse struct {
	Response []struct {
		ID              int    `json:"id"`
		FirstName       string `json:"first_name"`
		LastName        string `json:"last_name"`
		CanAccessClosed bool   `json:"can_access_closed"`
		IsClosed        bool   `json:"is_closed"`
	} `json:"response"`
}

func (h *Handler) registerVkOAuth(router *gin.RouterGroup) {
	router.GET("/vk", h.OAuthVK)
	router.GET("/vk/me", h.MeVk)
}

func (h *Handler) OAuthVK(c *gin.Context)  {
	urlReferer := c.Request.Referer()
	scope := strings.Join(h.oauth.VkScopes, "+")

	Url, err := url.Parse(h.oauth.VkAuthUri)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
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


func (h *Handler) MeVk(c *gin.Context) {
	code := c.Query("code")
	clientUrl := c.Query("state")
	if code == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("no correct code"))
		return
	}

	Url, err := url.Parse(h.oauth.VkTokenUri)
	if err != nil {
		panic("boom")
	}
	parameters := url.Values{}
	parameters.Set("client_id", h.oauth.VkClientId)
	parameters.Set("client_secret", h.oauth.VkClientSecret)
	parameters.Set("redirect_uri", h.oauth.VkRedirectUri)
	parameters.Set("code", code)

	req, _ := http.NewRequest("POST", Url.String(), strings.NewReader(parameters.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// urlR := fmt.Sprintf("https://oauth.vk.com/access_token?client_id=%s&client_secret=%s&redirect_uri=%s&code=%s",
	// h.oauth.VkClientId, h.oauth.VkClientSecret, h.oauth.VkRedirectUri, code)
	// req, _ := http.NewRequest("POST", urlR, nil)
	// resp, err := http.DefaultClient.Do(req)
	// if err != nil {
	// 	newErrorResponse(c, http.StatusInternalServerError, err.Error())
	// 	return
	// }
	defer resp.Body.Close()




	token := struct {
		AccessToken string `json:"access_token"`
	}{}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := json.Unmarshal(bytes, &token); err != nil { // Parse []byte to go struct pointer
		c.AbortWithError(http.StatusBadRequest, err)
	}

	Url, err = url.Parse(h.oauth.VkUserinfoUri)
	if err != nil {
		panic("boom")
	}
	parameters = url.Values{}
	parameters.Set("access_token", token.AccessToken)
	parameters.Set("v", "5.131")

	req, _ = http.NewRequest("POST", Url.String(), strings.NewReader(parameters.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// urlR = fmt.Sprintf("https://api.vk.com/method/%s?v=5.131&access_token=%s", "users.get", token.AccessToken)
	// req, err = http.NewRequest("GET", urlR, nil)
	// if err != nil {
	// 	newErrorResponse(c, http.StatusInternalServerError, err.Error())
	// 	return
	// }
	// resp, err = http.DefaultClient.Do(req)
	// if err != nil {
	// 	newErrorResponse(c, http.StatusInternalServerError, err.Error())
	// 	return
	// }
	defer resp.Body.Close()

	bytes, err = io.ReadAll(resp.Body)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	var bodyResponse VKBodyResponse
	if err := json.Unmarshal(bytes, &bodyResponse); err != nil { // Parse []byte to go struct pointer
		c.AbortWithError(http.StatusBadRequest, err)
	}

	input := domain.SignInInput{
		Login:    bodyResponse.Response[0].FirstName,
		Strategy: "jwt",
		Password: "",
		VkId:     fmt.Sprintf("%d", bodyResponse.Response[0].ID),
	}

	user, err := h.services.Authorization.ExistAuth(input)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if user.Login == "" {
		_, err = h.services.Authorization.CreateAuth(input)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
	}

	tokens, err := h.services.Authorization.SignIn(input)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	// tokenAPI, err := h.services.Authorization.GenerateToken(input)
	// if err != nil {
	// 	newErrorResponse(c, http.StatusInternalServerError, err.Error())
	// 	return
	// }
	Url, err = url.Parse(clientUrl)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}
	parameters = url.Values{}
	parameters.Add("token", tokens.AccessToken)
	Url.RawQuery = parameters.Encode()
	// c.Redirect(http.StatusMovedPermanently, path)
	c.SetCookie("jwt-handmade", tokens.RefreshToken, h.oauth.TimeExpireCookie, "/", c.Request.URL.Hostname(), false, true)
	c.Redirect(http.StatusFound, Url.String())

}