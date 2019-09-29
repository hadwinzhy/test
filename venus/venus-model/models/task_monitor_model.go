package models

type TaskMonitor struct {
	BaseModel
	TaskName    string `gorm:"type:varchar(50)" json:"task_name"`
	Result      string `gorm:"type:varchar(10)" json:"result"`
	Information string `json:"information"`
}
