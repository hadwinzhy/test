package models

import (
	"strconv"
	"strings"

	"siren/venus/venus-controller/controllers/errors"

	"github.com/jinzhu/gorm"
)

type SmRegion struct {
	BaseModel
	CompanyID   uint        `gorm:"index" json:"company_id"`
	Name        string      `gorm:"type:varchar(30);" json:"name"`
	FloorID     uint        `gorm:"index" json:"floor_id"`
	DeviceCount uint        `gorm:"type:integer" json:"device_count"`
	Location    string      `gorm:"type:varchar(50);" json:"location"`
	SmUvGroups  []SmUvGroup `gorm:"many2many:sm_regions_uvgroups;"`
	RegionTypes string      `gorm:"type:varchar(20);" json:"region_types"`
	Floor       SmFloor
}

type SmRegionBasicSerializer struct {
	BaseSerializer
	Name              string                      `json:"name"`
	RegionTypesName   []string                    `json:"region_types_name"`
	RegionTypesDetail []SmUVGroupSimpleSerializer `json:"region_types_detail"`
	DeviceCount       uint                        `json:"device_count"`
	Location          string                      `json:"location"`
	FloorID           uint                        `json:"floor_id"`
	// Floor
}

func (region *SmRegion) BasicSerialize() SmRegionBasicSerializer {
	typesDetail := make([]SmUVGroupSimpleSerializer, len(region.SmUvGroups))
	typesName := make([]string, len(region.SmUvGroups))
	for i, group := range region.SmUvGroups {
		typesDetail[i] = group.SimpleSerialize()
		typesName[i] = group.TypeName()
	}

	if len(typesName) == 0 { // 一个降级策略，没有load uvgroup时，使用regionType

		for i := 1; i <= 5; i++ {
			s := strconv.Itoa(i)
			if strings.Contains(region.RegionTypes, s) {
				typesName = append(typesName, UVGROUP_TYPENAMEMAP[i])
			}
		}
	}

	return SmRegionBasicSerializer{
		BaseSerializer:    region.BaseModel.Serialize(),
		Name:              region.Name,
		RegionTypesName:   typesName,
		RegionTypesDetail: typesDetail,
		DeviceCount:       region.DeviceCount,
		Location:          region.Location,
		FloorID:           region.FloorID,
	}
}

// TODO: 需要在楼层/商铺等依赖的组删除时，触发regiontypes缓存列的更新
// UpdateUvGroupAssociation 更新region对应的去重组
func (region *SmRegion) UpdateUvGroupAssociation(tx *gorm.DB, uvGroups []SmUvGroup) *errors.Error {
	if err := tx.Model(region).Association("SmUvGroups").Replace(uvGroups).Error; err != nil {
		structedErr := errors.MakeDBError(err.Error())
		return &structedErr
	}
	uvGroupIDs := make([]string, len(uvGroups))
	regionTypesMap := make(map[uint]bool)

	for i, group := range uvGroups {
		uvGroupIDs[i] = strconv.Itoa(int(group.ID))
		regionTypesMap[group.RegionType] = true
	}

	regionTypesCacheStr := ""
	for k := range regionTypesMap {
		regionTypesCacheStr += strconv.Itoa(int(k))
	}

	region.RegionTypes = regionTypesCacheStr

	if err := tx.Save(&region).Error; err != nil {
		structedErr := errors.MakeDBError(err.Error())
		return &structedErr
	}

	return nil
}

// Delete 处理了region区域的删除操作
func (region *SmRegion) Delete(tx *gorm.DB) *errors.Error {
	// 删除区域中的设备
	var devices []Device
	tx.Where("company_id = ?", region.CompanyID).Where("sm_region_id = ?", region.ID).Find(&devices)

	for _, device := range devices {
		if err := device.Delete(tx); err != nil {
			structedErr := errors.MakeDBError(err.Error())
			return &structedErr
		}
	}

	// 删除关联
	if err := tx.Model(&region).Association("SmUvGroups").Clear().Error; err != nil {

		structedErr := errors.MakeDBError(err.Error())
		return &structedErr
	}

	// 楼层 region_count - 1, 删除 楼层 关联
	if region.Floor.ID != 0 {
		if region.Floor.RegionCount > 0 {
			if err := tx.Model(&region.Floor).Update("region_count", gorm.Expr("region_count - ?", 1)).Error; err != nil {
				structedErr := errors.MakeDBError(err.Error())
				return &structedErr
			}
		}

		if err := tx.Model(&region).Association("Floor").Clear().Error; err != nil {

			structedErr := errors.MakeDBError(err.Error())
			return &structedErr
		}
	}

	if err := tx.Delete(&region).Error; err != nil {

		structedErr := errors.MakeDBError(err.Error())
		return &structedErr
	}

	return nil
}

var RegionTypePriorityMap = map[uint]uint{
	1: 1,
	2: 2,
	3: 4, // 公共区域优先级比商铺低
	4: 3,
	5: 5,
}

func (region *SmRegionBasicSerializer) LeastType() uint {
	leastNum := uint(255)

	for _, detail := range region.RegionTypesDetail {
		regionType := RegionTypePriorityMap[detail.RegionType]
		if regionType < leastNum {
			leastNum = regionType
		}
	}
	return leastNum
}
