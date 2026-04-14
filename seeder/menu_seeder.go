package seeder

import (
	"api-admin-panel/internal/domain/entity"

	"gorm.io/gorm"
)

func seedMenus(db *gorm.DB) error {
	// Gunakan FirstOrCreate per nama agar aman dijalankan berulang kali
	// dan tidak bergantung pada jumlah record yang ada
	createMenu := func(menu *entity.Menu) error {
		return db.Where("name = ?", menu.Name).FirstOrCreate(menu).Error
	}

	// ── Root menus ────────────────────────────────────────────────────────────
	dashboard := entity.Menu{Name: "Dashboard", Icon: "ri-home-6-line", Link: "dashboard", LinkAlias: "dashboard", Parent: 0, Order: 1, IsActive: true}
	if err := createMenu(&dashboard); err != nil {
		return err
	}

	user := entity.Menu{Name: "User", Icon: "ri-user-line", Link: "users", LinkAlias: "users", Parent: 0, Order: 2, IsActive: true}
	if err := createMenu(&user); err != nil {
		return err
	}

	profile := entity.Menu{Name: "Profile", Link: "profile", LinkAlias: "profile", Parent: 0, Order: 0, IsActive: true}
	if err := createMenu(&profile); err != nil {
		return err
	}

	settings := entity.Menu{Name: "Settings", Icon: "ri-tools-line", Parent: 0, Order: 100, IsActive: true}
	if err := createMenu(&settings); err != nil {
		return err
	}

	// ── Sub menus Settings ────────────────────────────────────────────────────
	subSettings := []entity.Menu{
		{Name: "Icon", Link: "settings/icons", LinkAlias: "icons", Parent: int64(settings.ID), Order: 1, IsActive: true},
		{Name: "Status", Link: "settings/statuses", LinkAlias: "statuses", Parent: int64(settings.ID), Order: 2, IsActive: true},
		{Name: "Role", Link: "settings/roles", LinkAlias: "roles", Parent: int64(settings.ID), Order: 3, IsActive: true},
		{Name: "Permission", Link: "settings/permissions", LinkAlias: "permissions", Parent: int64(settings.ID), Order: 4, IsActive: true},
		{Name: "Menu", Link: "settings/menus", LinkAlias: "menus", Parent: int64(settings.ID), Order: 5, IsActive: true},
	}

	for i := range subSettings {
		if err := createMenu(&subSettings[i]); err != nil {
			return err
		}
	}

	// ── Sub menu Menu ─────────────────────────────────────────────────────────
	var menuRecord entity.Menu
	if err := db.Where("name = ?", "Menu").First(&menuRecord).Error; err != nil {
		return err
	}

	assignMenuPermission := entity.Menu{
		Name:      "Assign Menu Permission",
		Link:      "settings/menus/assign-menu-permissions",
		LinkAlias: "assign-menu-permissions",
		Parent:    int64(menuRecord.ID),
		Order:     0,
		IsActive:  true,
	}
	if err := createMenu(&assignMenuPermission); err != nil {
		return err
	}

	return nil
}
