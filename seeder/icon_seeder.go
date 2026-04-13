package seeder

import (
	"api-admin-panel/internal/domain/entity"

	"gorm.io/gorm"
)

func seedIcons(db *gorm.DB) error {
	icons := []entity.Icon{
		{ID: "home", Name: "Home", Section: "sidebar", IsActive: true},
		{ID: "user", Name: "User", Section: "sidebar", IsActive: true},
		{ID: "settings", Name: "Settings", Section: "sidebar", IsActive: true},
		{ID: "shield", Name: "Shield", Section: "sidebar", IsActive: true},
		{ID: "menu", Name: "Menu", Section: "sidebar", IsActive: true},
		{ID: "lock", Name: "Lock", Section: "sidebar", IsActive: true},
	}
	for _, v := range icons {
		if err := db.Where("id = ?", v.ID).FirstOrCreate(&v).Error; err != nil {
			return err
		}
	}
	return nil
}
