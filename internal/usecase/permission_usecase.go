package usecase

import (
	"errors"

	"api-admin-panel/internal/domain/entity"
	"api-admin-panel/internal/domain/repository"
	"gorm.io/gorm"
)

type PermissionUsecase interface {
	GetAll() ([]entity.Permission, error)
	GetByID(id string) (*entity.Permission, error)
	Create(req *CreatePermissionRequest) (*entity.Permission, error)
	Update(id string, req *UpdatePermissionRequest) (*entity.Permission, error)
	Delete(id string) error
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

func (u *permissionUsecase) GetAll() ([]entity.Permission, error) {
	return u.repo.FindAll()
}

func (u *permissionUsecase) GetByID(id string) (*entity.Permission, error) {
	perm, err := u.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return perm, nil
}

func (u *permissionUsecase) Create(req *CreatePermissionRequest) (*entity.Permission, error) {
	perm := &entity.Permission{
		ID:       req.ID,
		Name:     req.Name,
		IsActive: true,
	}
	if err := u.repo.Create(perm); err != nil {
		return nil, err
	}
	return perm, nil
}

func (u *permissionUsecase) Update(id string, req *UpdatePermissionRequest) (*entity.Permission, error) {
	perm, err := u.GetByID(id)
	if err != nil {
		return nil, err
	}
	perm.Name = req.Name
	perm.IsActive = req.IsActive
	if err := u.repo.Update(perm); err != nil {
		return nil, err
	}
	return perm, nil
}

func (u *permissionUsecase) Delete(id string) error {
	if _, err := u.GetByID(id); err != nil {
		return err
	}
	return u.repo.Delete(id)
}
