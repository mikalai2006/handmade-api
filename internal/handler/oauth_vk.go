package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/handmade/internal/domain"
	"github.com/spf13/viper"
)

var (
	vk_auth_path = "https://oauth.vk.com/authorize"
	clientID     = os.Getenv("VK_CLIENT_ID")
	clientSecret = os.Getenv("VK_CLIENT_SECRET")
	redirectURI  = "http://localhost:8000/me"
	scope        = []string{"account"}
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

func (h *Handler) Me(c *gin.Context) {
	code := c.Query("code")
	clientUrl := c.Query("state")
	if code == "" {
		newErrorResponse(c, http.StatusInternalServerError, "No correct auth")
		return
	}

	urlR := fmt.Sprintf("https://oauth.vk.com/access_token?client_id=%s&client_secret=%s&redirect_uri=%s&code=%s",
		clientID, clientSecret, redirectURI, code)
	req, _ := http.NewRequest("POST", urlR, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	defer resp.Body.Close()

	token := struct {
		AccessToken string `json:"access_token"`
	}{}

	bytes, _ := io.ReadAll(resp.Body)
	json.Unmarshal(bytes, &token)

	urlR = fmt.Sprintf("https://api.vk.com/method/%s?v=5.131&access_token=%s", "users.get", token.AccessToken)
	req, err = http.NewRequest("GET", urlR, nil)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	defer resp.Body.Close()

	bytes, err = io.ReadAll(resp.Body)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	var bodyResponse VKBodyResponse
	if err := json.Unmarshal(bytes, &bodyResponse); err != nil { // Parse []byte to go struct pointer
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	input := domain.Auth{
		Login:    bodyResponse.Response[0].FirstName,
		Strategy: "jwt",
		Password: "",
		VkId:     fmt.Sprintf("%d", bodyResponse.Response[0].ID),
	}

	user, err := h.services.Authorization.ExistAuth(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if user.Login == "" {
		_, err = h.services.Authorization.CreateAuth(input)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	tokens, err := h.services.Authorization.SignIn(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	// tokenAPI, err := h.services.Authorization.GenerateToken(input)
	// if err != nil {
	// 	newErrorResponse(c, http.StatusInternalServerError, err.Error())
	// 	return
	// }
	Url, err := url.Parse(clientUrl)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	parameters := url.Values{}
	parameters.Add("token", tokens.AccessToken)
	Url.RawQuery = parameters.Encode()
	// c.Redirect(http.StatusMovedPermanently, path)
	c.SetCookie("jwt-handmade", tokens.RefreshToken, viper.GetInt("oauth.timeExpireCookie"), "/", c.Request.URL.Hostname(), false, true)
	c.Redirect(http.StatusFound, Url.String())

}