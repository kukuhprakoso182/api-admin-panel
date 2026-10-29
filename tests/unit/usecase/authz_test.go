package usecase_test

import (
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

// ── Status akun ───────────────────────────────────────────────────────────────

func TestLogin_InactiveUserForbidden(t *testing.T) {
	mockRepo := new(testhelper.MockUserRepository)
	mockJWT := new(testhelper.MockJWTService)

	hashed, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	mockRepo.On("FindByEmail", mock.Anything, "sus@example.com").Return(&entity.User{
		ID: "u1", Email: "sus@example.com", Password: string(hashed), StatusID: "suspended",
	}, nil)

	uc := usecase.NewAuthUsecase(mockRepo, mockJWT)
	_, err := uc.Login(ctx, &usecase.LoginRequest{Email: "sus@example.com", Password: "password123"})

	assert.ErrorIs(t, err, domain.ErrForbidden)
	mockJWT.AssertNotCalled(t, "GenerateTokenPair")
}

func TestLogin_InactiveUserWrongPasswordStaysInvalidCreds(t *testing.T) {
	mockRepo := new(testhelper.MockUserRepository)
	mockJWT := new(testhelper.MockJWTService)

	hashed, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	mockRepo.On("FindByEmail", mock.Anything, "sus@example.com").Return(&entity.User{
		ID: "u1", Email: "sus@example.com", Password: string(hashed), StatusID: "suspended",
	}, nil)

	uc := usecase.NewAuthUsecase(mockRepo, mockJWT)
	_, err := uc.Login(ctx, &usecase.LoginRequest{Email: "sus@example.com", Password: "salah"})

	// status akun tidak boleh bocor ke pihak yang tidak tahu password
	assert.ErrorIs(t, err, domain.ErrInvalidCreds)
}

func TestRefreshToken_InactiveUserUnauthorized(t *testing.T) {
	mockRepo := new(testhelper.MockUserRepository)
	mockJWT := new(testhelper.MockJWTService)

	mockJWT.On("ValidateRefreshToken", "rt").Return(&pkgjwt.Claims{UserID: "u1"}, nil)
	mockRepo.On("FindByID", mock.Anything, "u1").Return(&entity.User{ID: "u1", StatusID: "inactive"}, nil)

	uc := usecase.NewAuthUsecase(mockRepo, mockJWT)
	_, err := uc.RefreshToken(ctx, "rt")

	assert.ErrorIs(t, err, domain.ErrUnauthorized)
	mockJWT.AssertNotCalled(t, "GenerateTokenPair")
}

func TestRefreshToken_Success(t *testing.T) {
	mockRepo := new(testhelper.MockUserRepository)
	mockJWT := new(testhelper.MockJWTService)

	mockJWT.On("ValidateRefreshToken", "rt").Return(&pkgjwt.Claims{UserID: "u1"}, nil)
	mockRepo.On("FindByID", mock.Anything, "u1").Return(&entity.User{
		ID: "u1", Email: "a@example.com", StatusID: "active", Roles: []entity.Role{{ID: "admin"}},
	}, nil)
	mockJWT.On("GenerateTokenPair", "u1", "admin", "a@example.com").Return(&pkgjwt.TokenPair{AccessToken: "a"}, nil)

	uc := usecase.NewAuthUsecase(mockRepo, mockJWT)
	pair, err := uc.RefreshToken(ctx, "rt")

	assert.NoError(t, err)
	assert.Equal(t, "a", pair.AccessToken)
}

// ── Hierarki role (child_roles) ───────────────────────────────────────────────

func newGuardedUserUsecase(repo *testhelper.MockUserRepository) usecase.UserUsecase {
	adminChildren := "admin"
	superChildren := "super_admin,admin"

	roles := new(testhelper.MockRoleRepository)
	roles.On("FindByID", mock.Anything, "admin").
		Return(&entity.Role{ID: "admin", IsActive: true, ChildRoles: &adminChildren}, nil)
	roles.On("FindByID", mock.Anything, "super_admin").
		Return(&entity.Role{ID: "super_admin", IsActive: true, ChildRoles: &superChildren}, nil)

	return usecase.NewUserUsecase(repo, usecase.WithRoleRepository(roles))
}

func adminActor() usecase.Actor {
	return usecase.Actor{UserID: "actor-admin", RoleIDs: []string{"admin"}}
}

func createReq(roleIDs ...string) *usecase.CreateUserRequest {
	return &usecase.CreateUserRequest{
		StatusID: "active", RoleIDs: roleIDs, Name: "X", Email: "x@example.com", Password: "password123",
	}
}

func TestUserUsecase_Create_AdminCannotAssignSuperAdmin(t *testing.T) {
	repo := new(testhelper.MockUserRepository)
	uc := newGuardedUserUsecase(repo)

	_, err := uc.Create(usecase.WithActor(ctx, adminActor()), createReq("super_admin"))

	assert.ErrorIs(t, err, domain.ErrForbidden)
	repo.AssertNotCalled(t, "Create")
}

func TestUserUsecase_Create_AdminCanAssignAdmin(t *testing.T) {
	repo := new(testhelper.MockUserRepository)
	repo.On("Create", mock.Anything, mock.Anything, []string{"admin"}).Return(nil)
	repo.On("FindByID", mock.Anything, mock.Anything).Return(&entity.User{ID: "new"}, nil)
	uc := newGuardedUserUsecase(repo)

	_, err := uc.Create(usecase.WithActor(ctx, adminActor()), createReq("admin"))

	assert.NoError(t, err)
}

func TestUserUsecase_Create_WithoutActorIsDenied(t *testing.T) {
	repo := new(testhelper.MockUserRepository)
	uc := newGuardedUserUsecase(repo)

	_, err := uc.Create(ctx, createReq("admin"))

	assert.ErrorIs(t, err, domain.ErrForbidden)
	repo.AssertNotCalled(t, "Create")
}

func TestUserUsecase_Update_AdminCannotEditSuperAdminUser(t *testing.T) {
	repo := new(testhelper.MockUserRepository)
	repo.On("FindByID", mock.Anything, "boss").Return(&entity.User{
		ID: "boss", Roles: []entity.Role{{ID: "super_admin"}},
	}, nil)
	uc := newGuardedUserUsecase(repo)

	_, err := uc.Update(usecase.WithActor(ctx, adminActor()), "boss", &usecase.UpdateUserRequest{
		StatusID: "active", RoleIDs: []string{"admin"}, Name: "Boss", Email: "boss@example.com", Password: "newpassword123",
	})

	assert.ErrorIs(t, err, domain.ErrForbidden)
	repo.AssertNotCalled(t, "Update")
}

func TestUserUsecase_Delete_CannotDeleteSelf(t *testing.T) {
	repo := new(testhelper.MockUserRepository)
	uc := newGuardedUserUsecase(repo)

	err := uc.Delete(usecase.WithActor(ctx, adminActor()), "actor-admin")

	assert.ErrorIs(t, err, domain.ErrForbidden)
	repo.AssertNotCalled(t, "Delete")
}

func TestUserUsecase_Delete_AdminCannotDeleteSuperAdmin(t *testing.T) {
	repo := new(testhelper.MockUserRepository)
	repo.On("FindByID", mock.Anything, "boss").Return(&entity.User{
		ID: "boss", Roles: []entity.Role{{ID: "super_admin"}},
	}, nil)
	uc := newGuardedUserUsecase(repo)

	err := uc.Delete(usecase.WithActor(ctx, adminActor()), "boss")

	assert.ErrorIs(t, err, domain.ErrForbidden)
	repo.AssertNotCalled(t, "Delete")
}
