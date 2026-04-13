package mysql

import (
	"api-admin-panel/internal/domain/entity"
	"api-admin-panel/internal/domain/repository"
	"gorm.io/gorm"
)

type iconRepository struct {
	db *gorm.DB
}

func NewIconRepository(db *gorm.DB) repository.IconRepository {
	return &iconRepository{db: db}
}

func (r *iconRepository) FindAll() ([]entity.Icon, error) {
	var icons []entity.Icon
	if err := r.db.Find(&icons).Error; err != nil {
		return nil, err
	}
	return icons, nil
}

func (r *iconRepository) FindByID(id string) (*entity.Icon, error) {
	var icon entity.Icon
	if err := r.db.First(&icon, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &icon, nil
}

func (r *iconRepository) Create(icon *entity.Icon) error {
	return r.db.Create(icon).Error
}

func (r *iconRepository) Update(icon *entity.Icon) error {
	return r.db.Save(icon).Error
}

func (r *iconRepository) Delete(id string) error {
	return r.db.Delete(&entity.Icon{}, "id = ?", id).Error
}
