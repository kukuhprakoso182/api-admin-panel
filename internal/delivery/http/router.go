package http

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"api-admin-panel/config"
	"api-admin-panel/internal/delivery/http/handler"
	"api-admin-panel/internal/delivery/http/middleware"
	pkgjwt "api-admin-panel/pkg/jwt"
	"api-admin-panel/pkg/logger"
)

// Router bertanggung jawab hanya pada pendaftaran route.
type Router struct {
	engine     *gin.Engine
	log        *logger.Logger
	cfg        *config.Config
	jwtService pkgjwt.JWTService
	handlers   *handler.Handlers
}

func NewRouter(log *logger.Logger, cfg *config.Config, jwtService pkgjwt.JWTService, handlers *handler.Handlers) *Router {
	return &Router{
		engine:     gin.New(),
		log:        log,
		cfg:        cfg,
		jwtService: jwtService,
		handlers:   handlers,
	}
}

func (r *Router) Setup() *gin.Engine {
	_ = r.engine.SetTrustedProxies([]string{"127.0.0.1", "::1"})

	// ── Global middleware (urutan penting) ────────────────────────────────────
	r.engine.Use(gin.Recovery())                            // 1. Recover dari panic
	r.engine.Use(middleware.SecurityHeaders())              // 2. Tambah security headers
	r.engine.Use(middleware.RequestID())                    // 3. Inject X-Request-ID
	r.engine.Use(middleware.Timeout(30 * time.Second))      // 4. Timeout per request
	r.engine.Use(middleware.MaxRequestSize(4 * 1024 * 1024)) // 5. Max body 4MB
	r.engine.Use(r.corsMiddleware())                        // 6. CORS
	r.engine.Use(middleware.HTTPLogger(r.log))              // 7. Log request/response

	// ── Health check ──────────────────────────────────────────────────────────
	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "env": r.cfg.App.Env})
	})

	// ── API v1 ────────────────────────────────────────────────────────────────
	api := r.engine.Group("/api/v1")

	// Public — rate limit ketat di auth endpoint
	auth := api.Group("/auth")
	auth.Use(middleware.RateLimit(10, time.Minute)) // maks 10 request/menit/IP
	{
		auth.POST("/login", r.handlers.Auth.Login)
		auth.POST("/refresh", r.handlers.Auth.RefreshToken)
	}

	// Protected — semua endpoint butuh JWT
	protected := api.Group("/")
	protected.Use(middleware.JWTAuth(r.jwtService))
	protected.Use(middleware.RateLimitByUser(120, time.Minute)) // maks 120 req/menit/user
	{
		r.registerCRUD(protected, "/users", r.handlers.User)
		r.registerCRUD(protected, "/icons", r.handlers.Icon)
		r.registerCRUD(protected, "/statuses", r.handlers.Status)
		r.registerCRUD(protected, "/roles", r.handlers.Role)
		r.registerCRUD(protected, "/menus", r.handlers.Menu)
		r.registerCRUD(protected, "/permissions", r.handlers.Permission)
	}

	return r.engine
}

// corsMiddleware membuat konfigurasi CORS dari config.
func (r *Router) corsMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     r.cfg.CORS.AllowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     r.cfg.CORS.AllowHeaders,
		ExposeHeaders:    []string{"X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}

// CRUDHandler adalah interface yang harus dipenuhi setiap handler
// agar bisa didaftarkan secara generik oleh registerCRUD.
type CRUDHandler interface {
	GetAll(c *gin.Context)
	GetByID(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

func (r *Router) registerCRUD(group *gin.RouterGroup, path string, h CRUDHandler) {
	g := group.Group(path)
	g.GET("", h.GetAll)
	g.GET("/:id", h.GetByID)
	g.POST("", h.Create)
	g.PUT("/:id", h.Update)
	g.DELETE("/:id", h.Delete)
}
