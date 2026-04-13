package mysql

import (
	"api-admin-panel/internal/domain/entity"
	"api-admin-panel/internal/domain/repository"
	"gorm.io/gorm"
)

type roleMenuPermissionRepository struct {
	db *gorm.DB
}

func NewRoleMenuPermissionRepository(db *gorm.DB) repository.RoleMenuPermissionRepository {
	return &roleMenuPermissionRepository{db: db}
}

func (r *roleMenuPermissionRepository) FindByRoleID(roleID string) ([]entity.RoleMenuPermission, error) {
	var result []entity.RoleMenuPermission
	if err := r.db.
		Preload("Menu").
		Preload("Permission").
		Where("role_id = ?", roleID).
		Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (r *roleMenuPermissionRepository) FindByRoleAndMenu(roleID string, menuID uint64) ([]entity.RoleMenuPermission, error) {
	var result []entity.RoleMenuPermission
	if err := r.db.
		Preload("Permission").
		Where("role_id = ? AND menu_id = ?", roleID, menuID).
		Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (r *roleMenuPermissionRepository) Create(rmp *entity.RoleMenuPermission) error {
	return r.db.Create(rmp).Error
}

func (r *roleMenuPermissionRepository) Delete(roleID string, menuID uint64, permissionID string) error {
	return r.db.
		Where("role_id = ? AND menu_id = ? AND permission_id = ?", roleID, menuID, permissionID).
		Delete(&entity.RoleMenuPermission{}).Error
}
