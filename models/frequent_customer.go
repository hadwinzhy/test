package models

import "time"

type FrequentCustomerCount struct {
	BaseModel
	PersonUUID     string    `gorm:"type:varchar" json:"person_uuid"`
	Day            time.Time `gorm:"type:timestamp with time zone" json:"day"`
	CapturedAt     time.Time `gorm:"type:timestamp with time zone" json:"captured_at"`
	TimeEventVisit int       `gorm:"type:bit(5)" json:"time_event_visit"`
}

func (FrequentCustomerCount) TableName() string {
	return "frequent_customer_count"
}
