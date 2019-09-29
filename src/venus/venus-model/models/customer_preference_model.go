package models

import (
	"time"
)

type CustomerPreference struct {
	BaseModel
	Day              time.Time `gorm:"type:timestamp with time zone" json:"day"`
	Gender           int       `gorm:"type:integer" json:"gender"`
	AgeInterval      string    `gorm:"type:varchar" json:"age_interval"`
	Count            int       `gorm:"type:integer" json:"count"`
	SmShopID         uint      `gorm:"index"`
	SmBusinessTypeID uint      `gorm:"index"`
	CompanyID        uint      `gorm:"index"`
}

func (CustomerPreference) TableName() string {
	return "customer_preferences"
}

type CustomerPreferences []CustomerPreference
