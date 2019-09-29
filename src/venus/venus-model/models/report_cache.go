package models

import (
	"strconv"
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
)

// this is a query cache
type ReportCache struct {
	BaseModel
	// Type      string `gorm:"type:varchar(10);" json:"type"`
	FromDate  string `gorm:"type:varchar(20);" json:"from_date"`
	ToDate    string `gorm:"type:varchar(20);" json:"to_date"`
	Period    string `gorm:"type:varchar(10);" json:"period"`
	ShopID    uint   `gorm:"index" json:"shop_id"`
	DeviceID  uint   `gorm:"index" json:"device_id"`
	CompanyID uint   `gorm:"index" json:"company_id"`
	RowCount  uint   `gorm:"type:integer;" json:"row_count"`
}

// every reportcacherow is an element of the table
type ReportCacheRow struct {
	BaseModel
	ReportCacheID uint           `gorm:"index" json:"report_cache_id"`
	RowIndex      uint           `gorm:"type:integer" json:"row_index"`
	Time          time.Time      `gorm:"type:timestamp with time zone" json:"time"`
	EventCount    uint           `gorm:"type:integer" json:"event_count"`
	CustomerCount uint           `gorm:"type:integer" json:"customer_count"`
	MaleCount     uint           `gorm:"type:integer" json:"male_count"`
	FemaleCount   uint           `gorm:"type:integer" json:"female_count"`
	VIPCount      uint           `gorm:"type:integer" json:"vip_count"`
	AgeCount      postgres.Jsonb `gorm:"jsonb;" json:"age_count"`
}

func timestampToTime(from string) time.Time {
	i, _ := strconv.ParseInt(from, 10, 64)
	value := time.Unix(i, 0)
	return value
}

var expireInterval = 10

func (rc *ReportCache) IsCacheExpired() bool {
	toTime := timestampToTime(rc.ToDate)
	if rc.UpdatedAt.Before(toTime) {
		if time.Now().After(toTime) { // a day has passed, must be recalculated
			return true
		} else {
			if time.Now().Sub(rc.UpdatedAt).Minutes() > float64(expireInterval) {
				return true
			}
		}
	}
	return false
}
