package models

import (
	"time"
)

type UniqueReportEvent struct {
	ID              uint      `gorm:"primary_key" json:"id"`
	Hour            time.Time `gorm:"type:timestamp with time zone;index;" json:"hour"`
	CustomerInCount uint      `gorm:"default:0" json:"customer_in_count"`
	SmUvGroupID     uint      `gorm:"index" json:"uv_group_id"`
	CompanyID       uint      `gorm:"index;not null" json:"company_id"`
	MaleCount       uint      `gorm:"default:0" json:"male_count"`
	FemaleCount     uint      `gorm:"default:0" json:"female_count"`
}
