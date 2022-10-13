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

func (h *Handler) registerGoogleOAuth(router *gin.RouterGroup) {
	router.GET("/google", h.OAuthGoogle)
	router.GET("/google/me", h.MeGoogle)
}

func (h *Handler) OAuthGoogle(c *gin.Context) {
	urlReferer := c.Request.Referer()
	scope := strings.Join(h.oauth.GoogleScopes, " ")

	Url, err := url.Parse(h.oauth.GoogleAuthUri)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
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

func (h *Handler) MeGoogle(c *gin.Context) {

	code := c.Query("code")
	clientUrl := c.Query("state")

	if code == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("no correct code"))
		return
	}

	Url, err := url.Parse(h.oauth.GoogleTokenUri)
	if err != nil {
		panic("boom")
	}
	parameters := url.Values{}
	parameters.Set("client_id", h.oauth.GoogleClientId)
	parameters.Set("redirect_uri", h.oauth.GoogleRedirectUri)
	parameters.Set("client_secret", h.oauth.GoogleClientSecret)
	parameters.Set("code", code)
	parameters.Set("grant_type", "authorization_code")

	req, _ := http.NewRequest("POST", Url.String(), strings.NewReader(parameters.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
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
		return
	}

	UrlInfo, err := url.Parse(h.oauth.GoogleUserinfoUri)
	if err != nil {
		panic("boom")
	}
	r, _ := http.NewRequest(http.MethodGet, UrlInfo.String(), nil) // URL-encoded payload
	bearerToken := fmt.Sprintf("Bearer %s", token.AccessToken)
	r.Header.Add("Authorization", bearerToken)
	// r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err = http.DefaultClient.Do(r)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	defer resp.Body.Close()

	bytes, err = io.ReadAll(resp.Body)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	var bodyResponse GoogleUserInfo
	if err := json.Unmarshal(bytes, &bodyResponse); err != nil { // Parse []byte to go struct pointer
		c.AbortWithError(http.StatusBadRequest, err)
	}

	input := domain.SignInInput{
		Login:    bodyResponse.Email,
		Strategy: "jwt",
		Password: "",
		GoogleId: bodyResponse.Sub,
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

	Url, err = url.Parse(clientUrl)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}
	parameters = url.Values{}
	parameters.Add("token", tokens.AccessToken)
	Url.RawQuery = parameters.Encode()
	c.SetCookie("jwt-handmade", tokens.RefreshToken, h.oauth.TimeExpireCookie, "/", c.Request.URL.Hostname(), false, true)
	c.Redirect(http.StatusFound, Url.String())
}