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

// UserIDFrom mengambil user_id hasil JWTAuth dari context.
func UserIDFrom(c *gin.Context) string {
	return c.GetString(ContextKeyUserID)
}

// RoleIDsFrom mengambil daftar role (claim role_id berformat CSV) dari context.
func RoleIDsFrom(c *gin.Context) []string {
	raw := c.GetString(ContextKeyRoleID)
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	ids := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			ids = append(ids, p)
		}
	}
	return ids
}
