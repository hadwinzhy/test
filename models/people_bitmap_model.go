package models

import "time"

type FrequentCustomerPeopleBitMap struct {
	BaseModel
	FrequentCustomerGroupID uint      `gorm:"index"`
	PersonID                string    `gorm:"type:varchar(32)"`
	BitMap                  string    `gorm:"type:BIT(32)"`
	Date                    time.Time `gorm:"type:date"`
}
