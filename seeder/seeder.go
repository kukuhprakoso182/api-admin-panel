package seeder

import (
	"log"

	"gorm.io/gorm"
)

// Run menjalankan semua seeder secara berurutan.
// Urutan penting — tabel dengan foreign key harus di-seed setelah tabel induknya.
func Run(db *gorm.DB) error {
	steps := []struct {
		name string
		fn   func(*gorm.DB) error
	}{
		{"Icons", seedIcons},
		{"Statuses", seedStatuses},
		{"Roles", seedRoles},
		{"Permissions", seedPermissions},
		{"Menus", seedMenus},
		{"Users", seedUsers},
		{"Role Menu Permissions", seedRoleMenuPermissions},
	}

	for _, s := range steps {
		if err := s.fn(db); err != nil {
			return err
		}
		log.Printf("  ✓ %s seeded", s.name)
	}

	log.Println("✅  Seeding completed")
	return nil
}
