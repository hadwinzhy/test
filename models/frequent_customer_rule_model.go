package models

type FrequentCustomerRule struct {
	BaseModel
	CompanyID       uint   `gorm:"index"`
	RuleType        string `gorm:"type:varchar(16);"`
	From            int    `gorm:"type:integer" json:"from"`
	To              int    `gorm:"type:integer" json:"to"`
	isHighFrequency bool   `gorm:"type:boolean"`
}
