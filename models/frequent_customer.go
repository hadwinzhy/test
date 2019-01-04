package models

import (
	"time"
)

type FrequentCustomerGroup struct {
	BaseModel
	CompanyID uint   `gorm:"type:integer;" json:"company_id"`
	ShopID    uint   `gorm:"type:integer;" json:"shop_id"`
	GroupUUID string `gorm:"type:varchar;" json:"group_uuid"`
}

func (FrequentCustomerGroup) TableName() string {
	return "frequent_customer_group"
}

type FrequentCustomerCount struct {
	BaseModel
	PersonUUID              string    `gorm:"type:varchar" json:"person_uuid"`
	Day                     time.Time `gorm:"type:timestamp with time zone" json:"day"`
	CapturedAt              time.Time `gorm:"type:timestamp with time zone" json:"captured_at"`
	EventVisitCount         int       `gorm:"type:integer" json:"event_visit_count"`
	FrequentCustomerGroupID uint
}

func (FrequentCustomerCount) TableName() string {
	return "frequent_customer_count"
}
