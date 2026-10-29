package seeder

import (
	"errors"
	"os"

	"api-admin-panel/internal/domain/entity"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// seedPassword mengambil password awal akun seed dari SEED_DEFAULT_PASSWORD.
// Fallback "password123" hanya diizinkan di luar production.
func seedPassword() (string, error) {
	if pw := os.Getenv("SEED_DEFAULT_PASSWORD"); pw != "" {
		return pw, nil
	}
	if os.Getenv("APP_ENV") == "production" {
		return "", errors.New("SEED_DEFAULT_PASSWORD must be set when seeding in production")
	}
	return "password123", nil
}

func seedUsers(db *gorm.DB) error {
	password, err := seedPassword()
	if err != nil {
		return err
	}

	type userSeed struct {
		Name     string
		Email    string
		Password string
		RoleIDs  []string
		StatusID string
	}

	users := []userSeed{
		{
			Name:     "Super Administrator",
			Email:    "superadmin@example.com",
			Password: password,
			RoleIDs:  []string{"super_admin", "admin"},
			StatusID: "active",
		},
		{
			Name:     "Administrator",
			Email:    "admin@example.com",
			Password: password,
			RoleIDs:  []string{"admin"},
			StatusID: "active",
		},
	}

	for _, u := range users {
		// Gunakan Find+Limit agar tidak trigger log "record not found" dari GORM
		var existing entity.User
		db.Where("email = ?", u.Email).Limit(1).Find(&existing)
		if existing.ID != "" {
			// User sudah ada, skip
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
