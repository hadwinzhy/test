package models

import (
	"encoding/json"
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
)

const (
	MallType  = "mall"
	StoreType = "store"
)
const (
	SHOP_CATEGORY_ENTRANCES = "entrances"
	SHOP_CATEGORY_KEYAREAS  = "key_area"
)

// Shop ...
type Shop struct {
	BaseModel
	Name               string `gorm:"type:varchar(50);not null" json:"name"`
	Location           string `gorm:"type:varchar(255)" json:"location"`
	ProvinceID         uint   `gorm:"index;" json:"province_id"`
	CityID             uint   `gorm:"index;" json:"city_id"`
	DistrictID         uint   `gorm:"index;" json:"district_id"`
	CompanyID          uint   `gorm:"index;not null" json:"company_id"`
	DeviceCount        uint   `gorm:"type:integer;not null;default:0;" json:"device_count"`
	PandoraDeviceCount uint   `gorm:"type:integer;" json:"pandora_device_count"`
	ShopUUID           string `gorm:"type:varchar(16);" json:"shop_uuid"`
	ShopBuilding       string `gorm:"type:varchar(50);" json:"shop_building"`
	Delta              int    `gorm:"type:integer;default:4;" json:"delta"`
	Floor              []Floor
	ShopType           string `gorm:"type:varchar(64);" json:"shop_type"` // 零售类型 + 区域类型
	District           District
	Province           Province
	City               City
	CustomerGroups     []CustomerGroup

	IsVirtual bool `gorm:"type:bool;default:false;column:is_virtual" json:"is_virtual"` // 是否是虚拟的店铺

	BankingHours postgres.Jsonb `gorm:"type:jsonb" json:"banking_hours"` // 营业时间设置：门店版
}

func (s Shop) ToJsonb(bkHours BankingHours) postgres.Jsonb {
	var value postgres.Jsonb
	value.RawMessage, _ = json.Marshal(bkHours)
	return value
}

func (s Shop) ToBankHours() BankingHours {
	var value BankingHours
	if s.BankingHours.RawMessage == nil {
		value.StartHour = "00:00"
		value.EndHour = "23:59"
		return value
	}
	json.Unmarshal(s.BankingHours.RawMessage, &value)
	return value
}

// ShopBasicSerializer ...
type ShopBasicSerializer struct {
	ID           uint              `json:"id"`
	Name         string            `json:"name"`
	Location     string            `json:"location"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
	ProvinceID   uint              `json:"province_id"`
	CityID       uint              `json:"city_id"`
	DistrictID   uint              `json:"district_id"`
	CompanyID    uint              `json:"company_id"`
	ProvinceName string            `json:"province_name"`
	CityName     string            `json:"city_name"`
	DistrictName string            `json:"district_name"`
	ShopBuilding string            `json:"shop_building"`
	Floor        []FloorSerializer `json:"floor"`
	ShopType     string            `json:"shop_type"`
	Number       int               `json:"number"`
	ShopUUID     string            `json:"shop_uuid"`
	BankingHours BankingHours      `json:"banking_hours"`
}

// RegionBasicSerializer
type RegionBasicSerializer struct {
	ID          uint      `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Name        string    `json:"name"`
	ShopType    string    `json:"shop_type"`
	Location    string    `json:"location"`
	DeviceCount uint      `json:"device_count"`
	CompanyID   uint      `json:"company_id"`
	ShopUUID    string    `json:"shop_uuid"`
	Delta       int       `json:"delta"`
}

// ShopReadsenseSerializer ...
type ShopReadsenseSerializer struct {
	ShopBasicSerializer
	DeviceCount   uint `json:"device_count"`
	CustomerCount uint `json:"customer_count"`
}

type RegionWithNameSerializer struct {
	ShopBasicSerializer
	ChName string `json:"ch_name"`
}

// BasicSerialize ...
func (s *Shop) BasicSerialize(args ...interface{}) ShopBasicSerializer {
	floors := make([]FloorSerializer, len(s.Floor))
	platform := "normal"
	if len(args) != 0 && args[0] == "pandora" {
		platform = "pandora"
	}
	for i, floorObj := range s.Floor {
		floors[i] = floorObj.SerializeInPlatform(platform)
	}
	return ShopBasicSerializer{
		ID:           s.ID,
		Name:         s.Name,
		Location:     s.Location,
		CreatedAt:    s.CreatedAt,
		UpdatedAt:    s.UpdatedAt,
		ProvinceID:   s.ProvinceID,
		CityID:       s.CityID,
		DistrictID:   s.DistrictID,
		CompanyID:    s.CompanyID,
		ProvinceName: s.Province.Name,
		CityName:     s.City.Name,
		DistrictName: s.District.Name,
		ShopBuilding: s.ShopBuilding,
		Floor:        floors,
		ShopType:     s.ShopType,
		ShopUUID:     s.ShopUUID,
		Number:       s.Delta,
		BankingHours: s.ToBankHours(),
	}
}

// ReadsenseSerializer ...
func (s *Shop) ReadsenseSerializer() ShopReadsenseSerializer {
	result := ShopReadsenseSerializer{
		ShopBasicSerializer: s.BasicSerialize("normal"),
		DeviceCount:         s.DeviceCount,
	}
	return result
}

// PandoraEyeSerialize ...
func (s *Shop) PandoraEyeSerialize() ShopReadsenseSerializer {
	result := ShopReadsenseSerializer{
		ShopBasicSerializer: s.BasicSerialize("pandora"),
		DeviceCount:         s.PandoraDeviceCount,
	}
	return result
}

func (s *Shop) ShopMallBasicSerialize() RegionBasicSerializer {
	return RegionBasicSerializer{
		ID:          s.ID,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
		Name:        s.Name,
		ShopType:    s.ShopType,
		Location:    s.Location,
		DeviceCount: s.DeviceCount,
		CompanyID:   s.CompanyID,
		ShopUUID:    s.ShopUUID,
		Delta:       s.Delta,
	}
}

// RegionWithName
func (s *Shop) RegionSerializer(name string) RegionWithNameSerializer {
	return RegionWithNameSerializer{
		ShopBasicSerializer: s.BasicSerialize(),
		ChName:              name,
	}

}

// Floor ...
type Floor struct {
	ID                 uint   `gorm:"primary_key;" json:"floor_id"`
	Name               string `gorm:"type:varchar(64); " json:"floor_name"`
	ShopID             uint
	DeviceCount        uint `gorm:"type:integer;" json:"floor_device_count"`
	PandoraDeviceCount uint `gorm:"type:integer;" json:"pandora_device_count"`
}

// FloorSerializer ...
type FloorSerializer struct {
	FloorID          uint   `json:"floor_id"`
	FloorName        string `json:"floor_name"`
	FloorDeviceCount uint   `json:"floor_device_count"`
}

// SerializeInPlatform ...
func (f *Floor) SerializeInPlatform(platform string) FloorSerializer {
	deviceCount := f.DeviceCount
	if platform == "pandora" {
		deviceCount = f.PandoraDeviceCount
	}
	return FloorSerializer{
		FloorID:          f.ID,
		FloorName:        f.Name,
		FloorDeviceCount: deviceCount,
	}
}

//////////////////////////////////////////////////
