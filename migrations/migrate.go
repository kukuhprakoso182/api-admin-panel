package migrations

import (
	"api-admin-panel/internal/domain/entity"

	"gorm.io/gorm"
)

// Run menjalankan auto-migration semua tabel sesuai urutan dependency.
func Run(db *gorm.DB, fresh bool) error {
	tables := []interface{}{
		&entity.Icon{},
		&entity.Status{},
		&entity.Role{},
		&entity.Menu{},
		&entity.Permission{},
		&entity.User{},
		&entity.UserRole{},
		&entity.RoleMenuPermission{},
	}

	if fresh {
		if err := db.Migrator().DropTable(tables...); err != nil {
			return err
		}
	}

	for _, table := range tables {
		if err := db.AutoMigrate(table); err != nil {
			return err
		}
	}
	return nil
}
