package repository

import "api-admin-panel/internal/domain/entity"

type UserRepository interface {
	FindAll(page, limit int) ([]entity.User, int64, error)
	FindByID(id string) (*entity.User, error)
	FindByEmail(email string) (*entity.User, error)
	Create(user *entity.User, roleIDs []string) error
	Update(user *entity.User, roleIDs []string) error
	Delete(id string) error
}
