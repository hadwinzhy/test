package models

import "time"

type FunnelCache struct {
	BaseModel
	CompanyID    uint      `gorm:"index"`
	Period       string    `gorm:"type:varchar(30);default:'day'"`
	Date         time.Time `gorm:"type:date"`
	SrcUvGroupID uint      `gorm:"index"`
	DstUvGroupID uint      `gorm:"index"`
	SrcCount     uint      `gorm:"type:integer;"`
	DstCount     uint      `gorm:"type:integer;"`
	SrcUvGroup   SmUvGroup
	DstUvGroup   SmUvGroup
}
