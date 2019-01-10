package models

// Event 是一个简化的model
type Event struct {
	BaseModel
	CustomerID   uint   `gorm:"index" json:"customer_id"`
	ShopID       uint   `gorm:"index;not null" json:"shop_id"`
	DeviceID     uint   `gorm:"index;not null" json:"device_id"`
	DeviceName   string `gorm:"type:varchar(50)" json:"device_name"`
	CompanyID    uint   `gorm:"index;not null" json:"company_id"`
	Age          uint   `gorm:"default:0" json:"age"`
	Gender       uint   `json:"gender"` // 0 女， 1男
	OriginalFace string `gorm:"type:varchar(255)" json:"original_face"`
	PersonID     string `gorm:"index;type:varchar(255)" json:"person_id"`
}
