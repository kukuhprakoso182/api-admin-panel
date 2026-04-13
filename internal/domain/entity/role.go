package entity

type Role struct {
	ID         string  `gorm:"column:id;primaryKey;size:50" json:"id"`
	Name       string  `gorm:"column:name;size:50;not null" json:"name"`
	ChildRoles *string `gorm:"column:child_roles;type:text" json:"child_roles,omitempty"`
	IsActive   bool    `gorm:"column:is_active;not null;default:true" json:"is_active"`
}

func (Role) TableName() string { return "roles" }
