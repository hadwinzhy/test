package models

// Profile belongs to admin
type Profile struct {
	BaseModel
	FullName string `gorm:"type:varchar(50)" json:"full_name"`
	AdminID  uint   `gorm:"index" json: "admin_id"`
}
