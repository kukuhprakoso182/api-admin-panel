package mysql

import (
	"api-admin-panel/internal/domain/entity"
	"api-admin-panel/internal/domain/repository"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindAll(page, limit int) ([]entity.User, int64, error) {
	var users []entity.User
	var total int64
	offset := (page - 1) * limit

	if err := r.db.Model(&entity.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := r.db.
		Preload("Status").
		Preload("Roles").
		Offset(offset).Limit(limit).
		Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *userRepository) FindByID(id string) (*entity.User, error) {
	var user entity.User
	if err := r.db.
		Preload("Status").
		Preload("Roles").
		First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	if err := r.db.
		Preload("Status").
		Preload("Roles").
		Where("email = ?", email).
		First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(user *entity.User, roleIDs []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		return r.syncRoles(tx, user, roleIDs)
	})
}

func (r *userRepository) Update(user *entity.User, roleIDs []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(user).Error; err != nil {
			return err
		}
		return r.syncRoles(tx, user, roleIDs)
	})
}

func (r *userRepository) Delete(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Hapus pivot terlebih dahulu
		if err := tx.Where("user_id = ?", id).Delete(&entity.UserRole{}).Error; err != nil {
			return err
		}
		return tx.Delete(&entity.User{}, "id = ?", id).Error
	})
}

// syncRoles menghapus semua role lama lalu insert yang baru.
func (r *userRepository) syncRoles(tx *gorm.DB, user *entity.User, roleIDs []string) error {
	if err := tx.Where("user_id = ?", user.ID).Delete(&entity.UserRole{}).Error; err != nil {
		return err
	}
	for _, rid := range roleIDs {
		ur := entity.UserRole{UserID: user.ID, RoleID: rid}
		if err := tx.Create(&ur).Error; err != nil {
			return err
		}
	}
	return nil
}
