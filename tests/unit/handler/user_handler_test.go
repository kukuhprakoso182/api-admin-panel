package handler_test

import (
	"net/http"
	"testing"

	"api-admin-panel/internal/delivery/http/handler"
	"api-admin-panel/internal/domain/entity"
	"api-admin-panel/internal/usecase"
	"api-admin-panel/tests/testhelper"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserUsecase
type MockUserUsecase struct {
	mock.Mock
}

func (m *MockUserUsecase) GetAll(page, limit int) ([]entity.User, int64, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]entity.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserUsecase) GetByID(id string) (*entity.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserUsecase) Create(req *usecase.CreateUserRequest) (*entity.User, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserUsecase) Update(id string, req *usecase.UpdateUserRequest) (*entity.User, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserUsecase) Delete(id string) error {
	return m.Called(id).Error(0)
}

func TestUserHandler_GetAll_Success(t *testing.T) {
	mockUsecase := new(MockUserUsecase)

	users := []entity.User{
		{ID: "1", Name: "Alice", Email: "alice@example.com"},
		{ID: "2", Name: "Bob", Email: "bob@example.com"},
	}

	mockUsecase.On("GetAll", 1, 10).Return(users, int64(2), nil)

	h := handler.NewUserHandler(mockUsecase)
	r := testhelper.NewRouter()
	r.GET("/users", h.GetAll)

	req := testhelper.NewRequest(http.MethodGet, "/users?page=1&limit=10", nil)
	w := testhelper.Do(r, req)
	body := testhelper.ParseBody(w)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, true, body["success"])
	assert.NotNil(t, body["pagination"])
	mockUsecase.AssertExpectations(t)
}

func TestUserHandler_GetByID_Success(t *testing.T) {
	mockUsecase := new(MockUserUsecase)

	user := &entity.User{ID: "abc-123", Name: "Alice", Email: "alice@example.com"}
	mockUsecase.On("GetByID", "abc-123").Return(user, nil)

	h := handler.NewUserHandler(mockUsecase)
	r := testhelper.NewRouter()
	r.GET("/users/:id", h.GetByID)

	req := testhelper.NewRequest(http.MethodGet, "/users/abc-123", nil)
	w := testhelper.Do(r, req)
	body := testhelper.ParseBody(w)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, true, body["success"])
	mockUsecase.AssertExpectations(t)
}

func TestUserHandler_GetByID_NotFound(t *testing.T) {
	mockUsecase := new(MockUserUsecase)

	mockUsecase.On("GetByID", "ghost-id").Return(nil, usecase.ErrNotFound)

	h := handler.NewUserHandler(mockUsecase)
	r := testhelper.NewRouter()
	r.GET("/users/:id", h.GetByID)

	req := testhelper.NewRequest(http.MethodGet, "/users/ghost-id", nil)
	w := testhelper.Do(r, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockUsecase.AssertExpectations(t)
}

func TestUserHandler_Create_Success(t *testing.T) {
	mockUsecase := new(MockUserUsecase)

	created := &entity.User{
		ID:    "new-uuid",
		Name:  "New User",
		Email: "newuser@example.com",
	}

	mockUsecase.On("Create", mock.AnythingOfType("*usecase.CreateUserRequest")).Return(created, nil)

	h := handler.NewUserHandler(mockUsecase)
	r := testhelper.NewRouter()
	r.POST("/users", h.Create)

	req := testhelper.NewRequest(http.MethodPost, "/users", map[string]interface{}{
		"name":      "New User",
		"email":     "newuser@example.com",
		"password":  "password123",
		"status_id": "active",
		"role_ids":  []string{"admin"},
	})
	w := testhelper.Do(r, req)
	body := testhelper.ParseBody(w)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, true, body["success"])
	mockUsecase.AssertExpectations(t)
}

func TestUserHandler_Create_ValidationError(t *testing.T) {
	mockUsecase := new(MockUserUsecase)

	h := handler.NewUserHandler(mockUsecase)
	r := testhelper.NewRouter()
	r.POST("/users", h.Create)

	// Email tidak valid, password kurang dari 8 char
	req := testhelper.NewRequest(http.MethodPost, "/users", map[string]interface{}{
		"name":  "Bad User",
		"email": "bukan-email",
	})
	w := testhelper.Do(r, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockUsecase.AssertNotCalled(t, "Create")
}

func TestUserHandler_Delete_Success(t *testing.T) {
	mockUsecase := new(MockUserUsecase)

	mockUsecase.On("Delete", "user-to-delete").Return(nil)

	h := handler.NewUserHandler(mockUsecase)
	r := testhelper.NewRouter()
	r.DELETE("/users/:id", h.Delete)

	req := testhelper.NewRequest(http.MethodDelete, "/users/user-to-delete", nil)
	w := testhelper.Do(r, req)
	body := testhelper.ParseBody(w)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, true, body["success"])
	mockUsecase.AssertExpectations(t)
}

func TestUserHandler_Delete_NotFound(t *testing.T) {
	mockUsecase := new(MockUserUsecase)

	mockUsecase.On("Delete", "ghost-id").Return(usecase.ErrNotFound)

	h := handler.NewUserHandler(mockUsecase)
	r := testhelper.NewRouter()
	r.DELETE("/users/:id", h.Delete)

	req := testhelper.NewRequest(http.MethodDelete, "/users/ghost-id", nil)
	w := testhelper.Do(r, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockUsecase.AssertExpectations(t)
}
