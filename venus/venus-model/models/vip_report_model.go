package models

import (
	"time"
)

type SmReportEvent struct {
}

type VIPReportEvent struct {
	BaseModel
	Hour        time.Time `gorm:"type:timestamp with time zone" json:"hour"`
	CompanyID   uint      `gorm:"index" json:"company_id"`
	MaleCount   uint      `gorm:"type:integer" json:"male_count"`
	FemaleCount uint      `gorm:"type:integer" json:"female_count"`
	QueryKey    string    `gorm:"type:varchar(20)" json:"query_key" `
	QueryValue  string    `gorm:"index;type:varchar(20)" json:"query_value"`

	CustomerGroupID uint `gorm:"index" json:"customer_group_id"`
}

type EnhancedUniqueReportEvent struct {
	BaseModel
	Hour        time.Time `gorm:"type:timestamp with time zone" json:"hour"`
	CompanyID   uint      `gorm:"index" json:"company_id"`
	MaleCount   uint      `gorm:"type:integer" json:"male_count"`
	FemaleCount uint      `gorm:"type:integer" json:"female_count"`
	QueryKey    string    `gorm:"type:varchar(20)" json:"query_key" `
	QueryValue  string    `gorm:"index;type:varchar(20)" json:"query_value"`
}

// QueryKey 和 QueryValue应该做一个双列的索引
