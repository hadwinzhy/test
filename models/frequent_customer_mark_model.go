package models

// FrequentCustomerMark 存储一些回头客的标记信息，并且将信息返回给customer model
type FrequentCustomerMark struct {
	ID       uint   `gorm:"primary_key" json:"id"`
	PersonID string `gorm:"index" json:"person_id"`
	Note     string `gorm:"type:varchar(30)" json:"note"`
	Name     string `gorm:"type:varchar(30)" json:"name"`
}
