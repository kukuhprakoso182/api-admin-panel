package usecase_test

import (
	"context"
	"testing"

	"api-admin-panel/internal/domain"
	"api-admin-panel/internal/domain/entity"
	"api-admin-panel/internal/usecase"
	pkgjwt "api-admin-panel/pkg/jwt"
	"api-admin-panel/tests/testhelper"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

var ctx = context.Background()

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

	mockRepo.On("FindByEmail", mock.Anything, "admin@example.com").Return(user, nil)
	mockJWT.On("GenerateTokenPair", "user-123", "admin", "admin@example.com").
		Return(&pkgjwt.TokenPair{
			AccessToken:  "access-token-xxx",
			RefreshToken: "refresh-token-xxx",
			ExpiresAt:    9999999999,
		}, nil)

	uc := usecase.NewAuthUsecase(mockRepo, mockJWT)
	result, err := uc.Login(ctx, &usecase.LoginRequest{
		Email:    "admin@example.com",
		Password: "password123",
	})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "access-token-xxx", result.AccessToken)
	mockRepo.AssertExpectations(t)
	mockJWT.AssertExpectations(t)
}

func TestLogin_InvalidEmail(t *testing.T) {
	mockRepo := new(testhelper.MockUserRepository)
	mockJWT := new(testhelper.MockJWTService)

	mockRepo.On("FindByEmail", mock.Anything, "notfound@example.com").
		Return(nil, domain.ErrNotFound)

	uc := usecase.NewAuthUsecase(mockRepo, mockJWT)
	_, err := uc.Login(ctx, &usecase.LoginRequest{
		Email:    "notfound@example.com",
		Password: "password123",
	})

	// "email tidak ada" harus tidak bisa dibedakan dari "password salah"
	assert.ErrorIs(t, err, usecase.ErrInvalidCreds)
	mockRepo.AssertExpectations(t)
	mockJWT.AssertNotCalled(t, "GenerateTokenPair")
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

	mockRepo.On("FindByEmail", mock.Anything, "test@example.com").Return(user, nil)

	uc := usecase.NewAuthUsecase(mockRepo, mockJWT)
	_, err := uc.Login(ctx, &usecase.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	})

	assert.ErrorIs(t, err, usecase.ErrInvalidCreds)
	mockRepo.AssertExpectations(t)
	mockJWT.AssertNotCalled(t, "GenerateTokenPair")
}

func TestLogin_DatabaseError(t *testing.T) {
	mockRepo := new(testhelper.MockUserRepository)
	mockJWT := new(testhelper.MockJWTService)

	mockRepo.On("FindByEmail", mock.Anything, "admin@example.com").
		Return(nil, assert.AnError)

	uc := usecase.NewAuthUsecase(mockRepo, mockJWT)
	_, err := uc.Login(ctx, &usecase.LoginRequest{
		Email:    "admin@example.com",
		Password: "password123",
	})

	// error database tidak boleh disamarkan sebagai kredensial salah
	assert.ErrorIs(t, err, assert.AnError)
	assert.NotErrorIs(t, err, usecase.ErrInvalidCreds)
	mockJWT.AssertNotCalled(t, "GenerateTokenPair")
}

func TestRefreshToken_InvalidToken(t *testing.T) {
	mockRepo := new(testhelper.MockUserRepository)
	mockJWT := new(testhelper.MockJWTService)

	mockJWT.On("ValidateToken", "bad-token").Return(nil, assert.AnError)

	uc := usecase.NewAuthUsecase(mockRepo, mockJWT)
	_, err := uc.RefreshToken(ctx, "bad-token")

	assert.ErrorIs(t, err, usecase.ErrUnauthorized)
	mockRepo.AssertNotCalled(t, "FindByID")
}

func TestRefreshToken_UserDeleted(t *testing.T) {
	mockRepo := new(testhelper.MockUserRepository)
	mockJWT := new(testhelper.MockJWTService)

	mockJWT.On("ValidateToken", "valid-token").Return(&pkgjwt.Claims{UserID: "ghost"}, nil)
	mockRepo.On("FindByID", mock.Anything, "ghost").Return(nil, domain.ErrNotFound)

	uc := usecase.NewAuthUsecase(mockRepo, mockJWT)
	_, err := uc.RefreshToken(ctx, "valid-token")

	assert.ErrorIs(t, err, usecase.ErrUnauthorized)
	mockJWT.AssertNotCalled(t, "GenerateTokenPair")
}

// ── User Usecase Tests ────────────────────────────────────────────────────────

func TestUserUsecase_GetByID_NotFound(t *testing.T) {
	mockRepo := new(testhelper.MockUserRepository)

	mockRepo.On("FindByID", mock.Anything, "nonexistent-id").
		Return(nil, domain.ErrNotFound)

	uc := usecase.NewUserUsecase(mockRepo)
	_, err := uc.GetByID(ctx, "nonexistent-id")

	assert.ErrorIs(t, err, domain.ErrNotFound)
	mockRepo.AssertExpectations(t)
}

func TestUserUsecase_GetAll(t *testing.T) {
	mockRepo := new(testhelper.MockUserRepository)

	expected := []entity.User{
		{ID: "1", Name: "Alice", Email: "alice@example.com"},
		{ID: "2", Name: "Bob", Email: "bob@example.com"},
	}

	mockRepo.On("FindAll", mock.Anything, 1, 10).Return(expected, int64(2), nil)

	uc := usecase.NewUserUsecase(mockRepo)
	users, total, err := uc.GetAll(ctx, 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, users, 2)
	mockRepo.AssertExpectations(t)
}

func TestUserUsecase_Create_HashesPassword(t *testing.T) {
	mockRepo := new(testhelper.MockUserRepository)

	var saved *entity.User
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.User"), []string{"admin"}).
		Run(func(args mock.Arguments) { saved = args.Get(1).(*entity.User) }).
		Return(nil)
	mockRepo.On("FindByID", mock.Anything, mock.AnythingOfType("string")).
		Return(&entity.User{ID: "x"}, nil)

	uc := usecase.NewUserUsecase(mockRepo)
	_, err := uc.Create(ctx, &usecase.CreateUserRequest{
		StatusID: "active",
		RoleIDs:  []string{"admin"},
		Name:     "Alice",
		Email:    "alice@example.com",
		Password: "password123",
	})

	assert.NoError(t, err)
	assert.NotEqual(t, "password123", saved.Password) // tidak boleh tersimpan plaintext
	assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(saved.Password), []byte("password123")))
	mockRepo.AssertExpectations(t)
}

func TestUserUsecase_Create_DuplicateEmail(t *testing.T) {
	mockRepo := new(testhelper.MockUserRepository)

	mockRepo.On("Create", mock.Anything, mock.Anything, mock.Anything).
		Return(domain.ErrAlreadyExists)

	uc := usecase.NewUserUsecase(mockRepo)
	_, err := uc.Create(ctx, &usecase.CreateUserRequest{
		StatusID: "active",
		RoleIDs:  []string{"admin"},
		Name:     "Alice",
		Email:    "alice@example.com",
		Password: "password123",
	})

	assert.ErrorIs(t, err, domain.ErrAlreadyExists)
	mockRepo.AssertNotCalled(t, "FindByID")
}

func TestUserUsecase_Delete_NotFound(t *testing.T) {
	mockRepo := new(testhelper.MockUserRepository)

	// Usecase tidak lagi memanggil FindByID dulu; repository yang mengembalikan ErrNotFound
	mockRepo.On("Delete", mock.Anything, "ghost-id").Return(domain.ErrNotFound)

	uc := usecase.NewUserUsecase(mockRepo)
	err := uc.Delete(ctx, "ghost-id")

	assert.ErrorIs(t, err, domain.ErrNotFound)
	mockRepo.AssertNotCalled(t, "FindByID")
	mockRepo.AssertExpectations(t)
}

// ── Icon Usecase Tests ────────────────────────────────────────────────────────

func TestIconUsecase_GetAll(t *testing.T) {
	mockRepo := new(testhelper.MockIconRepository)

	expected := []entity.Icon{
		{ID: "ri-home-line", Name: "home", Section: "remixicon"},
		{ID: "ri-user-line", Name: "user", Section: "remixicon"},
	}

	mockRepo.On("FindAll", mock.Anything).Return(expected, nil)

	uc := usecase.NewIconUsecase(mockRepo)
	icons, err := uc.GetAll(ctx)

	assert.NoError(t, err)
	assert.Len(t, icons, 2)
	assert.Equal(t, "ri-home-line", icons[0].ID)
	mockRepo.AssertExpectations(t)
}

func TestIconUsecase_GetByID_NotFound(t *testing.T) {
	mockRepo := new(testhelper.MockIconRepository)

	mockRepo.On("FindByID", mock.Anything, "ri-ghost").Return(nil, domain.ErrNotFound)

	uc := usecase.NewIconUsecase(mockRepo)
	_, err := uc.GetByID(ctx, "ri-ghost")

	assert.ErrorIs(t, err, domain.ErrNotFound)
	mockRepo.AssertExpectations(t)
}

func TestIconUsecase_Create(t *testing.T) {
	mockRepo := new(testhelper.MockIconRepository)

	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Icon")).Return(nil)

	uc := usecase.NewIconUsecase(mockRepo)
	icon, err := uc.Create(ctx, &usecase.CreateIconRequest{
		ID:      "ri-star-line",
		Name:    "star",
		Section: "remixicon",
	})

	assert.NoError(t, err)
	assert.Equal(t, "ri-star-line", icon.ID)
	assert.True(t, icon.IsActive)
	mockRepo.AssertExpectations(t)
}

func TestIconUsecase_Update_NotFound(t *testing.T) {
	mockRepo := new(testhelper.MockIconRepository)

	mockRepo.On("FindByID", mock.Anything, "ri-ghost").Return(nil, domain.ErrNotFound)

	uc := usecase.NewIconUsecase(mockRepo)
	_, err := uc.Update(ctx, "ri-ghost", &usecase.UpdateIconRequest{Name: "x", Section: "y"})

	assert.ErrorIs(t, err, domain.ErrNotFound)
	mockRepo.AssertNotCalled(t, "Update")
}

// ── Role Usecase Tests ────────────────────────────────────────────────────────

func TestRoleUsecase_GetAll(t *testing.T) {
	mockRepo := new(testhelper.MockRoleRepository)

	expected := []entity.Role{
		{ID: "super_admin", Name: "Super Admin"},
		{ID: "admin", Name: "Admin"},
	}

	mockRepo.On("FindAll", mock.Anything).Return(expected, nil)

	uc := usecase.NewRoleUsecase(mockRepo)
	roles, err := uc.GetAll(ctx)

	assert.NoError(t, err)
	assert.Len(t, roles, 2)
	mockRepo.AssertExpectations(t)
}