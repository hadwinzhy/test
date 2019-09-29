package models

import "time"

// VipRecord ...
type VipRecord struct {
	BaseModel
	Date            time.Time     `gorm:"type: timestamp with time zone" json:"date"`
	Gender          bool          `json:"gender"` // true是男，false是女
	CustomerID      uint          `gorm:"index" json:"customer_id"`
	ShopID          uint          `gorm:"index" json:"shop_id"`
	SmRegionID      uint          `gorm:"index" json:"sm_region_id"`
	EventID         uint          `gorm:"index" json:"event_id"`
	CompanyID       uint          `gorm:"index" json:"company_id"`
	CustomerGroupID uint          `gorm:"index" json:"customer_group_id"`
	Count           uint          `gorm:"type: integer" json:"count"`
	CaptureAt       time.Time     `gorm:"type: timestamp with time zone" json:"capture_at"`
	LastCaptureAt   time.Time     `gorm:"type: timestamp with time zone" json:"last_capture_at"`
	Customer        Customer      `gorm:"auto_preload"`
	Shop            Shop          `gorm:"auto_preload"`
	Event           Event         `gorm:"auto_preload"`
	CustomerGroup   CustomerGroup `gorm:"auto_preload"`
	SmRegion        SmRegion      `gorm:"auto_preload"`
}
