// internal/middleware/auth.go
package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/Ubaid-Rza-08/go-rest-api/internal/utils"
)

const (
	ContextUserID   = "userID"
	ContextUserRole = "userRole"
	ContextEmail    = "userEmail"
)

// Authenticate validates the Bearer JWT and injects claims into the Gin context.
// Protected routes must include this middleware.
func Authenticate(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Unauthorized(c, "authorization header is required")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			utils.Unauthorized(c, "authorization header format must be: Bearer <token>")
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(parts[1], jwtSecret)
		if err != nil {
			utils.Unauthorized(c, "invalid or expired token")
			c.Abort()
			return
		}

		// Inject claims for downstream handlers
		c.Set(ContextUserID, claims.UserID)
		c.Set(ContextUserRole, claims.Role)
		c.Set(ContextEmail, claims.Email)
		c.Next()
	}
}

// RequireRole restricts a route to users with one of the given roles.
// Must be used AFTER Authenticate.
func RequireRole(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]bool, len(roles))
	for _, r := range roles {
		allowed[r] = true
	}
	return func(c *gin.Context) {
		role, _ := c.Get(ContextUserRole)
		if !allowed[role.(string)] {
			utils.Forbidden(c, "you do not have permission to access this resource")
			c.Abort()
			return
		}
		c.Next()
	}
}