package mysql

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"api-admin-panel/internal/domain"
	"api-admin-panel/internal/domain/entity"
	"api-admin-panel/internal/domain/repository"
)

type statusRepository struct {
	db *gorm.DB
}

func NewStatusRepository(db *gorm.DB) repository.StatusRepository {
	return &statusRepository{db: db}
}

func (r *statusRepository) FindAll(ctx context.Context) ([]entity.Status, error) {
	var statuses []entity.Status
	if err := r.db.WithContext(ctx).Order("name asc").Find(&statuses).Error; err != nil {
		return nil, err
	}
	return statuses, nil
}

func (r *statusRepository) FindByID(ctx context.Context, id string) (*entity.Status, error) {
	var status entity.Status
	if err := r.db.WithContext(ctx).First(&status, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &status, nil
}

func (r *statusRepository) Create(ctx context.Context, status *entity.Status) error {
	if err := r.db.WithContext(ctx).Create(status).Error; err != nil {
		if isDuplicateKey(err) {
			return domain.ErrAlreadyExists
		}
		return err
	}
	return nil
}

func (r *statusRepository) Update(ctx context.Context, status *entity.Status) error {
	res := r.db.WithContext(ctx).
		Model(&entity.Status{}).
		Where("id = ?", status.ID).
		Select("Name", "IsActive"). // agar is_active=false ikut tersimpan
		Updates(status)
	if res.Error != nil {
		return res.Error
	}

	// RowsAffected = 0 bisa berarti "tidak ada yang berubah", bukan "tidak ditemukan"
	if res.RowsAffected == 0 {
		var count int64
		if err := r.db.WithContext(ctx).Model(&entity.Status{}).
			Where("id = ?", status.ID).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			return domain.ErrNotFound
		}
	}
	return nil
}

func (r *statusRepository) Delete(ctx context.Context, id string) error {
	res := r.db.WithContext(ctx).Delete(&entity.Status{}, "id = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}