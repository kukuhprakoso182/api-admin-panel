package usecase

import (
	"errors"

	"api-admin-panel/internal/domain/entity"
	"api-admin-panel/internal/domain/repository"
	"gorm.io/gorm"
)

type StatusUsecase interface {
	GetAll() ([]entity.Status, error)
	GetByID(id string) (*entity.Status, error)
	Create(req *CreateStatusRequest) (*entity.Status, error)
	Update(id string, req *UpdateStatusRequest) (*entity.Status, error)
	Delete(id string) error
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

func (u *statusUsecase) GetAll() ([]entity.Status, error) {
	return u.repo.FindAll()
}

func (u *statusUsecase) GetByID(id string) (*entity.Status, error) {
	status, err := u.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return status, nil
}

func (u *statusUsecase) Create(req *CreateStatusRequest) (*entity.Status, error) {
	status := &entity.Status{
		ID:       req.ID,
		Name:     req.Name,
		IsActive: true,
	}
	if err := u.repo.Create(status); err != nil {
		return nil, err
	}
	return status, nil
}

func (u *statusUsecase) Update(id string, req *UpdateStatusRequest) (*entity.Status, error) {
	status, err := u.GetByID(id)
	if err != nil {
		return nil, err
	}
	status.Name = req.Name
	status.IsActive = req.IsActive
	if err := u.repo.Update(status); err != nil {
		return nil, err
	}
	return status, nil
}

func (u *statusUsecase) Delete(id string) error {
	if _, err := u.GetByID(id); err != nil {
		return err
	}
	return u.repo.Delete(id)
}
