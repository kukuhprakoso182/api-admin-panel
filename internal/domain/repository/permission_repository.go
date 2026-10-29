package repository

import (
	"context"

	"api-admin-panel/internal/domain/entity"
)

type PermissionRepository interface {
	FindAll(ctx context.Context) ([]entity.Permission, error)
	FindByID(ctx context.Context, id string) (*entity.Permission, error)
	Create(ctx context.Context, permission *entity.Permission) error
	Update(ctx context.Context, permission *entity.Permission) error
	Delete(ctx context.Context, id string) error
}