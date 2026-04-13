package usecase

import (
	"errors"

	"github.com/google/uuid"
	"api-admin-panel/internal/domain/entity"
	"api-admin-panel/internal/domain/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserUsecase interface {
	GetAll(page, limit int) ([]entity.User, int64, error)
	GetByID(id string) (*entity.User, error)
	Create(req *CreateUserRequest) (*entity.User, error)
	Update(id string, req *UpdateUserRequest) (*entity.User, error)
	Delete(id string) error
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

func (u *userUsecase) GetAll(page, limit int) ([]entity.User, int64, error) {
	return u.repo.FindAll(page, limit)
}

func (u *userUsecase) GetByID(id string) (*entity.User, error) {
	user, err := u.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return user, nil
}

func (u *userUsecase) Create(req *CreateUserRequest) (*entity.User, error) {
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

	if err := u.repo.Create(user, req.RoleIDs); err != nil {
		return nil, err
	}

	// Reload agar relasi terisi
	return u.repo.FindByID(user.ID)
}

func (u *userUsecase) Update(id string, req *UpdateUserRequest) (*entity.User, error) {
	user, err := u.GetByID(id)
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

	if err := u.repo.Update(user, req.RoleIDs); err != nil {
		return nil, err
	}

	return u.repo.FindByID(user.ID)
}

func (u *userUsecase) Delete(id string) error {
	if _, err := u.GetByID(id); err != nil {
		return err
	}
	return u.repo.Delete(id)
}
