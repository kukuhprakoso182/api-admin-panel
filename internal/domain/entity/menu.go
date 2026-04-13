package entity

type Menu struct {
	ID        uint64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name      string `gorm:"column:name;size:100;not null" json:"name"`
	Link      string `gorm:"column:link;size:100;not null;default:'#'" json:"link"`
	LinkAlias string `gorm:"column:link_alias;size:100;not null;default:'#'" json:"link_alias"`
	Icon      string `gorm:"column:icon;size:100;not null;default:'#'" json:"icon"`
	Parent    int64  `gorm:"column:parent;not null;default:0" json:"parent"`
	Order     int64  `gorm:"column:order;not null;default:0" json:"order"`
	IsActive  bool   `gorm:"column:is_active;not null;default:true" json:"is_active"`
}

func (Menu) TableName() string { return "menus" }
