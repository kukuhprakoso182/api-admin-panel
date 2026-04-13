package seeder

import (
	"api-admin-panel/internal/domain/entity"

	"gorm.io/gorm"
)

func seedRoles(db *gorm.DB) error {
	childRoles := `["staff"]`
	roles := []entity.Role{
		{ID: "superadmin", Name: "Super Admin", IsActive: true},
		{ID: "admin", Name: "Admin", IsActive: true},
		{ID: "manager", Name: "Manager", ChildRoles: &childRoles, IsActive: true},
		{ID: "staff", Name: "Staff", IsActive: true},
	}
	for _, v := range roles {
		if err := db.Where("id = ?", v.ID).FirstOrCreate(&v).Error; err != nil {
			return err
		}
	}
	return nil
}
