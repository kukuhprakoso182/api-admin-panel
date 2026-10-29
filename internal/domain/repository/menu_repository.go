package repository

import (
	"context"

	"api-admin-panel/internal/domain/entity"
)

type MenuRepository interface {
	FindAll(ctx context.Context) ([]entity.Menu, error)
	FindByID(ctx context.Context, id uint64) (*entity.Menu, error)
	Create(ctx context.Context, menu *entity.Menu) error
	Update(ctx context.Context, menu *entity.Menu) error
	Delete(ctx context.Context, id uint64) error
}