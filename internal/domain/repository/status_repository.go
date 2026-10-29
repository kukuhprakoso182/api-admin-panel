package repository

import (
	"context"

	"api-admin-panel/internal/domain/entity"
)

type StatusRepository interface {
	FindAll(ctx context.Context) ([]entity.Status, error)
	FindByID(ctx context.Context, id string) (*entity.Status, error)
	Create(ctx context.Context, status *entity.Status) error
	Update(ctx context.Context, status *entity.Status) error
	Delete(ctx context.Context, id string) error
}