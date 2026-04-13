package mysql

import (
	"api-admin-panel/internal/domain/entity"
	"api-admin-panel/internal/domain/repository"
	"gorm.io/gorm"
)

type menuRepository struct {
	db *gorm.DB
}

func NewMenuRepository(db *gorm.DB) repository.MenuRepository {
	return &menuRepository{db: db}
}

func (r *menuRepository) FindAll() ([]entity.Menu, error) {
	var menus []entity.Menu
	if err := r.db.Order("`order` asc").Find(&menus).Error; err != nil {
		return nil, err
	}
	return menus, nil
}

func (r *menuRepository) FindByID(id uint64) (*entity.Menu, error) {
	var menu entity.Menu
	if err := r.db.First(&menu, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &menu, nil
}

func (r *menuRepository) Create(menu *entity.Menu) error {
	return r.db.Create(menu).Error
}

func (r *menuRepository) Update(menu *entity.Menu) error {
	return r.db.Save(menu).Error
}

func (r *menuRepository) Delete(id uint64) error {
	return r.db.Delete(&entity.Menu{}, "id = ?", id).Error
}
