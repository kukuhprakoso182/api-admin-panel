package usecase

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"api-admin-panel/internal/domain"
	"api-admin-panel/internal/domain/entity"
	"api-admin-panel/internal/domain/repository"
)

type UserUsecase interface {
	GetAll(ctx context.Context, page, limit int) ([]entity.User, int64, error)
	GetByID(ctx context.Context, id string) (*entity.User, error)
	Create(ctx context.Context, req *CreateUserRequest) (*entity.User, error)
	Update(ctx context.Context, id string, req *UpdateUserRequest) (*entity.User, error)
	Delete(ctx context.Context, id string) error
}

type CreateUserRequest struct {
	StatusID string   `json:"status_id" binding:"required,max=50"`
	RoleIDs  []string `json:"role_ids"  binding:"required,min=1"`
	Name     string   `json:"name"      binding:"required,max=100"`
	Email    string   `json:"email"     binding:"required,email,max=100"`
	Password string   `json:"password"  binding:"required,min=8,max=72"`
}

type UpdateUserRequest struct {
	StatusID string   `json:"status_id" binding:"required,max=50"`
	RoleIDs  []string `json:"role_ids"  binding:"required,min=1"`
	Name     string   `json:"name"      binding:"required,max=100"`
	Email    string   `json:"email"     binding:"required,email,max=100"`
	Password string   `json:"password"  binding:"omitempty,min=8,max=72"`
}

type userUsecase struct {
	repo     repository.UserRepository
	roleRepo repository.RoleRepository
}

// UserUsecaseOption mengonfigurasi userUsecase.
type UserUsecaseOption func(*userUsecase)

// WithRoleRepository mengaktifkan pengecekan hierarki role (kolom roles.child_roles):
// seorang actor hanya boleh membuat/mengubah/menghapus user yang rolenya termasuk
// child_roles miliknya, dan hanya boleh meng-assign role tersebut.
// Tanpa opsi ini pengecekan dilewati (dipakai oleh unit test).
func WithRoleRepository(roleRepo repository.RoleRepository) UserUsecaseOption {
	return func(u *userUsecase) { u.roleRepo = roleRepo }
}

func NewUserUsecase(repo repository.UserRepository, opts ...UserUsecaseOption) UserUsecase {
	u := &userUsecase{repo: repo}
	for _, opt := range opts {
		opt(u)
	}
	return u
}

// allowedRoles menghitung role yang boleh dikelola actor = gabungan child_roles
// dari semua role aktif miliknya.
func (u *userUsecase) allowedRoles(ctx context.Context, actor Actor) (map[string]struct{}, error) {
	allowed := make(map[string]struct{})
	for _, rid := range actor.RoleIDs {
		role, err := u.roleRepo.FindByID(ctx, rid)
		if err != nil {
			if errors.Is(err, domain.ErrNotFound) {
				continue
			}
			return nil, err
		}
		if !role.IsActive || role.ChildRoles == nil {
			continue
		}
		for _, child := range strings.Split(*role.ChildRoles, ",") {
			if child = strings.TrimSpace(child); child != "" {
				allowed[child] = struct{}{}
			}
		}
	}
	return allowed, nil
}

// authorizeRoles memastikan semua roleIDs berada dalam hak kelola actor.
// Fail-closed: request tanpa actor di context ditolak.
func (u *userUsecase) authorizeRoles(ctx context.Context, roleIDs []string) error {
	if u.roleRepo == nil {
		return nil
	}
	actor, ok := ActorFrom(ctx)
	if !ok {
		return domain.ErrForbidden
	}
	allowed, err := u.allowedRoles(ctx, actor)
	if err != nil {
		return err
	}
	for _, rid := range roleIDs {
		if _, ok := allowed[rid]; !ok {
			return domain.ErrForbidden
		}
	}
	return nil
}

func roleIDsOf(user *entity.User) []string {
	ids := make([]string, 0, len(user.Roles))
	for _, r := range user.Roles {
		ids = append(ids, r.ID)
	}
	return ids
}

func (u *userUsecase) GetAll(ctx context.Context, page, limit int) ([]entity.User, int64, error) {
	return u.repo.FindAll(ctx, page, limit)
}

func (u *userUsecase) GetByID(ctx context.Context, id string) (*entity.User, error) {
	return u.repo.FindByID(ctx, id)
}

func (u *userUsecase) Create(ctx context.Context, req *CreateUserRequest) (*entity.User, error) {
	if err := u.authorizeRoles(ctx, req.RoleIDs); err != nil {
		return nil, err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		ID:       uuid.New().String(),
		StatusID: req.StatusID,
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashed),
	}

	if err := u.repo.Create(ctx, user, req.RoleIDs); err != nil {
		return nil, err
	}

	// Reload agar relasi terisi
	return u.repo.FindByID(ctx, user.ID)
}

func (u *userUsecase) Update(ctx context.Context, id string, req *UpdateUserRequest) (*entity.User, error) {
	user, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Actor harus berhak atas user target (role saat ini) DAN role baru yang diminta.
	if err := u.authorizeRoles(ctx, roleIDsOf(user)); err != nil {
		return nil, err
	}
	if err := u.authorizeRoles(ctx, req.RoleIDs); err != nil {
		return nil, err
	}

	user.StatusID = req.StatusID
	user.Name = req.Name
	user.Email = req.Email

	if req.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(hashed)
	}

	if err := u.repo.Update(ctx, user, req.RoleIDs); err != nil {
		return nil, err
	}

	return u.repo.FindByID(ctx, user.ID)
}

func (u *userUsecase) Delete(ctx context.Context, id string) error {
	if u.roleRepo != nil {
		actor, ok := ActorFrom(ctx)
		if !ok {
			return domain.ErrForbidden
		}
		if actor.UserID == id {
			return domain.ErrForbidden // tidak boleh menghapus akun sendiri
		}
		target, err := u.repo.FindByID(ctx, id)
		if err != nil {
			return err
		}
		if err := u.authorizeRoles(ctx, roleIDsOf(target)); err != nil {
			return err
		}
	}
	return u.repo.Delete(ctx, id)
}