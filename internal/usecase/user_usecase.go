package usecase

import (
	"context"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"api-admin-panel/internal/domain/entity"
	"api-admin-panel/internal/domain/repository"
)

type UserUsecase interface {
	GetAll(ctx context.Context, page, limit int) ([]entity.User, int64, error)
	GetByID(ctx context.Context, id string) (*entity.User, error)
	Create(ctx context.Context, req *CreateUserRequest) (*entity.User, error)
	Update(ctx context.Context, id string, req *UpdateUserRequest) (*entity.User, error)
	Delete(ctx context.Context, id string) error
}

type CreateUserRequest struct {
	StatusID string   `json:"status_id" binding:"required,max=50"`
	RoleIDs  []string `json:"role_ids"  binding:"required,min=1"`
	Name     string   `json:"name"      binding:"required,max=100"`
	Email    string   `json:"email"     binding:"required,email,max=100"`
	Password string   `json:"password"  binding:"required,min=8,max=100"`
}

type UpdateUserRequest struct {
	StatusID string   `json:"status_id" binding:"required,max=50"`
	RoleIDs  []string `json:"role_ids"  binding:"required,min=1"`
	Name     string   `json:"name"      binding:"required,max=100"`
	Email    string   `json:"email"     binding:"required,email,max=100"`
	Password string   `json:"password"  binding:"omitempty,min=8,max=100"`
}

type userUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) UserUsecase {
	return &userUsecase{repo: repo}
}

func (u *userUsecase) GetAll(ctx context.Context, page, limit int) ([]entity.User, int64, error) {
	return u.repo.FindAll(ctx, page, limit)
}

func (u *userUsecase) GetByID(ctx context.Context, id string) (*entity.User, error) {
	return u.repo.FindByID(ctx, id)
}

func (u *userUsecase) Create(ctx context.Context, req *CreateUserRequest) (*entity.User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		ID:       uuid.New().String(),
		StatusID: req.StatusID,
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashed),
	}

	if err := u.repo.Create(ctx, user, req.RoleIDs); err != nil {
		return nil, err
	}

	// Reload agar relasi terisi
	return u.repo.FindByID(ctx, user.ID)
}

func (u *userUsecase) Update(ctx context.Context, id string, req *UpdateUserRequest) (*entity.User, error) {
	user, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	user.StatusID = req.StatusID
	user.Name = req.Name
	user.Email = req.Email

	if req.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(hashed)
	}

	if err := u.repo.Update(ctx, user, req.RoleIDs); err != nil {
		return nil, err
	}

	return u.repo.FindByID(ctx, user.ID)
}

func (u *userUsecase) Delete(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}