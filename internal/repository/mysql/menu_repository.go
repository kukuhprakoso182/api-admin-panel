package mysql

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"api-admin-panel/internal/domain"
	"api-admin-panel/internal/domain/entity"
	"api-admin-panel/internal/domain/repository"
)

type menuRepository struct {
	db *gorm.DB
}

func NewMenuRepository(db *gorm.DB) repository.MenuRepository {
	return &menuRepository{db: db}
}

func (r *menuRepository) FindAll(ctx context.Context) ([]entity.Menu, error) {
	var menus []entity.Menu
	if err := r.db.WithContext(ctx).Order("`order` asc").Find(&menus).Error; err != nil {
		return nil, err
	}
	return menus, nil
}

func (r *menuRepository) FindByID(ctx context.Context, id uint64) (*entity.Menu, error) {
	var menu entity.Menu
	if err := r.db.WithContext(ctx).First(&menu, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &menu, nil
}

func (r *menuRepository) Create(ctx context.Context, menu *entity.Menu) error {
	if err := r.db.WithContext(ctx).Create(menu).Error; err != nil {
		if isDuplicateKey(err) {
			return domain.ErrAlreadyExists
		}
		return err
	}
	return nil
}

func (r *menuRepository) Update(ctx context.Context, menu *entity.Menu) error {
	res := r.db.WithContext(ctx).
		Model(&entity.Menu{}).
		Where("id = ?", menu.ID).
		// nama field struct, bukan nama kolom: GORM yang meng-quote `order`
		Select("Name", "Link", "LinkAlias", "Icon", "Parent", "Order", "IsActive").
		Updates(menu)
	if res.Error != nil {
		if isDuplicateKey(res.Error) {
			return domain.ErrAlreadyExists
		}
		return res.Error
	}

	// MySQL mengembalikan RowsAffected = 0 jika tidak ada nilai yang berubah,
	// jadi pastikan record memang tidak ada sebelum menyimpulkan "not found".
	if res.RowsAffected == 0 {
		var count int64
		if err := r.db.WithContext(ctx).Model(&entity.Menu{}).
			Where("id = ?", menu.ID).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			return domain.ErrNotFound
		}
	}
	return nil
}

func (r *menuRepository) Delete(ctx context.Context, id uint64) error {
	res := r.db.WithContext(ctx).Delete(&entity.Menu{}, "id = ?", id)
	if res.Error != nil {
		return translateWriteError(res.Error)
	}
	if res.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}