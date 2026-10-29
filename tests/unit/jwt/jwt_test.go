package jwt_test

import (
	"strings"
	"testing"
	"time"

	"api-admin-panel/config"
	pkgjwt "api-admin-panel/pkg/jwt"

	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testSecret = "0123456789abcdef0123456789abcdef"

func newService() pkgjwt.JWTService {
	return pkgjwt.NewJWTService(&config.JWTConfig{
		Secret:              testSecret,
		AccessExpireMinutes: time.Minute,
		RefreshExpireHours:  time.Hour,
	})
}

func TestAccessTokenValidatesAsAccessOnly(t *testing.T) {
	svc := newService()
	pair, err := svc.GenerateTokenPair("u1", "admin", "a@example.com")
	require.NoError(t, err)

	claims, err := svc.ValidateToken(pair.AccessToken)
	require.NoError(t, err)
	assert.Equal(t, "u1", claims.UserID)
	assert.Equal(t, pkgjwt.TokenTypeAccess, claims.TokenType)

	_, err = svc.ValidateRefreshToken(pair.AccessToken)
	assert.ErrorIs(t, err, pkgjwt.ErrInvalidToken)
}

func TestRefreshTokenCannotBeUsedAsAccessToken(t *testing.T) {
	svc := newService()
	pair, err := svc.GenerateTokenPair("u1", "admin", "a@example.com")
	require.NoError(t, err)

	_, err = svc.ValidateToken(pair.RefreshToken)
	assert.ErrorIs(t, err, pkgjwt.ErrInvalidToken)

	claims, err := svc.ValidateRefreshToken(pair.RefreshToken)
	require.NoError(t, err)
	assert.Equal(t, pkgjwt.TokenTypeRefresh, claims.TokenType)
}

func TestLegacyTokenWithoutTypeIsRejected(t *testing.T) {
	legacy := pkgjwt.Claims{
		UserID: "u1",
		RegisteredClaims: gojwt.RegisteredClaims{
			ExpiresAt: gojwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	signed, err := gojwt.NewWithClaims(gojwt.SigningMethodHS256, legacy).SignedString([]byte(testSecret))
	require.NoError(t, err)

	svc := newService()
	_, err = svc.ValidateToken(signed)
	assert.ErrorIs(t, err, pkgjwt.ErrInvalidToken)
	_, err = svc.ValidateRefreshToken(signed)
	assert.ErrorIs(t, err, pkgjwt.ErrInvalidToken)
}

func TestWrongSigningMethodIsRejected(t *testing.T) {
	claims := pkgjwt.Claims{
		UserID:    "u1",
		TokenType: pkgjwt.TokenTypeAccess,
		RegisteredClaims: gojwt.RegisteredClaims{
			ExpiresAt: gojwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	signed, err := gojwt.NewWithClaims(gojwt.SigningMethodHS512, claims).SignedString([]byte(testSecret))
	require.NoError(t, err)

	_, err = newService().ValidateToken(signed)
	assert.ErrorIs(t, err, pkgjwt.ErrInvalidToken)
}

func TestTamperedTokenIsRejected(t *testing.T) {
	svc := newService()
	pair, err := svc.GenerateTokenPair("u1", "admin", "a@example.com")
	require.NoError(t, err)

	parts := strings.Split(pair.AccessToken, ".")
	require.Len(t, parts, 3)
	tampered := parts[0] + "." + parts[1] + ".AAAA"

	_, err = svc.ValidateToken(tampered)
	assert.ErrorIs(t, err, pkgjwt.ErrInvalidToken)
}

func TestExpiredTokenReportsExpired(t *testing.T) {
	claims := pkgjwt.Claims{
		UserID:    "u1",
		TokenType: pkgjwt.TokenTypeAccess,
		RegisteredClaims: gojwt.RegisteredClaims{
			ExpiresAt: gojwt.NewNumericDate(time.Now().Add(-time.Minute)),
		},
	}
	signed, err := gojwt.NewWithClaims(gojwt.SigningMethodHS256, claims).SignedString([]byte(testSecret))
	require.NoError(t, err)

	_, err = newService().ValidateToken(signed)
	assert.ErrorIs(t, err, pkgjwt.ErrExpiredToken)
}
