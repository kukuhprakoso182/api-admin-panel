package seeder

import (
	"api-admin-panel/internal/domain/entity"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func seedUsers(db *gorm.DB) error {
	type userSeed struct {
		Name     string
		Email    string
		Password string
		RoleIDs  []string
		StatusID string
	}

	users := []userSeed{
		{
			Name: "Super Administrator", Email: "superadmin@example.com",
			Password: "password123", RoleIDs: []string{"superadmin", "admin"}, StatusID: "active",
		},
		{
			Name: "Administrator", Email: "admin@example.com",
			Password: "password123", RoleIDs: []string{"admin"}, StatusID: "active",
		},
		{
			Name: "Manager Staff", Email: "manager@example.com",
			Password: "password123", RoleIDs: []string{"manager", "staff"}, StatusID: "active",
		},
		{
			Name: "Staff User", Email: "staff@example.com",
			Password: "password123", RoleIDs: []string{"staff"}, StatusID: "active",
		},
	}

	for _, u := range users {
		var existing entity.User
		if err := db.Where("email = ?", u.Email).First(&existing).Error; err == nil {
			continue
		}

		hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		user := entity.User{
			ID:       uuid.New().String(),
			Name:     u.Name,
			Email:    u.Email,
			Password: string(hashed),
			StatusID: u.StatusID,
		}

		if err := db.Create(&user).Error; err != nil {
			return err
		}

		for _, rid := range u.RoleIDs {
			ur := entity.UserRole{UserID: user.ID, RoleID: rid}
			if err := db.Create(&ur).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
