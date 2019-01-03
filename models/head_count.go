package models

import "time"

type HistoryAgainCustomerEvents struct {
	BaseModel
	PersonID       uint      `gorm:"type:integer;" json:"person_id"`
	PersonUUID     string    `gorm:"type:varchar" json:"person_uuid"`
	Day            time.Time `gorm:"type:timestamp with time zone" json:"day"`
	TimeEventVisit string    `gorm:"type:bit(5)" json:"time_event_visit"`
}

func (HistoryAgainCustomerEvents) TableName() string {
	return "history_repeat_customer_events"
}
