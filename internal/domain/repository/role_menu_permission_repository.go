package repository

import (
	"context"

	"api-admin-panel/internal/domain/entity"
)

type RoleMenuPermissionRepository interface {
	FindByRoleID(ctx context.Context, roleID string) ([]entity.RoleMenuPermission, error)
	FindByRoleAndMenu(ctx context.Context, roleID string, menuID uint64) ([]entity.RoleMenuPermission, error)
	Create(ctx context.Context, rmp *entity.RoleMenuPermission) error
	Delete(ctx context.Context, roleID string, menuID uint64, permissionID string) error
	HasPermission(ctx context.Context, roleIDs []string, menuLink, permissionID string) (bool, error)
}