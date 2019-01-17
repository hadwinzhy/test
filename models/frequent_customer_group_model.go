package models

import "siren/pkg/database"

type FrequentCustomerGroup struct {
	BaseModel
	CompanyID     uint   `gorm:"index"`
	ShopID        uint   `gorm:"index"`
	GroupUUID     string `gorm:"type:varchar(32);"`
	DefaultNumber uint   `gorm:"type:integer"`
}

// FetchFrequentCustomerGroup获取公司,门店对应的组 因为有可能是门店版的companyID参数，所以返回值是slice
func FetchFrequentCustomerGroup(companyID uint, shopIDs []uint) []FrequentCustomerGroup {
	var results []FrequentCustomerGroup
	query := database.POSTGRES.Where("company_id = ?", companyID)

	if len(shopIDs) != 0 {
		query = query.Where("shop_id in (?)", shopIDs)
	}

	query.Find(&results)
	return results
}

type FrequentCustomerGroups []FrequentCustomerGroup
