package seeder

import (
	"api-admin-panel/internal/domain/entity"
	"fmt"

	"gorm.io/gorm"
)

func seedRoleMenuPermissions(db *gorm.DB) error {
	// Skip jika sudah ada data
	var count int64
	db.Model(&entity.RoleMenuPermission{}).Count(&count)
	if count > 0 {
		return nil
	}

	// Cari menu berdasarkan nama
	findMenu := func(name string) (uint64, error) {
		var m entity.Menu
		if err := db.Where("name = ?", name).First(&m).Error; err != nil {
			return 0, fmt.Errorf("menu '%s' not found: %w", name, err)
		}
		return m.ID, nil
	}

	dashboard, err := findMenu("Dashboard")
	if err != nil {
		return err
	}
	settings, err := findMenu("Settings")
	if err != nil {
		return err
	}
	role, err := findMenu("Role")
	if err != nil {
		return err
	}
	menu, err := findMenu("Menu")
	if err != nil {
		return err
	}
	permission, err := findMenu("Permission")
	if err != nil {
		return err
	}
	status, err := findMenu("Status")
	if err != nil {
		return err
	}
	icon, err := findMenu("Icon")
	if err != nil {
		return err
	}
	user, err := findMenu("User")
	if err != nil {
		return err
	}
	profile, err := findMenu("Profile")
	if err != nil {
		return err
	}
	assignMenuPermission, err := findMenu("Assign Menu Permission")
	if err != nil {
		return err
	}

	type rmpSeed struct {
		RoleID       string
		MenuID       uint64
		PermissionID string
	}

	seeds := []rmpSeed{
		// ── admin ──────────────────────────────────────────────────────────────
		{RoleID: "admin", MenuID: dashboard, PermissionID: "view"},
		{RoleID: "admin", MenuID: user, PermissionID: "view"},
		{RoleID: "admin", MenuID: user, PermissionID: "create"},
		{RoleID: "admin", MenuID: user, PermissionID: "update"},
		{RoleID: "admin", MenuID: user, PermissionID: "reset_password"},
		{RoleID: "admin", MenuID: user, PermissionID: "change_role"},
		{RoleID: "admin", MenuID: user, PermissionID: "change_status"},
		{RoleID: "admin", MenuID: profile, PermissionID: "view"},
		{RoleID: "admin", MenuID: profile, PermissionID: "change_profile"},
		{RoleID: "admin", MenuID: profile, PermissionID: "change_password"},

		// ── super_admin ────────────────────────────────────────────────────────
		{RoleID: "super_admin", MenuID: dashboard, PermissionID: "view"},
		{RoleID: "super_admin", MenuID: user, PermissionID: "view"},
		{RoleID: "super_admin", MenuID: user, PermissionID: "create"},
		{RoleID: "super_admin", MenuID: user, PermissionID: "update"},
		{RoleID: "super_admin", MenuID: user, PermissionID: "reset_password"},
		{RoleID: "super_admin", MenuID: user, PermissionID: "change_role"},
		{RoleID: "super_admin", MenuID: user, PermissionID: "change_status"},
		{RoleID: "super_admin", MenuID: profile, PermissionID: "view"},
		{RoleID: "super_admin", MenuID: profile, PermissionID: "change_profile"},
		{RoleID: "super_admin", MenuID: profile, PermissionID: "change_password"},
		{RoleID: "super_admin", MenuID: status, PermissionID: "view"},
		{RoleID: "super_admin", MenuID: status, PermissionID: "create"},
		{RoleID: "super_admin", MenuID: status, PermissionID: "update"},
		{RoleID: "super_admin", MenuID: status, PermissionID: "change_status"},
		{RoleID: "super_admin", MenuID: role, PermissionID: "view"},
		{RoleID: "super_admin", MenuID: role, PermissionID: "create"},
		{RoleID: "super_admin", MenuID: role, PermissionID: "update"},
		{RoleID: "super_admin", MenuID: role, PermissionID: "change_status"},
		{RoleID: "super_admin", MenuID: menu, PermissionID: "view"},
		{RoleID: "super_admin", MenuID: menu, PermissionID: "create"},
		{RoleID: "super_admin", MenuID: menu, PermissionID: "update"},
		{RoleID: "super_admin", MenuID: menu, PermissionID: "change_status"},
		{RoleID: "super_admin", MenuID: permission, PermissionID: "view"},
		{RoleID: "super_admin", MenuID: permission, PermissionID: "create"},
		{RoleID: "super_admin", MenuID: permission, PermissionID: "update"},
		{RoleID: "super_admin", MenuID: permission, PermissionID: "change_status"},
		{RoleID: "super_admin", MenuID: icon, PermissionID: "view"},
		{RoleID: "super_admin", MenuID: icon, PermissionID: "create"},
		{RoleID: "super_admin", MenuID: icon, PermissionID: "update"},
		{RoleID: "super_admin", MenuID: icon, PermissionID: "change_status"},
		{RoleID: "super_admin", MenuID: settings, PermissionID: "view"},
		{RoleID: "super_admin", MenuID: assignMenuPermission, PermissionID: "view"},
		{RoleID: "super_admin", MenuID: assignMenuPermission, PermissionID: "update"},
	}

	for _, s := range seeds {
		rmp := entity.RoleMenuPermission{
			RoleID:       s.RoleID,
			MenuID:       s.MenuID,
			PermissionID: s.PermissionID,
		}
		if err := db.Where("role_id = ? AND menu_id = ? AND permission_id = ?",
			s.RoleID, s.MenuID, s.PermissionID).
			FirstOrCreate(&rmp).Error; err != nil {
			return err
		}
	}
	return nil
}
