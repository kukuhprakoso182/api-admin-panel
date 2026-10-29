package handler_test

import (
	"context"
	"net/http"
	"testing"

	"api-admin-panel/internal/delivery/http/handler"
	"api-admin-panel/internal/usecase"
	"api-admin-panel/pkg/jwt"
	"api-admin-panel/tests/testhelper"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuthUsecase
type MockAuthUsecase struct {
	mock.Mock
}

func (m *MockAuthUsecase) Login(ctx context.Context, req *usecase.LoginRequest) (*jwt.TokenPair, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jwt.TokenPair), args.Error(1)
}

func (m *MockAuthUsecase) RefreshToken(ctx context.Context, token string) (*jwt.TokenPair, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jwt.TokenPair), args.Error(1)
}

func TestAuthHandler_Login_Success(t *testing.T) {
	mockUsecase := new(MockAuthUsecase)

	tokenPair := &jwt.TokenPair{
		AccessToken:  "access-token-xxx",
		RefreshToken: "refresh-token-xxx",
		ExpiresAt:    9999999999,
	}

	mockUsecase.On("Login", mock.Anything, mock.MatchedBy(func(req *usecase.LoginRequest) bool {
		return req.Email == "admin@example.com" && req.Password == "password123"
	})).Return(tokenPair, nil)

	h := handler.NewAuthHandler(mockUsecase)
	r := testhelper.NewRouter()
	r.POST("/auth/login", h.Login)

	req := testhelper.NewRequest(http.MethodPost, "/auth/login", map[string]string{
		"email":    "admin@example.com",
		"password": "password123",
	})
	w := testhelper.Do(r, req)
	body := testhelper.ParseBody(w)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, true, body["success"])
	assert.Equal(t, "login successful", body["message"])
	assert.NotNil(t, body["data"])
	mockUsecase.AssertExpectations(t)
}

func TestAuthHandler_Login_InvalidCredentials(t *testing.T) {
	mockUsecase := new(MockAuthUsecase)

	mockUsecase.On("Login", mock.Anything, mock.Anything).Return(nil, usecase.ErrInvalidCreds)

	h := handler.NewAuthHandler(mockUsecase)
	r := testhelper.NewRouter()
	r.POST("/auth/login", h.Login)

	req := testhelper.NewRequest(http.MethodPost, "/auth/login", map[string]string{
		"email":    "admin@example.com",
		"password": "wrongpassword",
	})
	w := testhelper.Do(r, req)
	body := testhelper.ParseBody(w)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, false, body["success"])
	mockUsecase.AssertExpectations(t)
}

func TestAuthHandler_Login_UnexpectedError(t *testing.T) {
	mockUsecase := new(MockAuthUsecase)

	mockUsecase.On("Login", mock.Anything, mock.Anything).
		Return(nil, assert.AnError)

	h := handler.NewAuthHandler(mockUsecase)
	r := testhelper.NewRouter()
	r.POST("/auth/login", h.Login)

	req := testhelper.NewRequest(http.MethodPost, "/auth/login", map[string]string{
		"email":    "admin@example.com",
		"password": "password123",
	})
	w := testhelper.Do(r, req)
	body := testhelper.ParseBody(w)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	// detail error internal tidak boleh bocor ke klien
	assert.Equal(t, "internal server error", body["message"])
	mockUsecase.AssertExpectations(t)
}

func TestAuthHandler_Login_ValidationError(t *testing.T) {
	mockUsecase := new(MockAuthUsecase)

	h := handler.NewAuthHandler(mockUsecase)
	r := testhelper.NewRouter()
	r.POST("/auth/login", h.Login)

	// Email tidak valid
	req := testhelper.NewRequest(http.MethodPost, "/auth/login", map[string]string{
		"email":    "bukan-email",
		"password": "password123",
	})
	w := testhelper.Do(r, req)
	body := testhelper.ParseBody(w)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, false, body["success"])
	mockUsecase.AssertNotCalled(t, "Login")
}

func TestAuthHandler_Login_MissingBody(t *testing.T) {
	mockUsecase := new(MockAuthUsecase)

	h := handler.NewAuthHandler(mockUsecase)
	r := testhelper.NewRouter()
	r.POST("/auth/login", h.Login)

	req := testhelper.NewRequest(http.MethodPost, "/auth/login", nil)
	w := testhelper.Do(r, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockUsecase.AssertNotCalled(t, "Login")
}

func TestAuthHandler_RefreshToken_Success(t *testing.T) {
	mockUsecase := new(MockAuthUsecase)

	tokenPair := &jwt.TokenPair{
		AccessToken:  "new-access-token",
		RefreshToken: "new-refresh-token",
		ExpiresAt:    9999999999,
	}

	mockUsecase.On("RefreshToken", mock.Anything, "valid-refresh-token").Return(tokenPair, nil)

	h := handler.NewAuthHandler(mockUsecase)
	r := testhelper.NewRouter()
	r.POST("/auth/refresh", h.RefreshToken)

	req := testhelper.NewRequest(http.MethodPost, "/auth/refresh", map[string]string{
		"refresh_token": "valid-refresh-token",
	})
	w := testhelper.Do(r, req)
	body := testhelper.ParseBody(w)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, true, body["success"])
	mockUsecase.AssertExpectations(t)
}

func TestAuthHandler_RefreshToken_InvalidToken(t *testing.T) {
	mockUsecase := new(MockAuthUsecase)

	mockUsecase.On("RefreshToken", mock.Anything, "expired-token").Return(nil, usecase.ErrUnauthorized)

	h := handler.NewAuthHandler(mockUsecase)
	r := testhelper.NewRouter()
	r.POST("/auth/refresh", h.RefreshToken)

	req := testhelper.NewRequest(http.MethodPost, "/auth/refresh", map[string]string{
		"refresh_token": "expired-token",
	})
	w := testhelper.Do(r, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockUsecase.AssertExpectations(t)
}