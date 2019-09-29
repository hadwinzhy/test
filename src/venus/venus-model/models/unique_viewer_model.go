package models

import (
	"time"
)

// UniqueViewer ...
type UniqueViewer struct {
	BaseModel
	CustomerID     uint      `gorm:"index" json:"customer_id"` // 除了会员是有的，其他全是0
	EventID        uint      `gorm:"index" json:"event_id"`
	CompanyID      uint      `gorm:"index" json:"company_id"`
	SmUVGroupID    uint      `gorm:"index" json:"uv_group_id"`
	VisitDate      time.Time `gorm:"type:date;index" json:"visit_date"`
	CaptureAt      time.Time `gorm:"type:timestamp with time zone;index" json:"capture_at"`
	LastCaptureAt  time.Time `gorm:"type:timestamp with time zone" json:"last_capture_at"`
	FirstCaptureAt time.Time `gorm:"type:timestamp with time zone" json:"first_capture_at"`
	Customer       Customer
	Event          Event
}

// CustomerEventSerialize ...
func (uniqV *UniqueViewer) CustomerEventSerialize() CustomerEventMallSerializer {
	return uniqV.Customer.CustomerEventMallSerialize(&uniqV.CaptureAt, &uniqV.LastCaptureAt, &uniqV.Event)
}
