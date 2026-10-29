package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"api-admin-panel/config"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

// Tipe token. Dicek saat validasi agar refresh token tidak bisa dipakai
// sebagai access token (dan sebaliknya).
const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

type Claims struct {
	UserID    string `json:"user_id"`
	RoleID    string `json:"role_id"`
	Email     string `json:"email"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

type JWTService interface {
	GenerateTokenPair(userID, roleID, email string) (*TokenPair, error)
	// ValidateToken memvalidasi ACCESS token (dipakai middleware).
	ValidateToken(tokenString string) (*Claims, error)
	// ValidateRefreshToken memvalidasi REFRESH token (dipakai endpoint refresh).
	ValidateRefreshToken(tokenString string) (*Claims, error)
}

type jwtService struct {
	cfg *config.JWTConfig
}

func NewJWTService(cfg *config.JWTConfig) JWTService {
	return &jwtService{cfg: cfg}
}

func (j *jwtService) sign(userID, roleID, email, tokenType string, now time.Time, ttl time.Duration) (string, time.Time, error) {
	expiry := now.Add(ttl)
	claims := Claims{
		UserID:    userID,
		RoleID:    roleID,
		Email:     email,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiry),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	signed, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(j.cfg.Secret))
	if err != nil {
		return "", time.Time{}, err
	}
	return signed, expiry, nil
}

func (j *jwtService) GenerateTokenPair(userID, roleID, email string) (*TokenPair, error) {
	now := time.Now()

	accessToken, accessExpiry, err := j.sign(userID, roleID, email, TokenTypeAccess, now, j.cfg.AccessExpireMinutes)
	if err != nil {
		return nil, err
	}

	refreshToken, _, err := j.sign(userID, roleID, email, TokenTypeRefresh, now, j.cfg.RefreshExpireHours)
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
	return j.parse(tokenString, TokenTypeAccess)
}

func (j *jwtService) ValidateRefreshToken(tokenString string) (*Claims, error) {
	return j.parse(tokenString, TokenTypeRefresh)
}

func (j *jwtService) parse(tokenString, expectedType string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(j.cfg.Secret), nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
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
	if claims.TokenType != expectedType {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
