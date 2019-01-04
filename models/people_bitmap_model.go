package models

import "time"

type FrequentCustomerPeopleBitMap struct {
	BaseModel
	FrequentCustomerPeopleID uint   `gorm:"index"`
	BitMap                   string `gorm:"type:BIT(32)"`
	FrequentCustomerPeople   FrequentCustomerPeople
}

type FrequentCustomerPeople struct {
	BaseModel
	FrequentCustomerGroupID uint      `gorm:"index"`
	PersonID                string    `gorm:"type:varchar(32)"`
	Date                    time.Time `gorm:"type:date"`
	Interval                uint      `gorm:"type:integer"`
	Frequency               uint      `gorm:"type:integer"`
}
