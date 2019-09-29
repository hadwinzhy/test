package models

// GetuiAccount DB model
type GetuiAccount struct {
	BaseModel
	CID        string `gorm:"type:varchar(50)" json:"cid"`
	DeviceType string `gorm:"type:varchar(50)" json:"device_type"`
	AdminID    uint   `gorm:"index" json:"admin_id"`
	AuthToken  string `gorm:"type:varchar(50)" json:"auth_token"`
	Admin      Admin
}

// GetuiTask DB model
type GetuiTask struct {
	BaseModel
	CID    string `gorm:"type:varchar(50)" json:"cid"`
	TaskID string `gorm:"type:varchar(128)" json:"task_id"`
	Info   string `gorm:"text" json:"info"`
}
