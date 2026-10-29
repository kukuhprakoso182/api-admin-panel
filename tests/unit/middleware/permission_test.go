package middleware_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"api-admin-panel/internal/delivery/http/middleware"
	"api-admin-panel/tests/testhelper"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type fakeChecker struct {
	allowed bool
	err     error

	called   bool
	gotRoles []string
	gotMenu  string
	gotPerm  string
}

func (f *fakeChecker) HasPermission(_ context.Context, roleIDs []string, menuLink, permissionID string) (bool, error) {
	f.called = true
	f.gotRoles, f.gotMenu, f.gotPerm = roleIDs, menuLink, permissionID
	return f.allowed, f.err
}

func routerWith(roles string, checker middleware.PermissionChecker, method string) *gin.Engine {
	r := testhelper.NewRouter()
	r.Use(func(c *gin.Context) {
		if roles != "" {
			c.Set(middleware.ContextKeyRoleID, roles)
		}
		c.Next()
	})
	r.Handle(method, "/users", middleware.RequirePermission(checker, "users"), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	return r
}

func TestRequirePermission_AllowsAndMapsMethod(t *testing.T) {
	cases := map[string]string{
		http.MethodGet:    "view",
		http.MethodPost:   "create",
		http.MethodPut:    "update",
		http.MethodDelete: "delete",
	}
	for method, perm := range cases {
		fc := &fakeChecker{allowed: true}
		w := testhelper.Do(routerWith("admin, super_admin", fc, method), testhelper.NewRequest(method, "/users", nil))

		assert.Equal(t, http.StatusOK, w.Code, method)
		assert.Equal(t, "users", fc.gotMenu, method)
		assert.Equal(t, perm, fc.gotPerm, method)
		assert.Equal(t, []string{"admin", "super_admin"}, fc.gotRoles, method)
	}
}

func TestRequirePermission_Denied(t *testing.T) {
	fc := &fakeChecker{allowed: false}
	w := testhelper.Do(routerWith("admin", fc, http.MethodGet), testhelper.NewRequest(http.MethodGet, "/users", nil))

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestRequirePermission_NoRolesDeniedWithoutQuery(t *testing.T) {
	fc := &fakeChecker{allowed: true}
	w := testhelper.Do(routerWith("", fc, http.MethodGet), testhelper.NewRequest(http.MethodGet, "/users", nil))

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.False(t, fc.called)
}

func TestRequirePermission_CheckerErrorIs500(t *testing.T) {
	fc := &fakeChecker{err: errors.New("db down")}
	w := testhelper.Do(routerWith("admin", fc, http.MethodGet), testhelper.NewRequest(http.MethodGet, "/users", nil))

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestRequirePermission_UnknownMethodDenied(t *testing.T) {
	fc := &fakeChecker{allowed: true}
	w := testhelper.Do(routerWith("admin", fc, http.MethodOptions), testhelper.NewRequest(http.MethodOptions, "/users", nil))

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.False(t, fc.called)
}
