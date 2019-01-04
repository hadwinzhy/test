package models

import (
	"time"
)

// BaseModel ...
type BaseModel struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `gorm:"type:timestamp with time zone" json:"created_at"`
	UpdatedAt time.Time  `gorm:"type:timestamp with time zone" json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

// TimeValue ...
type TimeValue struct {
	Time  time.Time `json:"time"`
	Value uint      `json:"value"`
}

// NameValue ...
type NameValue struct {
	Name  string `json:"name"`
	Value uint   `json:"value"`
	ID    uint   `json:"id"`
}

// BaseSerializer ...
type BaseSerializer struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func (base *BaseModel) Serialize() BaseSerializer {
	return BaseSerializer{
		ID:        base.ID,
		CreatedAt: base.CreatedAt,
		UpdatedAt: base.UpdatedAt,
		DeletedAt: base.DeletedAt,
	}
}

// ImageSerializer ...
type ImageSerializer struct {
	URL string `json:"url"`
}

// AgeGroup ...
type AgeGroup struct {
	Age   string `json:"age"`
	Count uint   `json:"count"`
}

// GenderGroup ...
type GenderGroup struct {
	Female  uint `json:"female"`
	Male    uint `json:"male"`
	Unknown uint `json:"unknown"`
	NoFace  uint `json:"no_face"`
}

// GetModels will return all models
func GetModels() []interface{} {
	return []interface{}{
		&FrequentCustomerCount{},
		&FrequentCustomerPeopleBitMap{},
		&FrequentCustomerGroup{},
		&FrequentCustomerReport{},
		&FrequentCustomerRule{},
	}
}
