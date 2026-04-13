package entity

type RoleMenuPermission struct {
	RoleID       string     `gorm:"column:role_id;primaryKey;size:50;not null" json:"role_id"`
	MenuID       uint64     `gorm:"column:menu_id;primaryKey;not null" json:"menu_id"`
	PermissionID string     `gorm:"column:permission_id;primaryKey;size:50;not null" json:"permission_id"`
	Role         *Role      `gorm:"foreignKey:RoleID;references:ID" json:"role,omitempty"`
	Menu         *Menu      `gorm:"foreignKey:MenuID;references:ID" json:"menu,omitempty"`
	Permission   *Permission `gorm:"foreignKey:PermissionID;references:ID" json:"permission,omitempty"`
}

func (RoleMenuPermission) TableName() string { return "role_menu_permission" }
