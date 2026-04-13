package entity

type Icon struct {
	ID       string `gorm:"column:id;primaryKey;size:50" json:"id"`
	Name     string `gorm:"column:name;size:50;not null"  json:"name"`
	Section  string `gorm:"column:section;size:50;not null" json:"section"`
	IsActive bool   `gorm:"column:is_active;not null;default:true" json:"is_active"`
}

func (Icon) TableName() string { return "icons" }
