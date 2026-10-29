package mysql

import (
	"context"

	"gorm.io/gorm"

	"api-admin-panel/internal/domain"
	"api-admin-panel/internal/domain/entity"
	"api-admin-panel/internal/domain/repository"
)

type roleMenuPermissionRepository struct {
	db *gorm.DB
}

func NewRoleMenuPermissionRepository(db *gorm.DB) repository.RoleMenuPermissionRepository {
	return &roleMenuPermissionRepository{db: db}
}

func (r *roleMenuPermissionRepository) FindByRoleID(ctx context.Context, roleID string) ([]entity.RoleMenuPermission, error) {
	var result []entity.RoleMenuPermission
	if err := r.db.WithContext(ctx).
		Preload("Menu").
		Preload("Permission").
		Where("role_id = ?", roleID).
		Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (r *roleMenuPermissionRepository) FindByRoleAndMenu(ctx context.Context, roleID string, menuID uint64) ([]entity.RoleMenuPermission, error) {
	var result []entity.RoleMenuPermission
	if err := r.db.WithContext(ctx).
		Preload("Permission").
		Where("role_id = ? AND menu_id = ?", roleID, menuID).
		Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (r *roleMenuPermissionRepository) Create(ctx context.Context, rmp *entity.RoleMenuPermission) error {
	if err := r.db.WithContext(ctx).Create(rmp).Error; err != nil {
		if isDuplicateKey(err) {
			return domain.ErrAlreadyExists
		}
		return err
	}
	return nil
}

func (r *roleMenuPermissionRepository) Delete(ctx context.Context, roleID string, menuID uint64, permissionID string) error {
	res := r.db.WithContext(ctx).
		Where("role_id = ? AND menu_id = ? AND permission_id = ?", roleID, menuID, permissionID).
		Delete(&entity.RoleMenuPermission{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *roleMenuPermissionRepository) HasPermission(
	ctx context.Context, roleIDs []string, menuLink, permissionID string,
) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.RoleMenuPermission{}).
		Joins("Role").
		Joins("Menu").
		Joins("Permission").
		Where("`role_menu_permission`.`role_id` IN ?", roleIDs).
		Where("`role_menu_permission`.`permission_id` = ?", permissionID).
		Where("`Menu`.`link` = ?", menuLink).
		Where("`Role`.`is_active` = ? AND `Menu`.`is_active` = ? AND `Permission`.`is_active` = ?", true, true, true).
		Count(&count).Error
	return count > 0, err
}