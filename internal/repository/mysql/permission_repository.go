package mysql

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"api-admin-panel/internal/domain"
	"api-admin-panel/internal/domain/entity"
	"api-admin-panel/internal/domain/repository"
)

type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) repository.PermissionRepository {
	return &permissionRepository{db: db}
}

func (r *permissionRepository) FindAll(ctx context.Context) ([]entity.Permission, error) {
	var permissions []entity.Permission
	if err := r.db.WithContext(ctx).Order("name asc").Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

func (r *permissionRepository) FindByID(ctx context.Context, id string) (*entity.Permission, error) {
	var permission entity.Permission
	if err := r.db.WithContext(ctx).First(&permission, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &permission, nil
}

func (r *permissionRepository) Create(ctx context.Context, permission *entity.Permission) error {
	if err := r.db.WithContext(ctx).Create(permission).Error; err != nil {
		if isDuplicateKey(err) {
			return domain.ErrAlreadyExists
		}
		return err
	}
	return nil
}

func (r *permissionRepository) Update(ctx context.Context, permission *entity.Permission) error {
	res := r.db.WithContext(ctx).
		Model(&entity.Permission{}).
		Where("id = ?", permission.ID).
		Select("Name", "IsActive"). // agar is_active=false ikut tersimpan
		Updates(permission)
	if res.Error != nil {
		return res.Error
	}

	// RowsAffected = 0 bisa berarti "tidak ada yang berubah", bukan "tidak ditemukan"
	if res.RowsAffected == 0 {
		var count int64
		if err := r.db.WithContext(ctx).Model(&entity.Permission{}).
			Where("id = ?", permission.ID).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			return domain.ErrNotFound
		}
	}
	return nil
}

func (r *permissionRepository) Delete(ctx context.Context, id string) error {
	res := r.db.WithContext(ctx).Delete(&entity.Permission{}, "id = ?", id)
	if res.Error != nil {
		return translateWriteError(res.Error)
	}
	if res.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}