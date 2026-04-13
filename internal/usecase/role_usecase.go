package usecase

import (
	"errors"

	"api-admin-panel/internal/domain/entity"
	"api-admin-panel/internal/domain/repository"
	"gorm.io/gorm"
)

type RoleUsecase interface {
	GetAll() ([]entity.Role, error)
	GetByID(id string) (*entity.Role, error)
	Create(req *CreateRoleRequest) (*entity.Role, error)
	Update(id string, req *UpdateRoleRequest) (*entity.Role, error)
	Delete(id string) error
}

type CreateRoleRequest struct {
	ID         string  `json:"id" binding:"required,max=50"`
	Name       string  `json:"name" binding:"required,max=50"`
	ChildRoles *string `json:"child_roles"`
}

type UpdateRoleRequest struct {
	Name       string  `json:"name" binding:"required,max=50"`
	ChildRoles *string `json:"child_roles"`
	IsActive   bool    `json:"is_active"`
}

type roleUsecase struct {
	repo repository.RoleRepository
}

func NewRoleUsecase(repo repository.RoleRepository) RoleUsecase {
	return &roleUsecase{repo: repo}
}

func (u *roleUsecase) GetAll() ([]entity.Role, error) {
	return u.repo.FindAll()
}

func (u *roleUsecase) GetByID(id string) (*entity.Role, error) {
	role, err := u.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return role, nil
}

func (u *roleUsecase) Create(req *CreateRoleRequest) (*entity.Role, error) {
	role := &entity.Role{
		ID:         req.ID,
		Name:       req.Name,
		ChildRoles: req.ChildRoles,
		IsActive:   true,
	}
	if err := u.repo.Create(role); err != nil {
		return nil, err
	}
	return role, nil
}

func (u *roleUsecase) Update(id string, req *UpdateRoleRequest) (*entity.Role, error) {
	role, err := u.GetByID(id)
	if err != nil {
		return nil, err
	}
	role.Name = req.Name
	role.ChildRoles = req.ChildRoles
	role.IsActive = req.IsActive
	if err := u.repo.Update(role); err != nil {
		return nil, err
	}
	return role, nil
}

func (u *roleUsecase) Delete(id string) error {
	if _, err := u.GetByID(id); err != nil {
		return err
	}
	return u.repo.Delete(id)
}
