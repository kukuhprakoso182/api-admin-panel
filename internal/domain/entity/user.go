package entity

import "time"

// User memiliki banyak role melalui tabel pivot user_roles.
type User struct {
	ID        string    `gorm:"column:id;primaryKey;type:varchar(36)"   json:"id"`
	StatusID  string    `gorm:"column:status_id;size:50;not null"        json:"status_id"`
	Name      string    `gorm:"column:name;size:100;not null"            json:"name"`
	Email     string    `gorm:"column:email;size:100;uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"column:password;size:100;not null"        json:"-"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"         json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"         json:"updated_at"`
	Status    *Status   `gorm:"foreignKey:StatusID;references:ID"        json:"status,omitempty"`
	Roles     []Role    `gorm:"many2many:user_roles;"                    json:"roles,omitempty"`
}

func (User) TableName() string { return "users" }
