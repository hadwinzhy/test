package models

import "time"

type FCPeopleBitMap struct {
	BaseModel
	FrequentCustomerGroupID uint      `gorm:"index"`
	PersonID                string    `gorm:"type:varchar(32)"`
	BitMap                  string    `gorm:"type:BIT(32)"`
	CurrentDate             time.Time `gorm:"type:date"`
}
