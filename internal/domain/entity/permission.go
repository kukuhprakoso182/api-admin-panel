package entity

type Permission struct {
	ID       string `gorm:"column:id;primaryKey;size:50" json:"id"`
	Name     string `gorm:"column:name;size:50;not null" json:"name"`
	IsActive bool   `gorm:"column:is_active;not null;default:true" json:"is_active"`
}

func (Permission) TableName() string { return "permissions" }
