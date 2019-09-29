package models

import "time"

type SmShop struct {
	BaseModel
	Name              string    `gorm:"type:varchar(30)" json:"name"`
	CompanyID         uint      `gorm:"index" json:"company_id"`
	BusinessTypeID    uint      `gorm:"index" json:"business_type_id"`
	DeviceCount       uint      `gorm:"type:integer" json:"device_count"`
	Location          string    `gorm:"type:varchar(128)" json:"location"`
	SmFloors          []SmFloor `gorm:"many2many:sm_shops_floors;" json:"sm_floors"`
	EntranceUvGroupID uint      `gorm:"index" json:"entrance_uv_group_id"`
	EntranceUvGroup   SmUvGroup `json:"entrance_uv_group"`
}

type SmShopSerializer struct {
	ID                uint       `json:"id"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at"`
	Name              string     `json:"name"`
	DeviceCount       uint       `json:"device_count"`
	Location          string     `json:"location"`
	BusinessTypeID    uint       `json:"business_type_id"`
	BusinessType      string     `json:"business_type"`
	SmFloors          []SmFloor  `json:"floors"`
	EntranceUvGroup   SmUvGroup  `json:"entrance_uv_group"`
	EntranceUvGroupID uint       `json:"entrance_uv_group_id"`
}

func (s SmShop) BasicSerializer(name string) SmShopSerializer {
	return SmShopSerializer{
		ID:                s.ID,
		CreatedAt:         s.CreatedAt,
		UpdatedAt:         s.UpdatedAt,
		DeletedAt:         s.DeletedAt,
		Name:              s.Name,
		DeviceCount:       s.DeviceCount,
		Location:          s.Location,
		BusinessTypeID:    s.BusinessTypeID,
		BusinessType:      name,
		SmFloors:          s.SmFloors,
		EntranceUvGroupID: s.EntranceUvGroupID,
		EntranceUvGroup:   s.EntranceUvGroup,
	}
}
