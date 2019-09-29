package models

type CompanyGroup struct {
	BaseModel
	Name string `gorm:"type:varchar(30)" json:"name"`
}
