package middleware

import (
	"github.com/gin-gonic/gin"
)

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func JSONAppErrorReporter() gin.HandlerFunc {
	return HadleError(gin.ErrorTypeAny)
}

func HadleError(errorType gin.ErrorType) gin.HandlerFunc {
	return func (c *gin.Context)  {
		c.Next()

		// Skip if no errors
		if c.Errors.Last() == nil {
			return
		}
		// public errors
		err := c.Errors.Last()
		if err == nil {
				return
		}

		if err != nil {
			c.JSON(-1, errorResponse{
				Code: c.Writer.Status(),
				Message: err.Error(),
			})
			c.Abort()
		}
	}
}