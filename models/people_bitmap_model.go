package models

import (
	"errors"
	"siren/pkg/database"
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
	PersonID                 uint   `gorm:"index"`
	BitMap                   string `gorm:"type:BIT(32)"`
	FrequentCustomerPeople   FrequentCustomerPeople
}

type FrequentCustomerPeople struct {
	BaseModel
	FrequentCustomerGroupID uint      `gorm:"index"`
	PersonID                string    `gorm:"type:varchar(32)"`
	Date                    time.Time `gorm:"type:date"`
	Interval                uint      `gorm:"type:integer"`
	Frequency               uint      `gorm:"type:integer"`

	customerType string // 隐藏字段，类型
}

// UpdateBitMap 更新对应的person的当天的bitmap并且返回出来
func (person *FrequentCustomerPeople) UpdateBitMap(today time.Time) (FrequentCustomerPeopleBitMap, error) {
	var bitMap FrequentCustomerPeopleBitMap
	database.POSTGRES.Preload("FrequentCustomerPeople").
		Where("person_id = ?").
		Order("id desc").
		First(&bitMap)

	if bitMap.ID == 0 || len(bitMap.BitMap) != 32 { // 以前从来没来过
		bitMap.FrequentCustomerPeopleID = person.ID
		bitMap.BitMap = "00000000000000000000000000000001"
		err := database.POSTGRES.Save(&bitMap).Error

		return bitMap, err
	} else {
		// 来过的话就重新计算一下bitMap保存下里
		var newBitMap FrequentCustomerPeopleBitMap
		newBitMap.FrequentCustomerPeopleID = person.ID
		days := (today.Add(time.Second).Sub(bitMap.FrequentCustomerPeople.Date)) / (86400 * time.Second) // +1s保证除尽

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
func (person *FrequentCustomerPeople) UpdateValueWithBitMap(bitMap *FrequentCustomerPeopleBitMap) {
	person.Frequency = uint(strings.Count(bitMap.BitMap, "1") - 1)
	lastIndex := strings.LastIndex(bitMap.BitMap[:len(bitMap.BitMap)-1], "1")
	if lastIndex != -1 {
		person.Interval = uint(31 - lastIndex)
	}
}

func (person *FrequentCustomerPeople) GetType() string {
	if person.customerType != "" {
		return person.customerType
	}

	if person.Frequency == 0 {
		person.customerType = FREQUENT_CUSTOMER_TYPE_NEW
	} else {
		// TODO: 根据公司获取高频规则，现在是默认规则
		var rule FrequentCustomerRule
		limit := rule.ReadableRule().Limit

		if person.Frequency >= limit {
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
