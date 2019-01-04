package workers

import (
	"siren/models"
	"siren/pkg/database"
	"siren/pkg/utils"
	"time"
)

type ComerBasicInfo struct {
	Interval uint
	Times    uint
}

func (comer *ComerBasicInfo) IsNewComer() bool {
	return comer.Interval == 0 && comer.Times == 0
}

func (comer *ComerBasicInfo) IsHighFrequentComer(ruleLimit uint) bool {
	return comer.Times > ruleLimit
}

func fetchFrequentCustomerGroup(companyID uint, shopID uint) (models.FrequentCustomerGroup, error) {
	var group models.FrequentCustomerGroup

	err := database.POSTGRES.FirstOrCreate(&group, models.FrequentCustomerGroup{
		CompanyID: companyID,
		ShopID:    shopID,
	}).Error

	return group, err
}

func StoreFrequentCustomerHandler(companyID uint, shopID uint, personID string, captureAt int64) {
	// 来了个新客

	// 0. 看看有没有frequent customer group
	fcGroup, err := fetchFrequentCustomerGroup(companyID, shopID)
	if err != nil {
		return
	}

	// 1. 首先看这组companyID shopID里有没有这个personID的bitmap，bitmap里记录了一个值，当天这人有没有来过
	var person models.FrequentCustomerPeople
	// var comer comerBasicInfo

	today := utils.CurrentDate(time.Now())

	database.POSTGRES.Where("frequent_customer_group_id = ?", fcGroup.ID).
		Where("person_id = ?", personID).
		Where("date = ?", today).
		First(&person)

	if person.ID == 0 { // 新建一个person
		// bitMap.BitMap = "00000000000000000000000000000001"
		person.FrequentCustomerGroupID = fcGroup.ID
		person.PersonID = personID
		person.Date = today
		err := database.POSTGRES.Save(&person).Error
		if err != nil {
			// return err
		}
	} else { // 已经有bitMap时的操作

		// bitMap.BitMap = string(append([]byte(bitMap.BitMap)[0:31], '1'))
		// firstDistance := 0
		// for i := 30; i >= 0; i-- {
		// 	if firstDistance == 0 && bitMap.BitMap[i] == '1' {
		// 		firstDistance = 31 - i
		// 	}
		// }
	}

	// 1.1 有的话，这就是一个来过的人，记在bitmap中更新那一行

	// 1.2 没有的话，就没来过，给bitmap添加一行

	// 2. 根据bitmap中的频次，记到当天的数据分布表中, 总人数，高频次数，低频次数，新客数，总到访间隔天数，总到访天数

	// 3. 取一下频率规则，判断是不是高频次的人

	// 3.1 是的话，根据来的captureAt时间，记到高频表里

}
