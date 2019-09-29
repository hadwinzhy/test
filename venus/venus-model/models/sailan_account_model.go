package models

// SailanAccount ...
type SailanAccount struct {
	BaseModel
	CID        string `gorm:"type:varchar(50);index;not null" json:"cid"`
	Username   string `gorm:"type:varchar(50);not null" json:"username"`
	Password   string `gorm:"type:varchar(50);" json:"password"`
	IsUsed     bool   `gorm:"default:false;not null" json:"is_used"`
	MacAddress string `gorm:"type:varchar(50);" json:"mac_address"`
	DeviceID   uint
}

// SailanAccountSerializer ...
type SailanAccountSerializer struct {
	BaseSerializer
	CID        string `json:"cid"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	DeviceID   uint   `json:"device_id"`
	MacAddress string `json:"mac_address"`
	IsUsed     bool   `json:"is_used"`
}

// BasicSerialize ...
func (x *SailanAccount) BasicSerialize() SailanAccountSerializer {
	return SailanAccountSerializer{
		BaseSerializer: BaseSerializer{
			ID:        x.ID,
			CreatedAt: x.CreatedAt,
			UpdatedAt: x.UpdatedAt,
		},
		CID:        x.CID,
		Username:   x.Username,
		Password:   x.Password,
		DeviceID:   x.DeviceID,
		IsUsed:     x.IsUsed,
		MacAddress: x.MacAddress,
	}
}

// RefSerialize ...
func (x *SailanAccount) RefSerialize() (ref *SailanAccountSerializer) {
	value := x.BasicSerialize()
	if value.ID > 0 {
		ref = &value
	}
	return ref
}
