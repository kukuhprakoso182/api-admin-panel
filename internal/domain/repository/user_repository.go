package repository

import (
	"context"

	"api-admin-panel/internal/domain/entity"
)

type UserRepository interface {
	FindAll(ctx context.Context, page, limit int) ([]entity.User, int64, error)
	FindByID(ctx context.Context, id string) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	Create(ctx context.Context, user *entity.User, roleIDs []string) error
	Update(ctx context.Context, user *entity.User, roleIDs []string) error
	Delete(ctx context.Context, id string) error
}