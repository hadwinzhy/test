package models

import (
	"encoding/json"
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
)

type ConversionTimeLine struct {
	BaseModel
	TimeLine  postgres.Jsonb `gorm:"column:time_line" json:"timeline"`
	CompanyID uint
}

func (ConversionTimeLine) TableName() string {
	return "conversion_timeline"
}

type TimeLineOneRule struct {
	From int `json:"from"`
	To   int `json:"to"`
}

type ConversionTimeLineSerializer struct {
	ID        uint              `json:"id"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	CompanyID uint              `json:"company_id"`
	TimeLine  []TimeLineOneRule `json:"time_line"`
}

func (c ConversionTimeLine) BasicSerializer() ConversionTimeLineSerializer {
	conversionTimeLineInfo := ConversionTimeLineSerializer{
		ID:        c.ID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		CompanyID: c.CompanyID,
	}
	conversionTimeLineInfo.TimeLine = c.GetTimeline()
	return conversionTimeLineInfo
}

func (c ConversionTimeLine) GetTimeline() []TimeLineOneRule {
	var timeLineRules []TimeLineOneRule
	if err := json.Unmarshal(c.TimeLine.RawMessage, &timeLineRules); err != nil {
		return nil
	}
	return timeLineRules
}
