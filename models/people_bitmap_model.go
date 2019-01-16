package models

import (
	"errors"
	"siren/pkg/database"
	"siren/pkg/utils"
	"strconv"
	"strings"
	"time"
)

const (
	FREQUENT_CUSTOMER_TYPE_HIGH = "high"
	FREQUENT_CUSTOMER_TYPE_LOW  = "low"
	FREQUENT_CUSTOMER_TYPE_NEW  = "new"
)

type FrequentCustomerPeopleBitMap struct {
	BaseModel
	FrequentCustomerPeopleID uint   `gorm:"index"`
	PersonID                 string `gorm:"index"`
	BitMap                   string `gorm:"type:BIT(32)"`
	FrequentCustomerPeople   FrequentCustomerPeople
}

type FrequentCustomerPeople struct {
	BaseModel
	FrequentCustomerGroupID uint      `gorm:"index"`
	PersonID                string    `gorm:"type:varchar(32);index"`
	Date                    time.Time `gorm:"type:date"`
	Hour                    time.Time `gorm:"type:timestamp with time zone"`
	Interval                uint      `gorm:"type:integer"`
	Frequency               uint      `gorm:"type:integer"`
	IsFrequentCustomer      bool      `gorm:"type:bool"`
	EventID                 uint      `gorm:"type:integer"`
	DefaultNumber           uint      `gorm:"type:integer;"`
	Event                   Event
	FrequentCustomerGroup   FrequentCustomerGroup

	customerType string // 隐藏字段，类型
}

// UpdateBitMap 更新对应的person的当天的bitmap并且返回出来
// personID是算法层面的personID
func (person *FrequentCustomerPeople) UpdateBitMap(personID string, today time.Time) (FrequentCustomerPeopleBitMap, error) {
	var bitMap FrequentCustomerPeopleBitMap
	database.POSTGRES.Preload("FrequentCustomerPeople").
		Where("person_id = ?", personID).
		Order("id desc").
		First(&bitMap)

	if bitMap.ID == 0 || len(bitMap.BitMap) != 32 { // 以前从来没来过
		bitMap.FrequentCustomerPeopleID = person.ID
		bitMap.PersonID = personID
		bitMap.BitMap = "00000000000000000000000000000001"
		err := database.POSTGRES.Save(&bitMap).Error

		return bitMap, err
	} else {
		// 来过的话就重新计算一下bitMap保存下里
		var newBitMap FrequentCustomerPeopleBitMap
		newBitMap.FrequentCustomerPeopleID = person.ID
		newBitMap.PersonID = personID
		lastDate := utils.CurrentDate(bitMap.FrequentCustomerPeople.Hour)
		days := (today.Add(time.Second).Sub(lastDate)) / (86400 * time.Second) // +1s保证除尽

		if days > 30 || days <= 0 {
			newBitMap.BitMap = "00000000000000000000000000000001"
		} else {
			bitMapNum, err := strconv.ParseInt(bitMap.BitMap, 2, 64)
			if err != nil {
				return newBitMap, err
			}

			bitMapNum = bitMapNum << uint(days)
			bitMapNum += 4611686018427387904 // 为了保证所有的数据字符串都是大于32位的，加了2^62

			newBitMapStr := strconv.FormatInt(bitMapNum, 2)

			if len(newBitMapStr) < 32 { // 没到32，往前补0
				return newBitMap, errors.New("bit error")
			} else {
				newBitMapStr = "00" + newBitMapStr[len(newBitMapStr)-30:len(newBitMapStr)-1] + "1" // 用不着32位，用30位，最后1位为1
			}

			newBitMap.BitMap = newBitMapStr
		}

		err := database.POSTGRES.Save(&newBitMap).Error
		return newBitMap, err
	}
}

// UpdateValueWithBitMap 根据bitMap，person的数据得到更新
func (person *FrequentCustomerPeople) UpdateValueWithBitMap(bitMap *FrequentCustomerPeopleBitMap, group *FrequentCustomerGroup) {
	person.Frequency = uint(strings.Count(bitMap.BitMap, "1"))
	lastIndex := strings.LastIndex(bitMap.BitMap[:len(bitMap.BitMap)-1], "1")
	if lastIndex != -1 {
		person.Interval = uint(31 - lastIndex)
		person.IsFrequentCustomer = true
	} else {
		person.IsFrequentCustomer = false
	}

	if person.IsFrequentCustomer {
		group.DefaultNumber++
		person.DefaultNumber = group.DefaultNumber
		database.POSTGRES.Save(group)
	}

	database.POSTGRES.Save(person)
}

func (person *FrequentCustomerPeople) GetType() string {
	if person.customerType != "" {
		return person.customerType
	}

	if person.Frequency <= 1 { // 一次以内就不算回头客
		person.customerType = FREQUENT_CUSTOMER_TYPE_NEW
	} else {
		var rule FrequentCustomerRule

		database.POSTGRES.Where("company_id = ?", person.FrequentCustomerGroup.CompanyID).First(&rule)

		limit := rule.ReadableRule().Limit

		if person.Frequency > limit {
			person.customerType = FREQUENT_CUSTOMER_TYPE_HIGH
		} else {
			person.customerType = FREQUENT_CUSTOMER_TYPE_LOW
		}
	}

	return person.customerType
}

func (person *FrequentCustomerPeople) IsHighFrequency() bool {
	return person.GetType() == FREQUENT_CUSTOMER_TYPE_HIGH
}
