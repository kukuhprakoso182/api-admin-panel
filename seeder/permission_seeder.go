package seeder

import (
	"api-admin-panel/internal/domain/entity"

	"gorm.io/gorm"
)

func seedPermissions(db *gorm.DB) error {
	permissions := []entity.Permission{
		{ID: "view", Name: "View", IsActive: true},
		{ID: "create", Name: "Create", IsActive: true},
		{ID: "update", Name: "Update", IsActive: true},
		{ID: "delete", Name: "Delete", IsActive: true},
		{ID: "change_profile", Name: "Change Profile", IsActive: true},
		{ID: "change_password", Name: "Change Password", IsActive: true},
		{ID: "change_role", Name: "Change Role", IsActive: true},
		{ID: "change_status", Name: "Change Status", IsActive: true},
		{ID: "reset_password", Name: "Reset Password", IsActive: true},
	}
	for _, v := range permissions {
		if err := db.Where("id = ?", v.ID).FirstOrCreate(&v).Error; err != nil {
			return err
		}
	}
	return nil
}
