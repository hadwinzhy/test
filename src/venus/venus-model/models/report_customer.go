package models

import (
	"time"
)

// ReportCustomer ..
type ReportCustomer struct {
	ID              uint      `gorm:"primary_key" json:"id"`
	Hour            time.Time `gorm:"type:timestamp with time zone" json:"hour"`
	InCount         uint      `gorm:"default:0" json:"in_count"`
	ShopID          uint      `gorm:"index;not null" json:"shop_id"`
	CustomerGroupID uint      `gorm:"index;" json:"customer_group_id"`
	CompanyID       uint      `gorm:"index;not null" json:"company_id"`
	CustomerGroup   CustomerGroup
}

// ReportCustomerSeriazlier ...
type ReportCustomerSeriazlier struct {
	Time              time.Time `json:"time"`
	InCount           uint      `json:"in_count"`
	CustomerGroupName string    `json:"customer_group_name"`
	CustomerGroupID   uint      `json:"customer_group_id"`
	IsVIPGroup        bool      `json:"is_vip_group"`
}

// BaseSerialize ...
func (c *ReportCustomer) BaseSerialize() ReportCustomerSeriazlier {
	return ReportCustomerSeriazlier{
		Time:              c.Hour,
		InCount:           c.InCount,
		CustomerGroupID:   c.CustomerGroupID,
		CustomerGroupName: c.CustomerGroup.Name,
		IsVIPGroup:        c.CustomerGroup.IsVipGroup(),
	}
}
