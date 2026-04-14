package seeder

import (
	"api-admin-panel/internal/domain/entity"

	"gorm.io/gorm"
)

func seedRoles(db *gorm.DB) error {
	superAdminChildRoles := "super_admin,admin"
	adminChildRoles := "admin"

	roles := []entity.Role{
		{ID: "super_admin", Name: "Super Admin", ChildRoles: &superAdminChildRoles, IsActive: true},
		{ID: "admin", Name: "Admin", ChildRoles: &adminChildRoles, IsActive: true},
	}

	for _, r := range roles {
		if err := db.Where("id = ?", r.ID).FirstOrCreate(&r).Error; err != nil {
			return err
		}
	}
	return nil
}
