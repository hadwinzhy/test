package models

import (
	"strings"
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
)

// CompanyConfig ...

const (
	COMPANY_PLUGIN_FULIDICHAN = "FULIDICHAN"
)

// CompanyConfig ...
type CompanyConfig struct {
	BaseModel
	IDShow                bool            `gorm:"default:false" json:"id_show"`
	AllFaceShow           bool            `gorm:"default:false" json:"all_face_show"`
	ShowFrequentCustomers bool            `gorm:"default:false" json:"show_frequent_customers"`
	EventCallback         string          `json:"event_callback"`
	DeviceCallback        string          `json:"device_callback"`
	CompanyID             uint            `gorm:"index;not null" json:"company_id"`
	ContentType           string          `json:"content_type"`
	Headers               postgres.Hstore `json:"headers"`
	Plugins               string          `gorm:"type:varchar(100);" json:"plugins"`
	NormalThreshold       uint            `gorm:"type:int;" json:"normal_threshold"`
	VIPThreshold          uint            `gorm:"type:int;" json:"vip_threshold"`
	TitanHostName         string          `gorm:"type:varchar(20);" json:"titan_host_name"`
}

// CompanyConfigBasicSerializer ...
type CompanyConfigBasicSerializer struct {
	ID                    uint      `json:"id"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
	IDShow                bool      `json:"id_show"`
	AllFaceShow           bool      `json:"all_face_show"`
	ShowFrequentCustomers bool      `json:"show_frequent_customers"`
	EventCallback         string    `json:"event_callback"`
	DeviceCallback        string    `json:"device_callback"`
	CompanyID             uint      `json:"company_id"`
	ContentType           string    `json:"content_type"`
	Plugins               string    `json:"plugins"`
	NormalThreshold       uint      `json:"normal_threshold"`
	VIPThreshold          uint      `json:"vip_threshold"`
	TitanHostName         string    `json:"titan_host_name"` // set two titan host name
}

const (
	TITAN_HOST_NAME_WALLE = "walle"
	TITAN_HOST_NAME_EVA   = "eva"
)

// BasicSerialize ...
func (c *CompanyConfig) BasicSerialize() CompanyConfigBasicSerializer {
	return CompanyConfigBasicSerializer{
		ID:                    c.ID,
		CreatedAt:             c.CreatedAt,
		UpdatedAt:             c.UpdatedAt,
		IDShow:                c.IDShow,
		AllFaceShow:           c.AllFaceShow,
		ShowFrequentCustomers: c.ShowFrequentCustomers,
		EventCallback:         c.EventCallback,
		DeviceCallback:        c.DeviceCallback,
		Plugins:               c.Plugins,
		CompanyID:             c.CompanyID,
	}
}

// CompanyConfigPluginInstalledSerializer ...
type CompanyConfigPluginInstalledSerializer struct {
	CompanyConfigBasicSerializer
	IsPluginInstalled bool `json:"is_plugin_installed"`
}

// PluginInstalledSerialize ...
func (c *CompanyConfig) PluginInstalledSerialize(pluginName string) CompanyConfigPluginInstalledSerializer {
	return CompanyConfigPluginInstalledSerializer{
		CompanyConfigBasicSerializer: c.BasicSerialize(),
		IsPluginInstalled:            strings.Contains(c.Plugins, pluginName),
	}
}
