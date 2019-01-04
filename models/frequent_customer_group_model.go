package models

type FrequentCustomerGroup struct {
	BaseModel
	CompanyID uint   `gorm:"index"`
	ShopID    uint   `gorm:"index"`
	GroupUUID string `gorm:"type:varchar(32);"`
}
