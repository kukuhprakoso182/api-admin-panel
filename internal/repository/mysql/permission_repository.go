package mysql

import (
	"api-admin-panel/internal/domain/entity"
	"api-admin-panel/internal/domain/repository"
	"gorm.io/gorm"
)

type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) repository.PermissionRepository {
	return &permissionRepository{db: db}
}

func (r *permissionRepository) FindAll() ([]entity.Permission, error) {
	var permissions []entity.Permission
	if err := r.db.Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

func (r *permissionRepository) FindByID(id string) (*entity.Permission, error) {
	var permission entity.Permission
	if err := r.db.First(&permission, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *permissionRepository) Create(permission *entity.Permission) error {
	return r.db.Create(permission).Error
}

func (r *permissionRepository) Update(permission *entity.Permission) error {
	return r.db.Save(permission).Error
}

func (r *permissionRepository) Delete(id string) error {
	return r.db.Delete(&entity.Permission{}, "id = ?", id).Error
}
