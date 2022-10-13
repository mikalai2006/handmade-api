package middleware

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

func SetUserIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)

	if header == "" {
		c.AbortWithError(http.StatusUnauthorized, errors.New("empty auth header"))
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		c.AbortWithError(http.StatusUnauthorized, errors.New("invalid auth header"))
		return
	}

	if len(headerParts[1]) == 0 {
		c.AbortWithError(http.StatusUnauthorized, errors.New("invalid auth header"))
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
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	id, err := tokenManager.Parse(headerParts[1])
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
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

func GetUserId(c *gin.Context) (string, error)  {
	id, ok := c.Get(userCtx)
	if !ok {
		return "", errors.New("user not found")
	}

	idInt, ok := id.(string)
	if !ok {
		return "", errors.New("user not found2")
	}

	return idInt, nil
}
