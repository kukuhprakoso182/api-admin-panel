package repository

import "api-admin-panel/internal/domain/entity"

type RoleMenuPermissionRepository interface {
	FindByRoleID(roleID string) ([]entity.RoleMenuPermission, error)
	FindByRoleAndMenu(roleID string, menuID uint64) ([]entity.RoleMenuPermission, error)
	Create(rmp *entity.RoleMenuPermission) error
	Delete(roleID string, menuID uint64, permissionID string) error
}
