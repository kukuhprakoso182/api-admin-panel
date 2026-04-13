package seeder

import (
	"api-admin-panel/internal/domain/entity"

	"gorm.io/gorm"
)

func seedStatuses(db *gorm.DB) error {
	statuses := []entity.Status{
		{ID: "active", Name: "Active", IsActive: true},
		{ID: "inactive", Name: "Inactive", IsActive: true},
		{ID: "suspended", Name: "Suspended", IsActive: true},
		{ID: "pending", Name: "Pending", IsActive: true},
	}
	for _, v := range statuses {
		if err := db.Where("id = ?", v.ID).FirstOrCreate(&v).Error; err != nil {
			return err
		}
	}
	return nil
}
