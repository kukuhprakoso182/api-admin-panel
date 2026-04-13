package seeder

import (
	"api-admin-panel/internal/domain/entity"

	"gorm.io/gorm"
)

func seedMenus(db *gorm.DB) error {
	var count int64
	db.Model(&entity.Menu{}).Count(&count)
	if count > 0 {
		return nil
	}

	menus := []entity.Menu{
		{Name: "Dashboard", Link: "/dashboard", LinkAlias: "/dashboard", Icon: "home", Parent: 0, Order: 1, IsActive: true},
		{Name: "User Management", Link: "#", LinkAlias: "#", Icon: "user", Parent: 0, Order: 2, IsActive: true},
		{Name: "Users", Link: "/users", LinkAlias: "/users", Icon: "user", Parent: 2, Order: 1, IsActive: true},
		{Name: "Roles", Link: "/roles", LinkAlias: "/roles", Icon: "shield", Parent: 2, Order: 2, IsActive: true},
		{Name: "System", Link: "#", LinkAlias: "#", Icon: "settings", Parent: 0, Order: 3, IsActive: true},
		{Name: "Icons", Link: "/icons", LinkAlias: "/icons", Icon: "menu", Parent: 5, Order: 1, IsActive: true},
		{Name: "Permissions", Link: "/permissions", LinkAlias: "/permissions", Icon: "lock", Parent: 5, Order: 2, IsActive: true},
	}
	for _, v := range menus {
		if err := db.Create(&v).Error; err != nil {
			return err
		}
	}
	return nil
}
