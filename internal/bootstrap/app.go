package bootstrap

import (
	"github.com/gin-gonic/gin"
	"api-admin-panel/config"
	httpdelivery "api-admin-panel/internal/delivery/http"
	"api-admin-panel/internal/delivery/http/handler"
	domainrepo "api-admin-panel/internal/domain/repository"
	mysqlrepo "api-admin-panel/internal/repository/mysql"
	"api-admin-panel/internal/usecase"
	"api-admin-panel/pkg/database"
	pkgjwt "api-admin-panel/pkg/jwt"
	"api-admin-panel/pkg/logger"
	"gorm.io/gorm"
)

// App menyimpan semua dependency yang sudah di-wire.
type App struct {
	Engine *gin.Engine
	DB     *gorm.DB
	Log    *logger.Logger
	Config *config.Config
}

// New menginisialisasi seluruh dependency dan mengembalikan App yang siap dijalankan.
func New(cfg *config.Config, log *logger.Logger) (*App, error) {
	// ── Gin mode ──────────────────────────────────────────────────────────────
	switch cfg.App.Env {
	case "production", "staging":
		gin.SetMode(gin.ReleaseMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	// ── Database ──────────────────────────────────────────────────────────────
	db, err := database.NewMySQLConnection(cfg)
	if err != nil {
		return nil, err
	}

	// ── Wire semua layer ──────────────────────────────────────────────────────
	jwtService := pkgjwt.NewJWTService(&cfg.JWT)
	repos      := initRepositories(db)
	usecases   := initUsecases(repos, jwtService)
	handlers   := initHandlers(usecases)

	// ── Router ────────────────────────────────────────────────────────────────
	router := httpdelivery.NewRouter(log, cfg, jwtService, handlers, repos.roleMenuPermission)
	engine := router.Setup()

	return &App{
		Engine: engine,
		DB:     db,
		Log:    log,
		Config: cfg,
	}, nil
}

// ── Repositories ──────────────────────────────────────────────────────────────

type repoContainer struct {
	icon               domainrepo.IconRepository
	status             domainrepo.StatusRepository
	role               domainrepo.RoleRepository
	menu               domainrepo.MenuRepository
	permission         domainrepo.PermissionRepository
	user               domainrepo.UserRepository
	roleMenuPermission domainrepo.RoleMenuPermissionRepository
}

func initRepositories(db *gorm.DB) *repoContainer {
	return &repoContainer{
		icon:               mysqlrepo.NewIconRepository(db),
		status:             mysqlrepo.NewStatusRepository(db),
		role:               mysqlrepo.NewRoleRepository(db),
		menu:               mysqlrepo.NewMenuRepository(db),
		permission:         mysqlrepo.NewPermissionRepository(db),
		user:               mysqlrepo.NewUserRepository(db),
		roleMenuPermission: mysqlrepo.NewRoleMenuPermissionRepository(db),
	}
}

// ── Usecases ──────────────────────────────────────────────────────────────────

type usecaseContainer struct {
	auth       usecase.AuthUsecase
	user       usecase.UserUsecase
	icon       usecase.IconUsecase
	status     usecase.StatusUsecase
	role       usecase.RoleUsecase
	menu       usecase.MenuUsecase
	permission usecase.PermissionUsecase
}

func initUsecases(r *repoContainer, jwtService pkgjwt.JWTService) *usecaseContainer {
	return &usecaseContainer{
		auth:       usecase.NewAuthUsecase(r.user, jwtService),
		user:       usecase.NewUserUsecase(r.user, usecase.WithRoleRepository(r.role)),
		icon:       usecase.NewIconUsecase(r.icon),
		status:     usecase.NewStatusUsecase(r.status),
		role:       usecase.NewRoleUsecase(r.role),
		menu:       usecase.NewMenuUsecase(r.menu),
		permission: usecase.NewPermissionUsecase(r.permission),
	}
}

// ── Handlers ──────────────────────────────────────────────────────────────────

func initHandlers(u *usecaseContainer) *handler.Handlers {
	return &handler.Handlers{
		Auth:       handler.NewAuthHandler(u.auth),
		User:       handler.NewUserHandler(u.user),
		Icon:       handler.NewIconHandler(u.icon),
		Status:     handler.NewStatusHandler(u.status),
		Role:       handler.NewRoleHandler(u.role),
		Menu:       handler.NewMenuHandler(u.menu),
		Permission: handler.NewPermissionHandler(u.permission),
	}
}
