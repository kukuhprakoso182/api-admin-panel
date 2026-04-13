package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	pkgjwt "api-admin-panel/pkg/jwt"
	"api-admin-panel/pkg/response"
)

const (
	ContextKeyUserID = "user_id"
	ContextKeyRoleID = "role_id"
	ContextKeyEmail  = "email"
)

func JWTAuth(jwtService pkgjwt.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "authorization header is required")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			response.Unauthorized(c, "invalid authorization header format")
			c.Abort()
			return
		}

		claims, err := jwtService.ValidateToken(parts[1])
		if err != nil {
			response.Unauthorized(c, err.Error())
			c.Abort()
			return
		}

		c.Set(ContextKeyUserID, claims.UserID)
		c.Set(ContextKeyRoleID, claims.RoleID)
		c.Set(ContextKeyEmail, claims.Email)
		c.Next()
	}
}
