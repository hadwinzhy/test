package models

import "time"

type EventMark struct {
	BaseModel
	EventID        uint   `gorm:"index" json:"event_id"`
	EventErrorType string `gorm:"varchar" json:"event_error_type"`
	Event          Event
}

type EventMarkSerializer struct {
	ID             uint       `json:"id"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at"`
	EventID        uint       `json:"event_id"`
	EventErrorType string     `json:"event_error_type"`
	//Event          Event
}

func (e *EventMark) BasicSerializer() EventMarkSerializer {
	return EventMarkSerializer{
		ID:             e.ID,
		CreatedAt:      e.CreatedAt,
		DeletedAt:      e.DeletedAt,
		UpdatedAt:      e.UpdatedAt,
		EventID:        e.EventID,
		EventErrorType: e.EventErrorType,
	}
}

type FullEventMarkSerializer struct {
	EventMarkSerializer
	DetailEvent EventBasicSerializer `json:"detail_event"`
}

func (e *EventMark) FullSerializer() FullEventMarkSerializer {
	return FullEventMarkSerializer{
		EventMarkSerializer: e.BasicSerializer(),
		DetailEvent:         e.Event.BasicSerialize(),
	}

}
