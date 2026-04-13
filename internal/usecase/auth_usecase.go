package usecase

import (
	"errors"
	"strings"

	"api-admin-panel/internal/domain/repository"
	pkgjwt "api-admin-panel/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthUsecase interface {
	Login(req *LoginRequest) (*pkgjwt.TokenPair, error)
	RefreshToken(refreshToken string) (*pkgjwt.TokenPair, error)
}

type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type authUsecase struct {
	userRepo   repository.UserRepository
	jwtService pkgjwt.JWTService
}

func NewAuthUsecase(userRepo repository.UserRepository, jwtService pkgjwt.JWTService) AuthUsecase {
	return &authUsecase{userRepo: userRepo, jwtService: jwtService}
}

func (u *authUsecase) Login(req *LoginRequest) (*pkgjwt.TokenPair, error) {
	user, err := u.userRepo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCreds
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCreds
	}

	// Kumpulkan semua role ID menjadi string terpisah koma
	roleIDs := make([]string, 0, len(user.Roles))
	for _, r := range user.Roles {
		roleIDs = append(roleIDs, r.ID)
	}

	return u.jwtService.GenerateTokenPair(user.ID, strings.Join(roleIDs, ","), user.Email)
}

func (u *authUsecase) RefreshToken(refreshToken string) (*pkgjwt.TokenPair, error) {
	claims, err := u.jwtService.ValidateToken(refreshToken)
	if err != nil {
		return nil, ErrUnauthorized
	}

	user, err := u.userRepo.FindByID(claims.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUnauthorized
		}
		return nil, err
	}

	roleIDs := make([]string, 0, len(user.Roles))
	for _, r := range user.Roles {
		roleIDs = append(roleIDs, r.ID)
	}

	return u.jwtService.GenerateTokenPair(user.ID, strings.Join(roleIDs, ","), user.Email)
}
