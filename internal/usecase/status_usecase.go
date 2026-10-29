package usecase

import (
	"context"

	"api-admin-panel/internal/domain/entity"
	"api-admin-panel/internal/domain/repository"
)

type StatusUsecase interface {
	GetAll(ctx context.Context) ([]entity.Status, error)
	GetByID(ctx context.Context, id string) (*entity.Status, error)
	Create(ctx context.Context, req *CreateStatusRequest) (*entity.Status, error)
	Update(ctx context.Context, id string, req *UpdateStatusRequest) (*entity.Status, error)
	Delete(ctx context.Context, id string) error
}

type CreateStatusRequest struct {
	ID   string `json:"id" binding:"required,max=50"`
	Name string `json:"name" binding:"required,max=50"`
}

type UpdateStatusRequest struct {
	Name     string `json:"name" binding:"required,max=50"`
	IsActive bool   `json:"is_active"`
}

type statusUsecase struct {
	repo repository.StatusRepository
}

func NewStatusUsecase(repo repository.StatusRepository) StatusUsecase {
	return &statusUsecase{repo: repo}
}

func (u *statusUsecase) GetAll(ctx context.Context) ([]entity.Status, error) {
	return u.repo.FindAll(ctx)
}

func (u *statusUsecase) GetByID(ctx context.Context, id string) (*entity.Status, error) {
	return u.repo.FindByID(ctx, id)
}

func (u *statusUsecase) Create(ctx context.Context, req *CreateStatusRequest) (*entity.Status, error) {
	status := &entity.Status{
		ID:       req.ID,
		Name:     req.Name,
		IsActive: true,
	}
	if err := u.repo.Create(ctx, status); err != nil {
		return nil, err
	}
	return status, nil
}

func (u *statusUsecase) Update(ctx context.Context, id string, req *UpdateStatusRequest) (*entity.Status, error) {
	status, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	status.Name = req.Name
	status.IsActive = req.IsActive

	if err := u.repo.Update(ctx, status); err != nil {
		return nil, err
	}
	return status, nil
}

func (u *statusUsecase) Delete(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}