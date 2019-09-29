package models

const (
	UVGROUP_REGION_TYPE_COMPANYENTRANCE = 1
	UVGROUP_REGION_TYPE_FLOORENTRANCE   = 2
	UVGROUP_REGION_TYPE_FLOORPUBLIC     = 3
	UVGROUP_REGION_TYPE_SHOPENTRANCE    = 4
	UVGROUP_REGION_TYPE_CUSTOMIZE       = 5
	UVGROUP_REGION_TYPE_COMPANYALL      = 6
)

var UVGROUP_TYPENAMEMAP = map[int]string{
	UVGROUP_REGION_TYPE_COMPANYENTRANCE: "商场出入口",
	UVGROUP_REGION_TYPE_FLOORENTRANCE:   "楼层出入口",
	UVGROUP_REGION_TYPE_FLOORPUBLIC:     "公共区域",
	UVGROUP_REGION_TYPE_SHOPENTRANCE:    "商铺出入口",
	UVGROUP_REGION_TYPE_CUSTOMIZE:       "自定义",
}

type SmUvGroup struct {
	BaseModel
	CompanyID  uint       `gorm:"index" json:"company_id"`
	SmRegions  []SmRegion `gorm:"many2many:sm_regions_uvgroups;" json:"sm_regions"`
	RegionType uint       `gorm:"type:integer;" json:"region_type"`
	RelatedID  uint       `gorm:"index" json:"related_id"`
	FloorID    uint       `gorm:"index" json:"floor_id"`
}

type SmUVGroupSimpleSerializer struct {
	RegionType uint `json:"region_type"`
	RelatedID  uint `json:"related_id"`
}

func (group *SmUvGroup) SimpleSerialize() SmUVGroupSimpleSerializer {
	return SmUVGroupSimpleSerializer{
		RegionType: group.RegionType,
		RelatedID:  group.RelatedID,
	}
}

func (group *SmUvGroup) TypeName() string {
	return UVGROUP_TYPENAMEMAP[int(group.RegionType)]
}

// func (s *SmUvGroup) SetDevicePackPandora(devicePackUUID string) bool {
// 	request := gorequest.New()
// 	sendString := `{"device_pack_uuid": "%s","delta": %d,"host_name": ""}`
// 	response, _, _ := request.Post(configs.FetchFieldValue("PANDORA_HOST") + "/v1/api/device_packs").
// 		Send(fmt.Sprintf(sendString, devicePackUUID, 24)).
// 		End()

// 	if response.StatusCode == http.StatusOK {
// 		// s.DevicePackUUID = devicePackUUID
// 		return true
// 	}
// 	return false
// }

// func (s *SmUvGroup) UnsetDevicePackPandora() bool {
// 	request := gorequest.New()
// 	response, _, _ := request.Delete(configs.FetchFieldValue("PANDORA_HOST") + "/v1/api/device_packs/" + s.DevicePackUUID).
// 		End()

// 	if response.StatusCode == http.StatusOK {
// 		return true
// 	}
// 	return false
// }
