package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"api-admin-panel/pkg/response"
)

// PermissionChecker dipenuhi oleh repository.RoleMenuPermissionRepository.
type PermissionChecker interface {
	HasPermission(ctx context.Context, roleIDs []string, menuLink, permissionID string) (bool, error)
}

// permissionForMethod memetakan HTTP method ke permission ID pada tabel permissions.
func permissionForMethod(method string) (string, bool) {
	switch method {
	case http.MethodGet, http.MethodHead:
		return "view", true
	case http.MethodPost:
		return "create", true
	case http.MethodPut, http.MethodPatch:
		return "update", true
	case http.MethodDelete:
		return "delete", true
	default:
		return "", false
	}
}

// RequirePermission menolak request kecuali salah satu role user (dari JWT)
// memiliki permission untuk menu dengan link menuLink. Harus dipasang SETELAH JWTAuth.
// Default-deny: method yang tidak dikenal juga ditolak.
func RequirePermission(checker PermissionChecker, menuLink string) gin.HandlerFunc {
	return func(c *gin.Context) {
		permissionID, ok := permissionForMethod(c.Request.Method)
		if !ok {
			response.Forbidden(c, "access forbidden")
			c.Abort()
			return
		}

		roleIDs := RoleIDsFrom(c)
		if len(roleIDs) == 0 {
			response.Forbidden(c, "access forbidden")
			c.Abort()
			return
		}

		allowed, err := checker.HasPermission(c.Request.Context(), roleIDs, menuLink, permissionID)
		if err != nil {
			_ = c.Error(err)
			if c.Request.Context().Err() == context.DeadlineExceeded {
				response.GatewayTimeout(c, "request timeout")
			} else {
				response.InternalServerError(c, "internal server error")
			}
			c.Abort()
			return
		}
		if !allowed {
			response.Forbidden(c, "access forbidden")
			c.Abort()
			return
		}
		c.Next()
	}
}
