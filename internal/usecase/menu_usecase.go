package usecase

import (
	"context"

	"api-admin-panel/internal/domain/entity"
	"api-admin-panel/internal/domain/repository"
)

type MenuUsecase interface {
	GetAll(ctx context.Context) ([]entity.Menu, error)
	GetByID(ctx context.Context, id uint64) (*entity.Menu, error)
	Create(ctx context.Context, req *CreateMenuRequest) (*entity.Menu, error)
	Update(ctx context.Context, id uint64, req *UpdateMenuRequest) (*entity.Menu, error)
	Delete(ctx context.Context, id uint64) error
}

type CreateMenuRequest struct {
	Name      string `json:"name" binding:"required,max=100"`
	Link      string `json:"link" binding:"required,max=100"`
	LinkAlias string `json:"link_alias" binding:"required,max=100"`
	Icon      string `json:"icon" binding:"required,max=100"`
	Parent    int64  `json:"parent"`
	Order     int64  `json:"order"`
}

type UpdateMenuRequest struct {
	Name      string `json:"name" binding:"required,max=100"`
	Link      string `json:"link" binding:"required,max=100"`
	LinkAlias string `json:"link_alias" binding:"required,max=100"`
	Icon      string `json:"icon" binding:"required,max=100"`
	Parent    int64  `json:"parent"`
	Order     int64  `json:"order"`
	IsActive  bool   `json:"is_active"`
}

type menuUsecase struct {
	repo repository.MenuRepository
}

func NewMenuUsecase(repo repository.MenuRepository) MenuUsecase {
	return &menuUsecase{repo: repo}
}

func (u *menuUsecase) GetAll(ctx context.Context) ([]entity.Menu, error) {
	return u.repo.FindAll(ctx)
}

func (u *menuUsecase) GetByID(ctx context.Context, id uint64) (*entity.Menu, error) {
	return u.repo.FindByID(ctx, id)
}

func (u *menuUsecase) Create(ctx context.Context, req *CreateMenuRequest) (*entity.Menu, error) {
	menu := &entity.Menu{
		Name:      req.Name,
		Link:      req.Link,
		LinkAlias: req.LinkAlias,
		Icon:      req.Icon,
		Parent:    req.Parent,
		Order:     req.Order,
		IsActive:  true,
	}
	if err := u.repo.Create(ctx, menu); err != nil {
		return nil, err
	}
	return menu, nil
}

func (u *menuUsecase) Update(ctx context.Context, id uint64, req *UpdateMenuRequest) (*entity.Menu, error) {
	menu, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	menu.Name = req.Name
	menu.Link = req.Link
	menu.LinkAlias = req.LinkAlias
	menu.Icon = req.Icon
	menu.Parent = req.Parent
	menu.Order = req.Order
	menu.IsActive = req.IsActive

	if err := u.repo.Update(ctx, menu); err != nil {
		return nil, err
	}
	return menu, nil
}

func (u *menuUsecase) Delete(ctx context.Context, id uint64) error {
	return u.repo.Delete(ctx, id)
}