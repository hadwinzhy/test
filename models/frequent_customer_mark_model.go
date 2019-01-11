package models

// FrequentCustomerMark 存储一些回头客的标记信息，并且将信息返回给customer model
type FrequentCustomerMark struct {
	ID       uint   `gorm:"primary_key"`
	PersonID string `gorm:"index"`
	Note     string `gorm:"type:varchar(30)"`
	Name     string `gorm:"type:varchar(30)"`
}
