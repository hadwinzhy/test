package models

import (
	"encoding/json"
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
)

type Warning struct {
	BaseModel
	CompanyID uint           `json:"company_id"`
	Company   Company        `json:"company"`
	Status    string         `gorm:"type:varchar" json:"status"`
	Settings  postgres.Jsonb `gorm:"type:jsonb;column:settings" json:"settings"`
}

func (Warning) TableName() string {
	return "warnings"
}

type WarningSerializer struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	CompanyID uint      `json:"company_id"`
	Company   Company   `json:"company"`
	UpdatedAt time.Time `json:"updated_at"`
	Settings  Settings  `json:"settings"`
	Status    string    `json:"status"`
}

func (w Warning) Serializer() WarningSerializer {
	return WarningSerializer{
		ID:        w.ID,
		CreatedAt: w.CreatedAt,
		UpdatedAt: w.UpdatedAt,
		Company:   w.Company,
		CompanyID: w.CompanyID,
		Settings:  w.ToSettingsHandle(),
		Status:    w.Status,
	}

}

type Setting struct {
	Status      string  `json:"status"`       // 开关状态: on | off
	Type        int     `json:"type"`         // 七个状态: 商场总出入口...
	CompareTime int     `json:"compare_time"` // 三个状态: 比较时间:昨天Y、上周W、上月M
	Trending    int     `json:"trending"`     // 两个状态: 上升 up , 下降 down
	Threshold   float32 `json:"threshold"`    // 数字: 大于0
}

type Settings []Setting

func (w Warning) ToSettingsHandle() Settings {
	var value Settings
	if err := json.Unmarshal(w.Settings.RawMessage, &value); err != nil {
		return nil
	}
	return value
}

func (w Warning) ToJsonbHandle(settings Settings) postgres.Jsonb {
	var (
		value postgres.Jsonb
		err   error
	)
	value.RawMessage, err = json.Marshal(settings)
	if err != nil {
		return value
	}
	return value
}
