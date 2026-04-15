package usecase_test

import (
	"testing"

	"api-admin-panel/internal/domain/entity"
	"api-admin-panel/internal/usecase"
	"api-admin-panel/tests/testhelper"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ── Auth Usecase Tests ────────────────────────────────────────────────────────

func TestLogin_Success(t *testing.T) {
	mockRepo := new(testhelper.MockUserRepository)
	mockJWT := new(testhelper.MockJWTService)

	hashed, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &entity.User{
		ID:       "user-123",
		Email:    "admin@example.com",
		Password: string(hashed),
		Roles:    []entity.Role{{ID: "admin"}},
	}

	mockRepo.On("FindByEmail", "admin@example.com").Return(user, nil)
	mockJWT.On("GenerateTokenPair", "user-123", "admin", "admin@example.com").
		Return(map[string]interface{}{"access_token": "token-xxx"}, nil)

	uc := usecase.NewAuthUsecase(mockRepo, mockJWT)
	_, err := uc.Login(&usecase.LoginRequest{
		Email:    "admin@example.com",
		Password: "password123",
	})

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestLogin_InvalidEmail(t *testing.T) {
	mockRepo := new(testhelper.MockUserRepository)
	mockJWT := new(testhelper.MockJWTService)

	mockRepo.On("FindByEmail", "notfound@example.com").
		Return(nil, gorm.ErrRecordNotFound)

	uc := usecase.NewAuthUsecase(mockRepo, mockJWT)
	_, err := uc.Login(&usecase.LoginRequest{
		Email:    "notfound@example.com",
		Password: "password123",
	})

	assert.ErrorIs(t, err, usecase.ErrInvalidCreds)
	mockRepo.AssertExpectations(t)
}

func TestLogin_WrongPassword(t *testing.T) {
	mockRepo := new(testhelper.MockUserRepository)
	mockJWT := new(testhelper.MockJWTService)

	hashed, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
	user := &entity.User{
		ID:       "user-123",
		Email:    "test@example.com",
		Password: string(hashed),
		Roles:    []entity.Role{{ID: "admin"}},
	}

	mockRepo.On("FindByEmail", "test@example.com").Return(user, nil)

	uc := usecase.NewAuthUsecase(mockRepo, mockJWT)
	_, err := uc.Login(&usecase.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	})

	assert.ErrorIs(t, err, usecase.ErrInvalidCreds)
	mockRepo.AssertExpectations(t)
	mockJWT.AssertNotCalled(t, "GenerateTokenPair")
}

// ── User Usecase Tests ────────────────────────────────────────────────────────

func TestUserUsecase_GetByID_NotFound(t *testing.T) {
	mockRepo := new(testhelper.MockUserRepository)

	mockRepo.On("FindByID", "nonexistent-id").
		Return(nil, gorm.ErrRecordNotFound)

	uc := usecase.NewUserUsecase(mockRepo)
	_, err := uc.GetByID("nonexistent-id")

	assert.ErrorIs(t, err, usecase.ErrNotFound)
	mockRepo.AssertExpectations(t)
}

func TestUserUsecase_GetAll(t *testing.T) {
	mockRepo := new(testhelper.MockUserRepository)

	expected := []entity.User{
		{ID: "1", Name: "Alice", Email: "alice@example.com"},
		{ID: "2", Name: "Bob", Email: "bob@example.com"},
	}

	mockRepo.On("FindAll", 1, 10).Return(expected, int64(2), nil)

	uc := usecase.NewUserUsecase(mockRepo)
	users, total, err := uc.GetAll(1, 10)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, users, 2)
	mockRepo.AssertExpectations(t)
}

func TestUserUsecase_Delete_NotFound(t *testing.T) {
	mockRepo := new(testhelper.MockUserRepository)

	mockRepo.On("FindByID", "ghost-id").Return(nil, gorm.ErrRecordNotFound)

	uc := usecase.NewUserUsecase(mockRepo)
	err := uc.Delete("ghost-id")

	assert.ErrorIs(t, err, usecase.ErrNotFound)
	mockRepo.AssertExpectations(t)
}

// ── Icon Usecase Tests ────────────────────────────────────────────────────────

func TestIconUsecase_GetAll(t *testing.T) {
	mockRepo := new(testhelper.MockIconRepository)

	expected := []entity.Icon{
		{ID: "ri-home-line", Name: "home", Section: "remixicon"},
		{ID: "ri-user-line", Name: "user", Section: "remixicon"},
	}

	mockRepo.On("FindAll").Return(expected, nil)

	uc := usecase.NewIconUsecase(mockRepo)
	icons, err := uc.GetAll()

	assert.NoError(t, err)
	assert.Len(t, icons, 2)
	assert.Equal(t, "ri-home-line", icons[0].ID)
	mockRepo.AssertExpectations(t)
}

func TestIconUsecase_GetByID_NotFound(t *testing.T) {
	mockRepo := new(testhelper.MockIconRepository)

	mockRepo.On("FindByID", "ri-ghost").Return(nil, gorm.ErrRecordNotFound)

	uc := usecase.NewIconUsecase(mockRepo)
	_, err := uc.GetByID("ri-ghost")

	assert.ErrorIs(t, err, usecase.ErrNotFound)
	mockRepo.AssertExpectations(t)
}

func TestIconUsecase_Create(t *testing.T) {
	mockRepo := new(testhelper.MockIconRepository)

	mockRepo.On("Create", mock.AnythingOfType("*entity.Icon")).Return(nil)

	uc := usecase.NewIconUsecase(mockRepo)
	icon, err := uc.Create(&usecase.CreateIconRequest{
		ID:      "ri-star-line",
		Name:    "star",
		Section: "remixicon",
	})

	assert.NoError(t, err)
	assert.Equal(t, "ri-star-line", icon.ID)
	assert.True(t, icon.IsActive)
	mockRepo.AssertExpectations(t)
}

// ── Role Usecase Tests ────────────────────────────────────────────────────────

func TestRoleUsecase_GetAll(t *testing.T) {
	mockRepo := new(testhelper.MockRoleRepository)

	expected := []entity.Role{
		{ID: "super_admin", Name: "Super Admin"},
		{ID: "admin", Name: "Admin"},
	}

	mockRepo.On("FindAll").Return(expected, nil)

	uc := usecase.NewRoleUsecase(mockRepo)
	roles, err := uc.GetAll()

	assert.NoError(t, err)
	assert.Len(t, roles, 2)
	mockRepo.AssertExpectations(t)
}
