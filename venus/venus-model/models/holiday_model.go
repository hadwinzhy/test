package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

// holiday no associations
type Holiday struct {
	gorm.Model
	EName      string `gorm:"type:varchar(32);column:en_name"`
	CName      string `gorm:"type:varchar(32);column:ch_name"`
	Count      int    `gorm:"type:integer" json:"count"`
	StartYear  int    `gorm:"type:integer;column:start_year"`
	StartMonth int    `gorm:"type:integer;column:start_month"`
	StartDay   int    `gorm:"type:integer;column:start_day"`
	EndYear    int    `gorm:"type:integer;column:end_year"`
	EndMonth   int    `gorm:"type:integer;column:end_month"`
	EndDay     int    `gorm:"type:integer;column:end_day"`
}

// activity with associations company
type Activity struct {
	BaseModel
	AName          string `gorm:"type:varchar(32);column:activity_name"`
	StartYear      int    `gorm:"type:integer;column:start_year"`
	StartMonth     int    `gorm:"type:varchar" json:"start_month"`
	StartDay       int    `gorm:"type:integer" json:"start_day"`
	EndYear        int    `gorm:"type:integer" json:"end_year"`
	EndMonth       int    `gorm:"type:varchar" json:"end_month"`
	EndDay         int    `gorm:"type:integer" json:"end_day"`
	Count          int    `gorm:"type:integer" json:"count"`
	DeleteOk       uint   `gorm:"default:0"`
	CompanyID      uint
	ShopID         uint
	CompanyGroupID uint
}

type HolidayBasicSerializer struct {
	ID        uint   `json:"id"`
	CName     string `json:"name"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Count     int    `json:"count_date"`
}

func (h *Holiday) HolidaySerializer() HolidayBasicSerializer {
	return HolidayBasicSerializer{
		ID:        h.ID,
		CName:     h.CName,
		StartDate: fmt.Sprintf("%d/%d/%d", h.StartYear, int(h.StartMonth), h.StartDay),
		EndDate:   fmt.Sprintf("%d/%d/%d", h.EndYear, int(h.EndMonth), h.EndDay),
		Count:     h.Count,
	}
}

type ActivityBasicSerializer struct {
	ID             uint   `json:"id"`
	AName          string `json:"name"`
	StartDate      string `json:"start_date"`
	EndDate        string `json:"end_date"`
	CountDate      int    `json:"count_date"`
	CompanyID      uint   `json:"company_id"`
	ShopID         uint   `json:"shop_id"`
	ShopName       string `json:"shop_name"`
	CompanyGroupID uint   `json:"company_group_id"`
}

func (a *Activity) ActivitySerializer() ActivityBasicSerializer {
	return ActivityBasicSerializer{
		ID:             a.ID,
		AName:          a.AName,
		StartDate:      fmt.Sprintf("%d/%d/%d", a.StartYear, int(a.StartMonth), a.StartDay),
		EndDate:        fmt.Sprintf("%d/%d/%d", a.EndYear, int(a.EndMonth), a.EndDay),
		CountDate:      a.Count,
		CompanyID:      a.CompanyID,
		ShopID:         a.ShopID,
		CompanyGroupID: a.CompanyGroupID,
	}
}
