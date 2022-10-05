package handler

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/handmade/pkg/auths"
)

const (
	authorizationHeader = "Authorization"
	userCtx = "userId"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)

	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		newErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}

	if len(headerParts[1]) == 0 {
		newErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}

	// parse token
	// userId, err := h.services.Authorization.ParseToken(headerParts[1])
	// if err != nil {
	// 	newErrorResponse(c, http.StatusUnauthorized, err.Error())
	// 	return
	// }
	tokenManager, err := auths.NewManager(os.Getenv("SIGNING_KEY"))
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	id, err := tokenManager.Parse(headerParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	c.Set(userCtx, id)
	// session := sessions.Default(c)
	// user := session.Get(userkey)
	// if user == nil {
	// 	// Abort the request with the appropriate error code
	// 	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	// 	return
	// }
	// logrus.Printf("user session= %s", user)
	// // Continue down the chain to handler etc
	// c.Next()
}

func getUserId(c *gin.Context) (string, error)  {
	id, ok := c.Get(userCtx)
	if !ok {
		newErrorResponse(c, http.StatusUnauthorized, "user id not found")
		return "", errors.New("user not found")
	}

	idInt, ok := id.(string)
	if !ok {
		newErrorResponse(c, http.StatusUnauthorized, "user id is of invalid type")
		return "", errors.New("user not found2")
	}

	return idInt, nil
}