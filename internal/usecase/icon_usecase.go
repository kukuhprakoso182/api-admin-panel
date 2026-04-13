package usecase

import (
	"errors"

	"api-admin-panel/internal/domain/entity"
	"api-admin-panel/internal/domain/repository"
	"gorm.io/gorm"
)

type IconUsecase interface {
	GetAll() ([]entity.Icon, error)
	GetByID(id string) (*entity.Icon, error)
	Create(req *CreateIconRequest) (*entity.Icon, error)
	Update(id string, req *UpdateIconRequest) (*entity.Icon, error)
	Delete(id string) error
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

func (u *iconUsecase) GetAll() ([]entity.Icon, error) {
	return u.repo.FindAll()
}

func (u *iconUsecase) GetByID(id string) (*entity.Icon, error) {
	icon, err := u.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return icon, nil
}

func (u *iconUsecase) Create(req *CreateIconRequest) (*entity.Icon, error) {
	icon := &entity.Icon{
		ID:       req.ID,
		Name:     req.Name,
		Section:  req.Section,
		IsActive: true,
	}
	if err := u.repo.Create(icon); err != nil {
		return nil, err
	}
	return icon, nil
}

func (u *iconUsecase) Update(id string, req *UpdateIconRequest) (*entity.Icon, error) {
	icon, err := u.GetByID(id)
	if err != nil {
		return nil, err
	}

	icon.Name = req.Name
	icon.Section = req.Section
	icon.IsActive = req.IsActive

	if err := u.repo.Update(icon); err != nil {
		return nil, err
	}
	return icon, nil
}

func (u *iconUsecase) Delete(id string) error {
	if _, err := u.GetByID(id); err != nil {
		return err
	}
	return u.repo.Delete(id)
}
