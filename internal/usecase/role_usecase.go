package usecase

import (
	"context"

	"api-admin-panel/internal/domain/entity"
	"api-admin-panel/internal/domain/repository"
)

type RoleUsecase interface {
	GetAll(ctx context.Context) ([]entity.Role, error)
	GetByID(ctx context.Context, id string) (*entity.Role, error)
	Create(ctx context.Context, req *CreateRoleRequest) (*entity.Role, error)
	Update(ctx context.Context, id string, req *UpdateRoleRequest) (*entity.Role, error)
	Delete(ctx context.Context, id string) error
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

func (u *roleUsecase) GetAll(ctx context.Context) ([]entity.Role, error) {
	return u.repo.FindAll(ctx)
}

func (u *roleUsecase) GetByID(ctx context.Context, id string) (*entity.Role, error) {
	return u.repo.FindByID(ctx, id)
}

func (u *roleUsecase) Create(ctx context.Context, req *CreateRoleRequest) (*entity.Role, error) {
	role := &entity.Role{
		ID:         req.ID,
		Name:       req.Name,
		ChildRoles: req.ChildRoles,
		IsActive:   true,
	}
	if err := u.repo.Create(ctx, role); err != nil {
		return nil, err
	}
	return role, nil
}

func (u *roleUsecase) Update(ctx context.Context, id string, req *UpdateRoleRequest) (*entity.Role, error) {
	role, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	role.Name = req.Name
	role.ChildRoles = req.ChildRoles
	role.IsActive = req.IsActive

	if err := u.repo.Update(ctx, role); err != nil {
		return nil, err
	}
	return role, nil
}

func (u *roleUsecase) Delete(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}