package repository

import (
	"context"

	"api-admin-panel/internal/domain/entity"
)

type RoleRepository interface {
	FindAll(ctx context.Context) ([]entity.Role, error)
	FindByID(ctx context.Context, id string) (*entity.Role, error)
	Create(ctx context.Context, role *entity.Role) error
	Update(ctx context.Context, role *entity.Role) error
	Delete(ctx context.Context, id string) error
}