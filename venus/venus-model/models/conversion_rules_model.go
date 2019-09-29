package models

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
)

type ConversionRule struct {
	BaseModel
	CompanyID      uint           `json:"company_id"`
	FromNickName   string         `gorm:"type:varchar; default:''" json:"from_nick_name"`
	FromRegionName string         `gorm:"type:varchar" json:"from_region_name"`
	FromRegionIDs  postgres.Jsonb `gorm:"column:from_region_ids"  json:"from_region_ids"`
	ToNickName     string         `gorm:"type:varchar; default:''" json:"to_nick_name"`
	ToRegionName   string         `gorm:"type:varchar" json:"to_region_name"`
	ToRegionIDs    postgres.Jsonb `gorm:"column:to_region_ids" json:"to_region_ids"`
}

func (ConversionRule) TableName() string {
	return "conversion_rules"
}

type ConversionRegionInfo struct {
	ID   uint   `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
}

type ConversionRegionInfos []ConversionRegionInfo

type ConversionRuleSerializer struct {
	ID             uint      `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	FromRegionInfo string    `json:"from_region_info"`
	ToRegionInfo   string    `json:"to_region_info"`
	CompanyID      uint      `json:"company_id"`
}

// show serializer for conversion rules
func (c ConversionRule) BasicSerializer() ConversionRuleSerializer {
	conversionShowInfo := ConversionRuleSerializer{
		ID:        c.ID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		CompanyID: c.CompanyID,
	}
	if c.FromNickName != "" {
		conversionShowInfo.FromRegionInfo = fmt.Sprintf("%s(%s)", c.FromNickName, c.FromRegionName)
	} else {
		conversionShowInfo.FromRegionInfo = fmt.Sprintf("%s", c.FromRegionName)
	}

	if c.ToNickName != "" {
		conversionShowInfo.ToRegionInfo = fmt.Sprintf("%s(%s)", c.ToNickName, c.ToRegionName)
	} else {
		conversionShowInfo.ToRegionInfo = fmt.Sprintf("%s", c.ToRegionName)
	}
	return conversionShowInfo
}
