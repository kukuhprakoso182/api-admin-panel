package testhelper

import (
	"api-admin-panel/internal/domain/entity"
	pkgjwt "api-admin-panel/pkg/jwt"

	"github.com/stretchr/testify/mock"
)

// ── MockUserRepository ────────────────────────────────────────────────────────

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindAll(page, limit int) ([]entity.User, int64, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]entity.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) FindByID(id string) (*entity.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(email string) (*entity.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) Create(user *entity.User, roleIDs []string) error {
	return m.Called(user, roleIDs).Error(0)
}

func (m *MockUserRepository) Update(user *entity.User, roleIDs []string) error {
	return m.Called(user, roleIDs).Error(0)
}

func (m *MockUserRepository) Delete(id string) error {
	return m.Called(id).Error(0)
}

// ── MockIconRepository ────────────────────────────────────────────────────────

type MockIconRepository struct {
	mock.Mock
}

func (m *MockIconRepository) FindAll() ([]entity.Icon, error) {
	args := m.Called()
	return args.Get(0).([]entity.Icon), args.Error(1)
}

func (m *MockIconRepository) FindByID(id string) (*entity.Icon, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Icon), args.Error(1)
}

func (m *MockIconRepository) Create(icon *entity.Icon) error {
	return m.Called(icon).Error(0)
}

func (m *MockIconRepository) Update(icon *entity.Icon) error {
	return m.Called(icon).Error(0)
}

func (m *MockIconRepository) Delete(id string) error {
	return m.Called(id).Error(0)
}

// ── MockRoleRepository ────────────────────────────────────────────────────────

type MockRoleRepository struct {
	mock.Mock
}

func (m *MockRoleRepository) FindAll() ([]entity.Role, error) {
	args := m.Called()
	return args.Get(0).([]entity.Role), args.Error(1)
}

func (m *MockRoleRepository) FindByID(id string) (*entity.Role, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Role), args.Error(1)
}

func (m *MockRoleRepository) Create(role *entity.Role) error {
	return m.Called(role).Error(0)
}

func (m *MockRoleRepository) Update(role *entity.Role) error {
	return m.Called(role).Error(0)
}

func (m *MockRoleRepository) Delete(id string) error {
	return m.Called(id).Error(0)
}

// ── MockJWTService ────────────────────────────────────────────────────────────
// Mengimplementasikan pkgjwt.JWTService interface

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
