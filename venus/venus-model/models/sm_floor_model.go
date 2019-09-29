package models

import "time"

type SmFloor struct {
	BaseModel
	Priority          uint      `gorm:"type:integer;" json:"priority"`
	Name              string    `gorm:"type:varchar(20)" json:"name"`
	CompanyID         uint      `gorm:"index" json:"company_id"`
	ShopCount         uint      `gorm:"type:integer;" json:"shop_count"`
	RegionCount       uint      `gorm:"type:integer;" json:"region_count"`
	EntranceUvGroupID uint      `gorm:"index" json:"entrance_uv_group_id"`
	EntranceUvGroup   SmUvGroup `json:"entrance_uv_group"`
	PublicUvGroupID   uint      `gorm:"index" json:"public_uv_group_id"`
	PublicUvGroup     SmUvGroup `json:"public_uv_group"`
	Shops             []SmShop  `gorm:"many2many:sm_shops_floors" json:"shops"`
}

type SmFloorSerializer struct {
	ID                uint       `json:"id"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at"`
	CompanyID         uint       `json:"company_id"`
	Priority          uint       `json:"priority"`
	Name              string     `json:"name"`
	ShopCount         uint       `json:"shop_count"`
	RegionCount       uint       `json:"region_count"`
	EntranceUvGroupID uint       `json:"entrance_uv_group_id"`
	EntranceUvGroup   SmUvGroup  `json:"entrance_uv_group"`
	PublicUvGroupID   uint       `json:"public_uv_group_id"`
	PublicUvGroup     SmUvGroup  `json:"public_uv_group"`
	Shops             []SmShop   `json:"shops"`
}

func (s SmFloor) BasicSerializer() SmFloorSerializer {
	return SmFloorSerializer{
		ID:                s.ID,
		CreatedAt:         s.CreatedAt,
		UpdatedAt:         s.UpdatedAt,
		DeletedAt:         s.DeletedAt,
		Priority:          s.Priority,
		CompanyID:         s.CompanyID,
		Name:              s.Name,
		ShopCount:         s.ShopCount,
		RegionCount:       s.RegionCount,
		EntranceUvGroupID: s.EntranceUvGroupID,
		EntranceUvGroup:   s.EntranceUvGroup,
		PublicUvGroupID:   s.PublicUvGroupID,
		PublicUvGroup:     s.PublicUvGroup,
		Shops:             s.Shops,
	}
}
