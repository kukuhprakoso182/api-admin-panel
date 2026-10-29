package usecase

import (
	"context"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"api-admin-panel/internal/domain"
	"api-admin-panel/internal/domain/entity"
	"api-admin-panel/internal/domain/repository"
	pkgjwt "api-admin-panel/pkg/jwt"
)

type AuthUsecase interface {
	Login(ctx context.Context, req *LoginRequest) (*pkgjwt.TokenPair, error)
	RefreshToken(ctx context.Context, refreshToken string) (*pkgjwt.TokenPair, error)
}

type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// dummyHash dipakai agar waktu respons login tetap sama
// walaupun email tidak terdaftar.
var dummyHash, _ = bcrypt.GenerateFromPassword([]byte("dummy-password"), bcrypt.DefaultCost)

type authUsecase struct {
	userRepo   repository.UserRepository
	jwtService pkgjwt.JWTService
}

func NewAuthUsecase(userRepo repository.UserRepository, jwtService pkgjwt.JWTService) AuthUsecase {
	return &authUsecase{userRepo: userRepo, jwtService: jwtService}
}

func (u *authUsecase) Login(ctx context.Context, req *LoginRequest) (*pkgjwt.TokenPair, error) {
	user, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			// Tetap jalankan bcrypt agar tidak bisa dibedakan dari "password salah"
			_ = bcrypt.CompareHashAndPassword(dummyHash, []byte(req.Password))
			return nil, ErrInvalidCreds
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCreds
	}

	return u.jwtService.GenerateTokenPair(user.ID, roleIDsCSV(user), user.Email)
}

func (u *authUsecase) RefreshToken(ctx context.Context, refreshToken string) (*pkgjwt.TokenPair, error) {
	claims, err := u.jwtService.ValidateToken(refreshToken)
	if err != nil {
		return nil, ErrUnauthorized
	}

	user, err := u.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, ErrUnauthorized // user sudah dihapus
		}
		return nil, err
	}

	// Role diambil ulang dari database, bukan dari token lama,
	// sehingga perubahan role langsung berlaku saat refresh.
	return u.jwtService.GenerateTokenPair(user.ID, roleIDsCSV(user), user.Email)
}

func roleIDsCSV(user *entity.User) string {
	ids := make([]string, 0, len(user.Roles))
	for _, r := range user.Roles {
		ids = append(ids, r.ID)
	}
	return strings.Join(ids, ",")
}