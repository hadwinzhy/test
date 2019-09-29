package models

import (
	"strconv"
	"siren/venus/venus-model/models/logger"

	raven "github.com/getsentry/raven-go"
	"github.com/jinzhu/gorm"
)

const (
	SPECIAL_GROUP_TYPE_FULIDICHAN = "FULIDICHAN"
)

// SpecialGroup ...
type SpecialGroup struct {
	BaseModel
	Name      string `gorm:"type:varchar(50);not null" json:"name"`
	CompanyID uint   `gorm:"index;not null" json:"company_id"`
	ShopID    uint   `gorm:"index;not null" json:"shop_id"`
	GroupID   string `gorm:"index" json:"group_id"`
	GroupType string `gorm:"default:'FULIDICHAN'" json:"group_type"`
	Company   Company
	Shop      Shop
}

func MakeNewSpecialGroup(tx *gorm.DB, shop Shop, devices []Device) *SpecialGroup {
	var group SpecialGroup
	var macAddresses []string
	for _, device := range devices {
		macAddresses = append(macAddresses, device.MacAddress)
	}

	groupUUID, err := CreatePandoraGroup(SPECIAL_GROUP_TYPE_FULIDICHAN, macAddresses, SPECIAL_GROUP_TYPE_FULIDICHAN, "")
	if err != nil {
		logger.Error(nil, "customer_group", "create", "CREATE_GROUP_IN_PANDORA_ERROR: ", err)
		raven.CaptureError(err, map[string]string{
			"action":     "CREATE_GROUP_IN_PANDORA_ERROR",
			"company_id": strconv.Itoa(int(shop.CompanyID)),
			"group_name": SPECIAL_GROUP_TYPE_FULIDICHAN,
			"detail":     err.Error(),
		})
		return nil
	}

	group.GroupID = groupUUID
	group.Name = SPECIAL_GROUP_TYPE_FULIDICHAN
	group.GroupType = SPECIAL_GROUP_TYPE_FULIDICHAN
	group.ShopID = shop.ID
	group.CompanyID = shop.CompanyID

	if err := tx.Create(&group).Error; err != nil {
		return nil
	}

	return &group
}
