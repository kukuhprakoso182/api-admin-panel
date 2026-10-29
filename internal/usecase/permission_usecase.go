package usecase

import (
	"context"

	"api-admin-panel/internal/domain/entity"
	"api-admin-panel/internal/domain/repository"
)

type PermissionUsecase interface {
	GetAll(ctx context.Context) ([]entity.Permission, error)
	GetByID(ctx context.Context, id string) (*entity.Permission, error)
	Create(ctx context.Context, req *CreatePermissionRequest) (*entity.Permission, error)
	Update(ctx context.Context, id string, req *UpdatePermissionRequest) (*entity.Permission, error)
	Delete(ctx context.Context, id string) error
}

type CreatePermissionRequest struct {
	ID   string `json:"id" binding:"required,max=50"`
	Name string `json:"name" binding:"required,max=50"`
}

type UpdatePermissionRequest struct {
	Name     string `json:"name" binding:"required,max=50"`
	IsActive bool   `json:"is_active"`
}

type permissionUsecase struct {
	repo repository.PermissionRepository
}

func NewPermissionUsecase(repo repository.PermissionRepository) PermissionUsecase {
	return &permissionUsecase{repo: repo}
}

func (u *permissionUsecase) GetAll(ctx context.Context) ([]entity.Permission, error) {
	return u.repo.FindAll(ctx)
}

func (u *permissionUsecase) GetByID(ctx context.Context, id string) (*entity.Permission, error) {
	return u.repo.FindByID(ctx, id)
}

func (u *permissionUsecase) Create(ctx context.Context, req *CreatePermissionRequest) (*entity.Permission, error) {
	perm := &entity.Permission{
		ID:       req.ID,
		Name:     req.Name,
		IsActive: true,
	}
	if err := u.repo.Create(ctx, perm); err != nil {
		return nil, err
	}
	return perm, nil
}

func (u *permissionUsecase) Update(ctx context.Context, id string, req *UpdatePermissionRequest) (*entity.Permission, error) {
	perm, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	perm.Name = req.Name
	perm.IsActive = req.IsActive

	if err := u.repo.Update(ctx, perm); err != nil {
		return nil, err
	}
	return perm, nil
}

func (u *permissionUsecase) Delete(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}