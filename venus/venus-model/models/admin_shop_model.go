package models

type AdminShop struct {
	BaseModel
	ShopID  uint `gorm:"index;not null" json:"shop_id"`
	AdminID uint `gorm:"index;not null" json:"admin_id"`
}
