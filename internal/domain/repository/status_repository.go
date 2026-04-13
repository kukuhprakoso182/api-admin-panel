package repository

import "api-admin-panel/internal/domain/entity"

type StatusRepository interface {
	FindAll() ([]entity.Status, error)
	FindByID(id string) (*entity.Status, error)
	Create(status *entity.Status) error
	Update(status *entity.Status) error
	Delete(id string) error
}
