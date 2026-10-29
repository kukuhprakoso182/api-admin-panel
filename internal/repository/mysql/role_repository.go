package mysql

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"api-admin-panel/internal/domain"
	"api-admin-panel/internal/domain/entity"
	"api-admin-panel/internal/domain/repository"
)

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) repository.RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) FindAll(ctx context.Context) ([]entity.Role, error) {
	var roles []entity.Role
	if err := r.db.WithContext(ctx).Order("name asc").Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *roleRepository) FindByID(ctx context.Context, id string) (*entity.Role, error) {
	var role entity.Role
	if err := r.db.WithContext(ctx).First(&role, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) Create(ctx context.Context, role *entity.Role) error {
	if err := r.db.WithContext(ctx).Create(role).Error; err != nil {
		if isDuplicateKey(err) {
			return domain.ErrAlreadyExists
		}
		return err
	}
	return nil
}

func (r *roleRepository) Update(ctx context.Context, role *entity.Role) error {
	res := r.db.WithContext(ctx).
		Model(&entity.Role{}).
		Where("id = ?", role.ID).
		// ChildRoles ikut di Select agar bisa di-set kembali ke NULL (nil)
		Select("Name", "ChildRoles", "IsActive").
		Updates(role)
	if res.Error != nil {
		return res.Error
	}

	// RowsAffected = 0 bisa berarti "tidak ada yang berubah", bukan "tidak ditemukan"
	if res.RowsAffected == 0 {
		var count int64
		if err := r.db.WithContext(ctx).Model(&entity.Role{}).
			Where("id = ?", role.ID).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			return domain.ErrNotFound
		}
	}
	return nil
}

func (r *roleRepository) Delete(ctx context.Context, id string) error {
	res := r.db.WithContext(ctx).Delete(&entity.Role{}, "id = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}