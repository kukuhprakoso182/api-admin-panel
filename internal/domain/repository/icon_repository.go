package repository

import (
	"context"

	"api-admin-panel/internal/domain/entity"
)

type IconRepository interface {
	FindAll(ctx context.Context) ([]entity.Icon, error)
	FindByID(ctx context.Context, id string) (*entity.Icon, error)
	Create(ctx context.Context, icon *entity.Icon) error
	Update(ctx context.Context, icon *entity.Icon) error
	Delete(ctx context.Context, id string) error
}