package models

import (
	"encoding/json"
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"

	"github.com/jinzhu/gorm"
)

const (
	COMPANY_TYPE_STORE        = "store"
	COMPANY_TYPE_SHOPPINGMALL = "shopping_mall"
	COMPANY_TYPE_COMPANYGROUP = "company_group"
)

// Company DB model
// VirtualEntranceShopID 当type是shopping mall时有，虚拟所有出入口的店，用于总出入口会员记录
// VirtualStoresShopID 当type是shopping mall时有，虚拟所有商铺的店，用于总店铺会员记录
type Company struct {
	BaseModel
	Name       string `gorm:"type:varchar(50);not null;unique_index" json:"name"`
	BrandTitle string `gorm:"type:varchar(50);" json:"brand_title"`
	Type       string `gorm:"type:varchar(20);" json:"type"`

	Shops          []Shop
	Devices        []Device
	CustomerGroups []CustomerGroup
	CompanyConfig  CompanyConfig
	Activity       []Activity

	EntranceUvGroupID uint `gorm:"index" json:"entrance_uv_group_id"`
	EntranceUvGroup   SmUvGroup
	DevicePackUUID    string

	CompanyGroupID uint           `gorm:"index" json:"company_group_id"`
	BankingHours   postgres.Jsonb `gorm:"type:jsonb" json:"banking_hours"`
	CompanyGroup   CompanyGroup
}

// CompanyBasicSerializer is basic return Serializer
type CompanyBasicSerializer struct {
	ID            uint                         `json:"id"`
	Name          string                       `json:"name"`
	BrandTitle    string                       `json:"brand_title"`
	CreatedAt     time.Time                    `json:"created_at"`
	UpdatedAt     time.Time                    `json:"updated_at"`
	CompanyConfig CompanyConfigBasicSerializer `json:"company_config"`
	BankingHours  BankingHours                 `json:"banking_hours"`
}

// BasicSerialize will set model to Serializer
func (c *Company) BasicSerialize() CompanyBasicSerializer {
	return CompanyBasicSerializer{
		ID:            c.ID,
		Name:          c.Name,
		BrandTitle:    c.BrandTitle,
		CreatedAt:     c.CreatedAt,
		UpdatedAt:     c.UpdatedAt,
		CompanyConfig: c.CompanyConfig.BasicSerialize(),
		BankingHours:  c.ToBankingHours(),
	}
}

// CompanyBasicSerializerV2 shopping mall类型需要地址和类型
type CompanyBasicSerializerV2 struct {
	CompanyBasicSerializer
	Type    string `json:"type"`
	Address string `json:"address"`
}

// BasicSerializeV2 ...
func (c *Company) BasicSerializeV2(address string) CompanyBasicSerializerV2 {
	return CompanyBasicSerializerV2{
		CompanyBasicSerializer: c.BasicSerialize(),
		Type:    c.Type,
		Address: address,
	}
}

// AfterCreate will create companyConfig
func (c *Company) AfterCreate(tx *gorm.DB) (err error) {
	if &c.CompanyConfig == nil {
		config := CompanyConfig{
			CompanyID: c.ID,
		}
		tx.Save(&config)
	}

	// 创建商场出入口的去重组
	if c.Type == COMPANY_TYPE_SHOPPINGMALL {
		uvGroup := SmUvGroup{
			CompanyID:  c.ID,
			RegionType: UVGROUP_REGION_TYPE_COMPANYENTRANCE,
			RelatedID:  0,
		}

		err = tx.Save(&uvGroup).Error
		if err != nil {
			return
		}

		err = tx.Save(c).Error
	}

	return
}

// ManuallyLoadConfig ...
func (c *Company) ManuallyLoadConfig(tx *gorm.DB) {
	if c.CompanyConfig.ID != 0 {
		return
	}

	var companyConfig CompanyConfig
	tx.Where("company_id = ?", c.ID).First(&companyConfig)
	if companyConfig.ID != 0 {
		c.CompanyConfig = companyConfig
	}
}

type BankingHours struct {
	StartHour string `json:"start_hour"`
	EndHour   string `json:"end_hour"`
}

func (c Company) ToJsonb(bkHours BankingHours) postgres.Jsonb {
	var value postgres.Jsonb
	value.RawMessage, _ = json.Marshal(bkHours)
	return value
}

func (c Company) ToBankingHours() BankingHours {
	var value BankingHours
	if c.BankingHours.RawMessage == nil {
		value.StartHour = "00:00"
		value.EndHour = "23:59"
		return value
	}
	json.Unmarshal(c.BankingHours.RawMessage, &value)
	if value.StartHour == "" || value.EndHour == "" {
		value.StartHour = "00:00"
		value.EndHour = "23:59"
	}
	return value

}
