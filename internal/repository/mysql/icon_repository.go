package mysql

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"api-admin-panel/internal/domain"
	"api-admin-panel/internal/domain/entity"
	"api-admin-panel/internal/domain/repository"
)

type iconRepository struct {
	db *gorm.DB
}

func NewIconRepository(db *gorm.DB) repository.IconRepository {
	return &iconRepository{db: db}
}

func (r *iconRepository) FindAll(ctx context.Context) ([]entity.Icon, error) {
	var icons []entity.Icon
	if err := r.db.WithContext(ctx).Order("section, name").Find(&icons).Error; err != nil {
		return nil, err
	}
	return icons, nil
}

func (r *iconRepository) FindByID(ctx context.Context, id string) (*entity.Icon, error) {
	var icon entity.Icon
	if err := r.db.WithContext(ctx).First(&icon, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &icon, nil
}

func (r *iconRepository) Create(ctx context.Context, icon *entity.Icon) error {
	if err := r.db.WithContext(ctx).Create(icon).Error; err != nil {
		if isDuplicateKey(err) {
			return domain.ErrAlreadyExists
		}
		return err
	}
	return nil
}

func (r *iconRepository) Update(ctx context.Context, icon *entity.Icon) error {
	res := r.db.WithContext(ctx).
		Model(&entity.Icon{}).
		Where("id = ?", icon.ID).
		Select("name", "section", "is_active"). // agar is_active=false ikut tersimpan
		Updates(icon)
	if res.Error != nil {
		return res.Error
	}

	// MySQL mengembalikan RowsAffected = 0 jika nilainya tidak berubah,
	// jadi cek keberadaan record dulu sebelum menyimpulkan "not found".
	if res.RowsAffected == 0 {
		var count int64
		if err := r.db.WithContext(ctx).Model(&entity.Icon{}).
			Where("id = ?", icon.ID).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			return domain.ErrNotFound
		}
	}
	return nil
}

func (r *iconRepository) Delete(ctx context.Context, id string) error {
	res := r.db.WithContext(ctx).Delete(&entity.Icon{}, "id = ?", id)
	if res.Error != nil {
		return translateWriteError(res.Error)
	}
	if res.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}
