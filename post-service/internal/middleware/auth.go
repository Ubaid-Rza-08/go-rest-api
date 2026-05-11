package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/Ubaid-Rza-08/post-service/internal/utils"
)

func Authenticate(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {

		auth := c.GetHeader("Authorization")

		if auth == "" {
			utils.Unauthorized(c, "missing authorization header")
			c.Abort()
			return
		}

		parts := strings.SplitN(auth, " ", 2)

		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.Unauthorized(c, "invalid auth format")
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(parts[1], secret)
		if err != nil {
			utils.Unauthorized(c, "invalid token")
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)

		c.Next()
	}
}