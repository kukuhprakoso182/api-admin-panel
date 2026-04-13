package mysql

import (
	"api-admin-panel/internal/domain/entity"
	"api-admin-panel/internal/domain/repository"
	"gorm.io/gorm"
)

type statusRepository struct {
	db *gorm.DB
}

func NewStatusRepository(db *gorm.DB) repository.StatusRepository {
	return &statusRepository{db: db}
}

func (r *statusRepository) FindAll() ([]entity.Status, error) {
	var statuses []entity.Status
	if err := r.db.Find(&statuses).Error; err != nil {
		return nil, err
	}
	return statuses, nil
}

func (r *statusRepository) FindByID(id string) (*entity.Status, error) {
	var status entity.Status
	if err := r.db.First(&status, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &status, nil
}

func (r *statusRepository) Create(status *entity.Status) error {
	return r.db.Create(status).Error
}

func (r *statusRepository) Update(status *entity.Status) error {
	return r.db.Save(status).Error
}

func (r *statusRepository) Delete(id string) error {
	return r.db.Delete(&entity.Status{}, "id = ?", id).Error
}
