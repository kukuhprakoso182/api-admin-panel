package usecase

import (
	"context"

	"api-admin-panel/internal/domain/entity"
	"api-admin-panel/internal/domain/repository"
)

type IconUsecase interface {
	GetAll(ctx context.Context) ([]entity.Icon, error)
	GetByID(ctx context.Context, id string) (*entity.Icon, error)
	Create(ctx context.Context, req *CreateIconRequest) (*entity.Icon, error)
	Update(ctx context.Context, id string, req *UpdateIconRequest) (*entity.Icon, error)
	Delete(ctx context.Context, id string) error
}

type CreateIconRequest struct {
	ID      string `json:"id" binding:"required,max=50"`
	Name    string `json:"name" binding:"required,max=50"`
	Section string `json:"section" binding:"required,max=50"`
}

type UpdateIconRequest struct {
	Name     string `json:"name" binding:"required,max=50"`
	Section  string `json:"section" binding:"required,max=50"`
	IsActive bool   `json:"is_active"`
}

type iconUsecase struct {
	repo repository.IconRepository
}

func NewIconUsecase(repo repository.IconRepository) IconUsecase {
	return &iconUsecase{repo: repo}
}

func (u *iconUsecase) GetAll(ctx context.Context) ([]entity.Icon, error) {
	return u.repo.FindAll(ctx)
}

func (u *iconUsecase) GetByID(ctx context.Context, id string) (*entity.Icon, error) {
	return u.repo.FindByID(ctx, id)
}

func (u *iconUsecase) Create(ctx context.Context, req *CreateIconRequest) (*entity.Icon, error) {
	icon := &entity.Icon{
		ID:       req.ID,
		Name:     req.Name,
		Section:  req.Section,
		IsActive: true,
	}
	if err := u.repo.Create(ctx, icon); err != nil {
		return nil, err
	}
	return icon, nil
}

func (u *iconUsecase) Update(ctx context.Context, id string, req *UpdateIconRequest) (*entity.Icon, error) {
	icon, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	icon.Name = req.Name
	icon.Section = req.Section
	icon.IsActive = req.IsActive

	if err := u.repo.Update(ctx, icon); err != nil {
		return nil, err
	}
	return icon, nil
}

func (u *iconUsecase) Delete(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}