package repository

import "api-admin-panel/internal/domain/entity"

type MenuRepository interface {
	FindAll() ([]entity.Menu, error)
	FindByID(id uint64) (*entity.Menu, error)
	Create(menu *entity.Menu) error
	Update(menu *entity.Menu) error
	Delete(id uint64) error
}
