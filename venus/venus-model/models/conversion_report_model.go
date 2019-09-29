package models

import (
	"encoding/json"
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
)

type ConversionReport struct {
	BaseModel
	CompanyID            uint           `gorm:"index" json:"company_id"`
	ConversionRuleID     uint           `gorm:"index" json:"conversion_rule_id"`
	Hour                 time.Time      `gorm:"type:timestamp with time zone" json:"date"`
	NewCustomerCount     uint           `json:"new_customer_count"`
	NewVIPCount          uint           `json:"new_vip_count"`
	NewCustomerConverted uint           `json:"new_customer_converted"`
	NewVIPConverted      uint           `json:"new_vip_converted"`
	NewMaleConverted     uint           `json:"new_male_converted"`
	NewFemaleConverted   uint           `json:"new_female_converted"`
	NewMaleAgeMap        postgres.Jsonb `gorm:"jsonb" json:"new_male_age_map"`
	NewFemaleAgeMap      postgres.Jsonb `gorm:"jsonb" json:"new_female_age_map"`
	NewDurationMap       postgres.Jsonb `gorm:"jsonb" json:"new_duration_map"`
	// IsCalculated         bool           `json:"is_calculated"`
}

type ConversionReportSerializer struct {
	ID                   uint      `json:"id"`
	CompanyID            uint      `json:"company_id"`
	ConversionRuleID     uint      `json:"conversion_rule_id"`
	Hour                 time.Time `json:"hour"`
	NewCustomerCount     uint      `json:"new_customer_count"`
	NewVIPCount          uint      `json:"new_vip_count"`
	NewCustomerConverted uint      `json:"new_customer_converted"`
	NewVIPConverted      uint      `json:"new_vip_converted"`
	NewMaleConverted     uint      `json:"new_male_converted"`
	NewFemaleConverted   uint      `json:"new_female_converted"`
	NewMaleAgeMap        Map       `json:"new_male_age_map"`
	NewFemaleAgeMap      Map       `json:"new_female_age_map"`
	NewDurationMap       Map       `json:"new_duration_map"`
}

type ConversionReports []ConversionReport

type Map map[string]int

func (c ConversionReport) ToMap(jsonbValue postgres.Jsonb) Map {
	var values = make(map[string]int)
	if err := json.Unmarshal(jsonbValue.RawMessage, &values); err != nil {
		return nil
	}
	return values
}

func (c ConversionReport) BasicSerialize() ConversionReportSerializer {
	return ConversionReportSerializer{
		ID:                   c.ID,
		CompanyID:            c.CompanyID,
		ConversionRuleID:     c.ConversionRuleID,
		Hour:                 c.Hour,
		NewCustomerCount:     c.NewCustomerCount,
		NewCustomerConverted: c.NewCustomerConverted,
		NewVIPCount:          c.NewVIPCount,
		NewVIPConverted:      c.NewVIPConverted,
		NewFemaleConverted:   c.NewFemaleConverted,
		NewMaleConverted:     c.NewMaleConverted,
		NewMaleAgeMap:        c.ToMap(c.NewMaleAgeMap),
		NewFemaleAgeMap:      c.ToMap(c.NewFemaleAgeMap),
		NewDurationMap:       c.ToMap(c.NewDurationMap),
	}

}
