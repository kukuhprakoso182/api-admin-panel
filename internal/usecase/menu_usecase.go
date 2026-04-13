package usecase

import (
	"errors"

	"api-admin-panel/internal/domain/entity"
	"api-admin-panel/internal/domain/repository"
	"gorm.io/gorm"
)

type MenuUsecase interface {
	GetAll() ([]entity.Menu, error)
	GetByID(id uint64) (*entity.Menu, error)
	Create(req *CreateMenuRequest) (*entity.Menu, error)
	Update(id uint64, req *UpdateMenuRequest) (*entity.Menu, error)
	Delete(id uint64) error
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

func (u *menuUsecase) GetAll() ([]entity.Menu, error) {
	return u.repo.FindAll()
}

func (u *menuUsecase) GetByID(id uint64) (*entity.Menu, error) {
	menu, err := u.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return menu, nil
}

func (u *menuUsecase) Create(req *CreateMenuRequest) (*entity.Menu, error) {
	menu := &entity.Menu{
		Name:      req.Name,
		Link:      req.Link,
		LinkAlias: req.LinkAlias,
		Icon:      req.Icon,
		Parent:    req.Parent,
		Order:     req.Order,
		IsActive:  true,
	}
	if err := u.repo.Create(menu); err != nil {
		return nil, err
	}
	return menu, nil
}

func (u *menuUsecase) Update(id uint64, req *UpdateMenuRequest) (*entity.Menu, error) {
	menu, err := u.GetByID(id)
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
	if err := u.repo.Update(menu); err != nil {
		return nil, err
	}
	return menu, nil
}

func (u *menuUsecase) Delete(id uint64) error {
	if _, err := u.GetByID(id); err != nil {
		return err
	}
	return u.repo.Delete(id)
}
