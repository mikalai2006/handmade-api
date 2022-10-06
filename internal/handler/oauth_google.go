package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/handmade/internal/domain"
	"github.com/spf13/viper"
)

var (
	GOOGLE_SCOPES        = []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"} // Права доступа
	GOOGLE_AUTH_URI      = "https://accounts.google.com/o/oauth2/auth"                                                                    // Посилання на аутентифікацію
	GOOGLE_TOKEN_URI     = "https://accounts.google.com/o/oauth2/token"                                                                   // Посилання на отримання токена
	GOOGLE_USER_INFO_URI = "https://www.googleapis.com/oauth2/v3/userinfo"                                                                // Посилання на отримання інформації про користувача
	// GOOGLE_CLIENT_ID     = os.Getenv("GOOGLE_CLIENT_ID")
	// GOOGLE_CLIENT_SECRET = os.Getenv("GOOGLE_CLIENT_SECRET")
	GOOGLE_REDIRECT_URI  = "http://localhost:8000/googleme"
)

type GoogleUserInfo struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Profile       string `json:"profile"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Gender        string `json:"gender"`
}

func (h *Handler) MeGoogle(c *gin.Context) {

	code := c.Query("code")
	clientUrl := c.Query("state")

	if code == "" {
		newErrorResponse(c, http.StatusInternalServerError, "No correct auth")
		return
	}

	Url, err := url.Parse(GOOGLE_TOKEN_URI)
	if err != nil {
		panic("boom")
	}
	parameters := url.Values{}
	parameters.Set("client_id", os.Getenv("GOOGLE_CLIENT_ID"))
	parameters.Set("redirect_uri", GOOGLE_REDIRECT_URI)
	parameters.Set("client_secret", os.Getenv("GOOGLE_CLIENT_SECRET"))
	parameters.Set("code", code)
	parameters.Set("grant_type", "authorization_code")

	req, _ := http.NewRequest("POST", Url.String(), strings.NewReader(parameters.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
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

	UrlInfo, err := url.Parse(GOOGLE_USER_INFO_URI)
	if err != nil {
		panic("boom")
	}
	r, _ := http.NewRequest(http.MethodGet, UrlInfo.String(), nil) // URL-encoded payload
	bearerToken := fmt.Sprintf("Bearer %s", token.AccessToken)
	r.Header.Add("Authorization", bearerToken)
	// r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err = http.DefaultClient.Do(r)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	defer resp.Body.Close()

	bytes, err = io.ReadAll(resp.Body)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	var bodyResponse GoogleUserInfo
	if err := json.Unmarshal(bytes, &bodyResponse); err != nil { // Parse []byte to go struct pointer
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	input := domain.Auth{
		Login:    bodyResponse.Email,
		Strategy: "jwt",
		Password: "",
		GoogleId: bodyResponse.Sub,
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

	Url, err = url.Parse(clientUrl)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	parameters = url.Values{}
	parameters.Add("token", tokens.AccessToken)
	Url.RawQuery = parameters.Encode()
	c.SetCookie("jwt-handmade", tokens.RefreshToken, viper.GetInt("oauth.timeExpireCookie"), "/", c.Request.URL.Hostname(), false, true)
	c.Redirect(http.StatusFound, Url.String())
}