package mysql

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"api-admin-panel/internal/domain"
	"api-admin-panel/internal/domain/entity"
	"api-admin-panel/internal/domain/repository"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindAll(ctx context.Context, page, limit int) ([]entity.User, int64, error) {
	var users []entity.User
	var total int64
	offset := (page - 1) * limit

	db := r.db.WithContext(ctx)

	if err := db.Model(&entity.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.
		Preload("Status").
		Preload("Roles").
		Order("name asc, id asc"). // id sebagai tie-breaker agar urutan stabil
		Offset(offset).Limit(limit).
		Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).
		Preload("Status").
		Preload("Roles").
		First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).
		Preload("Status").
		Preload("Roles").
		Where("email = ?", email).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, user *entity.User, roleIDs []string) error {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit(clause.Associations).Create(user).Error; err != nil {
			return err
		}
		return syncRoles(tx, user.ID, roleIDs)
	})
	return translateWriteError(err)
}

func (r *userRepository) Update(ctx context.Context, user *entity.User, roleIDs []string) error {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&entity.User{}).
			Where("id = ?", user.ID).
			Omit(clause.Associations). // jangan ikut menyimpan Status/Roles yang ter-preload
			Select("StatusID", "Name", "Email", "Password").
			Updates(user)
		if res.Error != nil {
			return res.Error
		}

		// RowsAffected = 0 bisa berarti "tidak ada yang berubah", bukan "tidak ditemukan"
		if res.RowsAffected == 0 {
			var count int64
			if err := tx.Model(&entity.User{}).Where("id = ?", user.ID).Count(&count).Error; err != nil {
				return err
			}
			if count == 0 {
				return domain.ErrNotFound
			}
		}
		return syncRoles(tx, user.ID, roleIDs)
	})
	return translateWriteError(err)
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Hapus pivot terlebih dahulu
		if err := tx.Where("user_id = ?", id).Delete(&entity.UserRole{}).Error; err != nil {
			return err
		}
		res := tx.Delete(&entity.User{}, "id = ?", id)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return domain.ErrNotFound // transaksi di-rollback
		}
		return nil
	})
	return translateWriteError(err)
}

// syncRoles menghapus semua role lama lalu insert yang baru dalam satu batch.
func syncRoles(tx *gorm.DB, userID string, roleIDs []string) error {
	if err := tx.Where("user_id = ?", userID).Delete(&entity.UserRole{}).Error; err != nil {
		return err
	}

	ids := uniqueStrings(roleIDs)
	if len(ids) == 0 {
		return nil
	}

	links := make([]entity.UserRole, 0, len(ids))
	for _, rid := range ids {
		links = append(links, entity.UserRole{UserID: userID, RoleID: rid})
	}
	return tx.Create(&links).Error
}

func uniqueStrings(in []string) []string {
	seen := make(map[string]struct{}, len(in))
	out := make([]string, 0, len(in))
	for _, s := range in {
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	return out
}