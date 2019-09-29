package models

import (
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

const (
	AUTH_PLACE_TYPE_COMPANYIDS    = "company_ids"
	AUTH_PLACE_TYPE_ALL           = "all"
	AUTH_PLACE_TYPE_BUSINESS_TYPE = "business_type"
	AUTH_PLACE_TYPE_BRANDSHOPS    = "brandshops"
)

// ExtendedAuthority 是shopping mall版用的数据，company_group也用
type ExtendedAuthority struct {
	BaseModel
	PlaceType         string        `gorm:"type:varchar(20)" json:"place_type"`
	Shortcut          string        `gorm:"type:varchar(20)" json:"shortcut"`
	SmBusinessTypeIDs pq.Int64Array `gorm:"type:integer[];" json:"sm_business_type_ids"`
	SmShopIDs         pq.Int64Array `gorm:"type:integer[];" json:"sm_shop_ids"`
	CompanyIDs        pq.Int64Array `gorm:"type:integer[]" json:"company_ids"`
	IsReadOnly        bool          `gorm:"type:bool" json:"is_read_only"`
}

type ExtendedAuthoritySerializer struct {
	BaseSerializer
	PlaceType       string  `json:"place_type"`
	Shortcut        string  `json:"shortcut"`
	BusinessTypeIDs []int64 `json:"business_type_ids"`
	ShopIDs         []int64 `json:"shop_ids"`
	CompanyIDs      []int64 `json:"company_ids"`
	IsReadOnly      bool    `json:"is_read_only"`
}

func (authority *ExtendedAuthority) Serialize() ExtendedAuthoritySerializer {
	return ExtendedAuthoritySerializer{
		BaseSerializer:  authority.BaseModel.Serialize(),
		PlaceType:       authority.PlaceType,
		Shortcut:        authority.Shortcut,
		BusinessTypeIDs: []int64(authority.SmBusinessTypeIDs),
		ShopIDs:         []int64(authority.SmShopIDs),
		IsReadOnly:      authority.IsReadOnly,
	}
}

type ExtendedAuthoritySerializerForRole struct {
	BaseSerializer
	ShowName        string  `json:"name"`
	CompanyNames    string  `json:"company_names"`
	PlaceNames      string  `json:"place_names"`
	BusinessTypeIDs []int64 `json:"business_type_ids"`
	ShopIDs         []int64 `json:"shop_ids"`
	CompanyIDs      []int64 `json:"company_ids"`
	IsReadOnly      bool    `json:"is_read_only"`
}

func fetchCompanyName(tx *gorm.DB, companyID uint) string {
	var company Company
	tx.Where("id = ?", companyID).First(&company)
	return company.Name
}

func (authority ExtendedAuthority) SerializeForRole(tx *gorm.DB, companyName string) ExtendedAuthoritySerializerForRole {
	var serializer ExtendedAuthoritySerializerForRole
	serializer.BaseSerializer = authority.BaseModel.Serialize()
	serializer.BusinessTypeIDs = authority.SmBusinessTypeIDs
	serializer.ShopIDs = authority.SmShopIDs
	serializer.CompanyIDs = authority.CompanyIDs
	serializer.IsReadOnly = authority.IsReadOnly

	if authority.PlaceType == AUTH_PLACE_TYPE_COMPANYIDS && len(authority.CompanyIDs) > 0 {
		var companies []Company
		var names []string
		tx.Where("id in (?)", []int64(authority.CompanyIDs)).Find(&companies).Pluck("name", &names)
		serializer.ShowName = strings.Join(names, ",")
		serializer.CompanyNames = strings.Join(names, ",")
		serializer.PlaceNames = "all"
		return serializer
	} else if len(authority.SmBusinessTypeIDs) > 1 {
		serializer.ShowName = "all"
		serializer.CompanyNames = companyName
		serializer.PlaceNames = serializer.ShowName
		return serializer
	} else if len(authority.SmBusinessTypeIDs) == 1 {
		var smBusinessType SmBusinessType
		if dbErr := tx.Where("id = ?", authority.SmBusinessTypeIDs[0]).First(&smBusinessType).Error; dbErr != nil {
			return serializer
		}
		serializer.BusinessTypeIDs = []int64{int64(smBusinessType.ID)}
		serializer.ShowName = smBusinessType.Name
		serializer.CompanyNames = companyName
		serializer.PlaceNames = serializer.ShowName
		return serializer
	} else if len(authority.SmShopIDs) >= 1 {
		var sid []int64
		for _, i := range authority.SmShopIDs {
			sid = append(sid, i)
		}
		var smShops []SmShop
		if dbErr := tx.Where("id in (?)", sid).Find(&smShops).Error; dbErr != nil {
			return serializer
		}
		var names []string
		var shopIDs []int64
		for _, i := range smShops {
			names = append(names, i.Name)
			shopIDs = append(shopIDs, int64(i.ID))
		}
		serializer.ShopIDs = shopIDs
		serializer.ShowName = strings.Join(names, ",")
		serializer.CompanyNames = companyName
		serializer.PlaceNames = serializer.ShowName
		return serializer
	} else {
		serializer.ShowName = "all" // 全部区域
		serializer.CompanyNames = companyName
		serializer.PlaceNames = serializer.ShowName
		return serializer
	}
}
