package testhelper

import (
	"context"

	"api-admin-panel/internal/domain/entity"
	pkgjwt "api-admin-panel/pkg/jwt"

	"github.com/stretchr/testify/mock"
)

// ── MockUserRepository ────────────────────────────────────────────────────────

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindAll(ctx context.Context, page, limit int) ([]entity.User, int64, error) {
	args := m.Called(ctx, page, limit)
	users, _ := args.Get(0).([]entity.User)
	total, _ := args.Get(1).(int64)
	return users, total, args.Error(2)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) Create(ctx context.Context, user *entity.User, roleIDs []string) error {
	return m.Called(ctx, user, roleIDs).Error(0)
}

func (m *MockUserRepository) Update(ctx context.Context, user *entity.User, roleIDs []string) error {
	return m.Called(ctx, user, roleIDs).Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

// ── MockIconRepository ────────────────────────────────────────────────────────

type MockIconRepository struct {
	mock.Mock
}

func (m *MockIconRepository) FindAll(ctx context.Context) ([]entity.Icon, error) {
	args := m.Called(ctx)
	icons, _ := args.Get(0).([]entity.Icon)
	return icons, args.Error(1)
}

func (m *MockIconRepository) FindByID(ctx context.Context, id string) (*entity.Icon, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Icon), args.Error(1)
}

func (m *MockIconRepository) Create(ctx context.Context, icon *entity.Icon) error {
	return m.Called(ctx, icon).Error(0)
}

func (m *MockIconRepository) Update(ctx context.Context, icon *entity.Icon) error {
	return m.Called(ctx, icon).Error(0)
}

func (m *MockIconRepository) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

// ── MockRoleRepository ────────────────────────────────────────────────────────

type MockRoleRepository struct {
	mock.Mock
}

func (m *MockRoleRepository) FindAll(ctx context.Context) ([]entity.Role, error) {
	args := m.Called(ctx)
	roles, _ := args.Get(0).([]entity.Role)
	return roles, args.Error(1)
}

func (m *MockRoleRepository) FindByID(ctx context.Context, id string) (*entity.Role, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Role), args.Error(1)
}

func (m *MockRoleRepository) Create(ctx context.Context, role *entity.Role) error {
	return m.Called(ctx, role).Error(0)
}

func (m *MockRoleRepository) Update(ctx context.Context, role *entity.Role) error {
	return m.Called(ctx, role).Error(0)
}

func (m *MockRoleRepository) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

// ── MockJWTService ────────────────────────────────────────────────────────────
// Mengimplementasikan pkgjwt.JWTService interface (tidak berubah, tanpa ctx)

type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) GenerateTokenPair(userID, roleID, email string) (*pkgjwt.TokenPair, error) {
	args := m.Called(userID, roleID, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pkgjwt.TokenPair), args.Error(1)
}

func (m *MockJWTService) ValidateToken(tokenString string) (*pkgjwt.Claims, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pkgjwt.Claims), args.Error(1)
}