package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"api-admin-panel/config"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken  = errors.New("token has expired")
)

type Claims struct {
	UserID string `json:"user_id"`
	RoleID string `json:"role_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

type JWTService interface {
	GenerateTokenPair(userID, roleID, email string) (*TokenPair, error)
	ValidateToken(tokenString string) (*Claims, error)
}

type jwtService struct {
	cfg *config.JWTConfig
}

func NewJWTService(cfg *config.JWTConfig) JWTService {
	return &jwtService{cfg: cfg}
}

func (j *jwtService) GenerateTokenPair(userID, roleID, email string) (*TokenPair, error) {
	accessExpiry := time.Now().Add(j.cfg.AccessExpireMinutes)

	accessClaims := Claims{
		UserID: userID,
		RoleID: roleID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpiry),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).
		SignedString([]byte(j.cfg.Secret))
	if err != nil {
		return nil, err
	}

	refreshClaims := Claims{
		UserID: userID,
		RoleID: roleID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.cfg.RefreshExpireHours)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).
		SignedString([]byte(j.cfg.Secret))
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    accessExpiry.Unix(),
	}, nil
}

func (j *jwtService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(j.cfg.Secret), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
