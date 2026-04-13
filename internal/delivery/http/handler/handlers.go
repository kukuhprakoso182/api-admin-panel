package handler

// Handlers mengumpulkan semua HTTP handler dalam satu struct.
// Digunakan oleh bootstrap untuk DI dan oleh router untuk registrasi route.
type Handlers struct {
	Auth       *AuthHandler
	User       *UserHandler
	Icon       *IconHandler
	Status     *StatusHandler
	Role       *RoleHandler
	Menu       *MenuHandler
	Permission *PermissionHandler
}
