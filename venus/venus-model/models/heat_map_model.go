package models

import (
	"encoding/json"
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
)

type HeatMap struct {
	BaseModel
	Background string         `gorm:"type:varchar" json:"background"`
	Width      int            `gorm:"type:integer" json:"width"`
	Height     int            `gorm:"type:integer" json:"height"`
	Points     postgres.Jsonb `gorm:"type:jsonb;column:points" json:"points"`
	SmFloorID  uint           `json:"sm_floor_id"`
	CompanyID  uint           `json:"company_id"`
	Company    Company
}

func (HeatMap) TableName() string {
	return "sm_heat_map"
}

type HeatMapSerializer struct {
	ID         uint      `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Background string    `json:"background"`
	Map        Points    `json:"points"`
	SmFloorID  uint      `json:"sm_floor_id"`
	CompanyID  uint      `json:"company_id"`
	Company    Company   `json:"company"`
	Width      int       `json:"width"`
	Height     int       `json:"height"`
}

func (h HeatMap) Serializer() HeatMapSerializer {
	return HeatMapSerializer{
		ID:         h.ID,
		CreatedAt:  h.CreatedAt.Round(time.Second),
		UpdatedAt:  h.UpdatedAt.Round(time.Second),
		Background: h.Background,
		Map:        h.JsonUnmarshal(),
		CompanyID:  h.CompanyID,
		Company:    h.Company,
		SmFloorID:  h.SmFloorID,
		Width:      h.Width,
		Height:     h.Height,
	}
}

func (h HeatMap) JsonUnmarshal() Points {
	var points Points
	if err := json.Unmarshal(h.Points.RawMessage, &points); err != nil {
		return nil
	}
	return points
}

func (h HeatMap) JsonMarshal(points Points) postgres.Jsonb {
	var values postgres.Jsonb
	values.RawMessage, _ = json.Marshal(points)
	return values
}

func (h HeatMap) GetSmRegionIDs() []uint {
	var points Points
	points = h.JsonUnmarshal()
	var ids []uint
	for _, one := range points {
		ids = append(ids, one.SmRegionID)
	}
	return ids
}

type Point struct {
	X          int  `json:"x"`
	Y          int  `json:"y"`
	SmRegionID uint `json:"sm_region_id"`
}

type Points []Point
