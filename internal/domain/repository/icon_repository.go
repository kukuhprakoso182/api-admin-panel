package repository

import "api-admin-panel/internal/domain/entity"

type IconRepository interface {
	FindAll() ([]entity.Icon, error)
	FindByID(id string) (*entity.Icon, error)
	Create(icon *entity.Icon) error
	Update(icon *entity.Icon) error
	Delete(id string) error
}
