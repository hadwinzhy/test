package models

// XcloudAccount ...
type XcloudAccount struct {
	BaseModel
	UID          string `gorm:"type:varchar(50);index;not null" json:"uid"`
	Password     string `gorm:"type:varchar(50);not null" json:"password"`
	SerialNumber string `gorm:"type:varchar(50);" json:"serial_number"`
	IsUsed       bool   `gorm:"default:false;not null" json:"is_used"`
	DeviceID     uint
}

// XcloudAccountSerializer ...
type XcloudAccountSerializer struct {
	BaseSerializer
	UID          string `json:"uid"`
	Password     string `json:"password"`
	SerialNumber string `json:"serial_number"`
	DeviceID     uint   `json:"device_id"`
	IsUsed       bool   `json:"is_used"`
}

// BasicSerialize ...
func (x *XcloudAccount) BasicSerialize() XcloudAccountSerializer {
	return XcloudAccountSerializer{
		BaseSerializer: BaseSerializer{
			ID:        x.ID,
			CreatedAt: x.CreatedAt,
			UpdatedAt: x.UpdatedAt,
		},
		UID:          x.UID,
		Password:     x.Password,
		SerialNumber: x.SerialNumber,
		DeviceID:     x.DeviceID,
		IsUsed:       x.IsUsed,
	}
}

// RefSerialize ...
func (x *XcloudAccount) RefSerialize() (ref *XcloudAccountSerializer) {
	value := x.BasicSerialize()
	if value.ID > 0 {
		ref = &value
	}
	return ref
}
