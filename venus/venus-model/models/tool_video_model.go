package models

// ToolVideo
type ToolVideo struct {
	BaseModel
	VideoURL string `gorm:"not null" json:"video_url"`
}
