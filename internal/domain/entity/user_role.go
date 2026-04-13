package entity

// UserRole adalah tabel pivot many-to-many antara User dan Role.
type UserRole struct {
	UserID string `gorm:"column:user_id;primaryKey;type:varchar(36);not null" json:"user_id"`
	RoleID string `gorm:"column:role_id;primaryKey;size:50;not null"          json:"role_id"`
	User   *User  `gorm:"foreignKey:UserID;references:ID"                     json:"user,omitempty"`
	Role   *Role  `gorm:"foreignKey:RoleID;references:ID"                     json:"role,omitempty"`
}

func (UserRole) TableName() string { return "user_roles" }
