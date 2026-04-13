package repository

import "api-admin-panel/internal/domain/entity"

type RoleRepository interface {
	FindAll() ([]entity.Role, error)
	FindByID(id string) (*entity.Role, error)
	Create(role *entity.Role) error
	Update(role *entity.Role) error
	Delete(id string) error
}
