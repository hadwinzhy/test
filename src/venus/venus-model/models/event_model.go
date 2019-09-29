package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

var (
	Male   uint
	Female uint
	In     string
	Out    string
)

func init() {
	Male = 1
	Female = 0
	In = "in"
	Out = "out"
}

// Event ...
type Event struct {
	BaseModel
	CustomerID   uint      `gorm:"index" json:"customer_id"`
	ShopID       uint      `gorm:"index;not null" json:"shop_id"`
	SmRegionID   uint      `gorm:"index" json:"sm_region_id"`
	EventType    string    `gorm:"type:varchar(100);not null" json:"event_type"`
	CaptureAt    time.Time `gorm:"type:timestamp with time zone" json:"capture_at"`
	DeviceID     uint      `gorm:"index;not null" json:"device_id"`
	DeviceName   string    `gorm:"type:varchar(50)" json:"device_name"`
	CompanyID    uint      `gorm:"index;not null" json:"company_id"`
	FaceID       string    `gorm:"index" json:"face_id"`
	Age          uint      `gorm:"default:0" json:"age"`
	Gender       uint      `json:"gender"` // 0 女， 1男
	OriginalFace string    `gorm:"type:varchar(255)" json:"original_face"`
	PersonID     string    `gorm:"index;type:varchar(255)" json:"person_id"`
	Status       string    `gorm:"type:varchar(100)" json:"status"`
	TrackID      string    `gorm:"index;type:varchar(100)" json:"track_id"`
	EventID      string    `gorm:"index;type:varchar(100)" json:"event_id"`
	FacesCount   int       `gorm:"default:0;not null" json:"faces_count"`
	HasSent      bool      `gorm:"default:false;not null" json:"has_sent"`
	BestFaceID   uint      `gorm:"index" json:"best_face_id"`
	Confidence   string    `gorm:"type:varchar(255)" json:"confidence"`
	FacePitch    string    `gorm:"type:varchar" json:"face_pitch"`
	FaceYaw      string    `gorm:"type:varchar" json:"face_yaw"`
	FaceRoll     string    `gorm:"type:varchar" json:"face_roll"`
	Shop         Shop
	SmRegion     SmRegion
	Device       Device
	Customer     Customer
}

type EventSerializer struct {
	ID                  uint            `json:"id"`
	PersonID            string          `json:"person_id"`
	DeviceID            uint            `json:"device_id"`
	CustomerID          uint            `json:"customer_id"`
	CustomerEventsCount uint            `json:"customer_events_count"`
	CustomerName        string          `json:"customer_name"`
	EventType           string          `json:"event_type"`
	OriginalFace        ImageSerializer `json:"original_face"`
	OriginalFaceURL     string          `json:"original_face_url"`
	Status              string          `json:"status"`
	CaptureAt           time.Time       `json:"capture_at"`
	CreatedAt           time.Time       `json:"created_at"`
	DeviceName          string          `json:"device_name"`
	TrackID             string          `json:"track_id"`
	Age                 uint            `json:"age"`
	Gender              uint            `json:"gender"`
	Confidence          string          `json:"confidence"`
	CustomerGroupID     uint            `json:"customer_group_id"`
	CustomerGroupName   string          `json:"customer_group_name"`
	FacePitch           string          `json:"face_pitch"`
	FaceYaw             string          `json:"face_yaw"`
	FaceRoll            string          `json:"face_roll"`
}

// EventBasicSerializer ...
type EventBasicSerializer struct {
	EventSerializer
	ShopID   uint   `json:"shop_id"`
	ShopName string `json:"shop_name"`
}

// EventMallSerializer 为shopping mall版本， 不对应商铺，对应region
type EventMallSerializer struct {
	EventSerializer
	SmRegionID   uint   `json:"region_id"`
	SmRegionName string `json:"region_name"`
}

type EventGeneralSerializer struct {
	EventSerializer
	ShopID           uint   `json:"shop_id"`
	ShopName         string `json:"shop_name"`
	SmRegionID       uint   `json:"region_id"`
	SmRegionName     string `json:"region_name"`
	DeviceMacAddress string `json:"device_mac_address"`
}

func (e *Event) Serialize() EventSerializer {
	return EventSerializer{
		ID:                  e.ID,
		PersonID:            e.PersonID,
		DeviceID:            e.DeviceID,
		CustomerID:          e.CustomerID,
		CustomerName:        e.Customer.Name,
		CustomerEventsCount: e.Customer.EventsCount,
		EventType:           e.EventType,
		Status:              e.Status,
		CaptureAt:           e.CaptureAt,
		CreatedAt:           e.CreatedAt.Round(time.Second),

		Age:             e.Age,
		Gender:          e.Gender,
		DeviceName:      e.DeviceName,
		TrackID:         e.TrackID,
		OriginalFaceURL: e.OriginalFace,
		Confidence:      e.Confidence,
		OriginalFace: ImageSerializer{
			URL: e.OriginalFace,
		},
		CustomerGroupID: e.Customer.CustomerGroupID,
		FacePitch:       e.FacePitch,
		FaceYaw:         e.FaceYaw,
		FaceRoll:        e.FaceRoll,
	}
}

// BasicSerialize ...
func (e *Event) BasicSerialize() EventBasicSerializer {
	return EventBasicSerializer{
		EventSerializer: e.Serialize(),
		ShopID:          e.ShopID,
		ShopName:        e.Shop.Name,
	}
}

// MallSerialize ...
func (e *Event) MallSerialize() EventMallSerializer {
	return EventMallSerializer{
		EventSerializer: e.Serialize(),
		SmRegionID:      e.SmRegionID,
		SmRegionName:    e.SmRegion.Name,
	}
}

// GeneralSerialize 超级管理员使用的serializer
func (e *Event) GeneralSerialize() EventGeneralSerializer {
	return EventGeneralSerializer{
		EventSerializer:  e.Serialize(),
		SmRegionID:       e.SmRegionID,
		SmRegionName:     e.SmRegion.Name,
		ShopID:           e.ShopID,
		ShopName:         e.Shop.Name,
		DeviceMacAddress: e.Device.MacAddress,
	}
}

// GetRelatedDevicePtr ...
func (e *Event) GetRelatedDevicePtr(tx *gorm.DB) *Device {
	if e.Device.ID == 0 {
		var device Device
		tx.Model(e).Related(&device)
		e.Device = device
	}

	if e.Device.ID == 0 {
		return nil
	}

	return &(e.Device)
}

// GetRelatedShopPtr ...
func (e *Event) GetRelatedShopPtr(tx *gorm.DB) *Shop {
	if e.Shop.ID == 0 {
		var shop Shop
		tx.Model(e).Related(&shop)
		e.Shop = shop
	}

	if e.Shop.ID == 0 {
		return nil
	}

	return &(e.Shop)
}
