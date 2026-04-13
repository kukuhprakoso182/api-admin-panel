package repository

import "api-admin-panel/internal/domain/entity"

type PermissionRepository interface {
	FindAll() ([]entity.Permission, error)
	FindByID(id string) (*entity.Permission, error)
	Create(permission *entity.Permission) error
	Update(permission *entity.Permission) error
	Delete(id string) error
}
