package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"industrial-platform-BE/internal/platform"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			status := platform.HTTPStatus(err)

			c.JSON(status, gin.H{
				"error": gin.H{
					"message": platform.AppErrorMessage(err),
					"code":    status,
				},
			})
		}
	}
}

func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.Error(&platform.AppError{
				Err:     platform.ErrInternal,
				Code:    http.StatusInternalServerError,
				Message: err,
			})
		} else {
			c.Error(platform.ErrInternal)
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	})
}
