package models

import "time"

type FrequentCustomerReport struct {
	BaseModel
	FrequentCustomerGroupID uint      `gorm:"index"`
	Date                    time.Time `gorm:"type:date"`
	HighFrequency           uint      `gorm:"type:integer"`
	LowFrequency            uint      `gorm:"type:integer"`
	NewComer                uint      `gorm:"type:integer"`
	SumInterval             uint      `gorm:"type:integer"`
	SumTimes                uint      `gorm:"type:integer"`
}

// 总人数，高频次数，低频次数，新客数，总到访间隔天数，总到访天数
