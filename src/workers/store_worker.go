package workers

import (
	"errors"
	"siren/models"
	"siren/pkg/database"
	"siren/pkg/utils"
	"time"
)

func fetchFrequentCustomerGroup(companyID uint, shopID uint) (models.FrequentCustomerGroup, error) {
	var group models.FrequentCustomerGroup

	err := database.POSTGRES.FirstOrCreate(&group, models.FrequentCustomerGroup{
		CompanyID: companyID,
		ShopID:    shopID,
	}).Error

	return group, err
}

func fetchFrequentCustomerPerson(group *models.FrequentCustomerGroup, personID string, today time.Time) (models.FrequentCustomerPeople, error) {
	var person models.FrequentCustomerPeople
	// var comer comerBasicInfo

	database.POSTGRES.Where("frequent_customer_group_id = ?", group.ID).
		Where("person_id = ?", personID).
		Where("date = ?", today).
		First(&person)

	if person.ID == 0 { // 新建一个person
		person.FrequentCustomerGroupID = group.ID
		person.PersonID = personID
		person.Date = today
		err := database.POSTGRES.Save(&person).Error
		if err != nil {
			return person, err
		}
	}
	return person, nil
}

func updateFrequentCustomerReport(person *models.FrequentCustomerPeople, groupID uint, today time.Time) error {
	var report models.FrequentCustomerReport
	database.POSTGRES.FirstOrCreate(
		&report,
		models.FrequentCustomerReport{
			FrequentCustomerGroupID: groupID,
			Date: today,
		},
	)

	switch person.GetType() {
	case models.FREQUENT_CUSTOMER_TYPE_HIGH:
		report.HighFrequency++
	case models.FREQUENT_CUSTOMER_TYPE_LOW:
		report.LowFrequency++
	case models.FREQUENT_CUSTOMER_TYPE_NEW:
		report.NewComer++
	default:
		return errors.New("回头客的类型不存在")
	}

	report.SumInterval += person.Interval
	report.SumTimes += person.Frequency

	err := database.POSTGRES.Save(&report).Error
	return err
}

func updateFrequentCustomerHighTimeTable(groupID uint, today time.Time, captureAt int64) {
	var table models.FrequentCustomerHighTimeTable
	database.POSTGRES.FirstOrCreate(
		&table,
		models.FrequentCustomerHighTimeTable{
			FrequentCustomerGroupID: groupID,
			Date: today,
		},
	)

	table.AddCount(time.Unix(captureAt, 0))
}

func StoreFrequentCustomerHandler(companyID uint, shopID uint, personID string, captureAt int64) {
	// 来了个新客

	// 0. 看看有没有frequent customer group
	fcGroup, err := fetchFrequentCustomerGroup(companyID, shopID)
	if err != nil {
		return
	}

	today := utils.CurrentDate(time.Now())
	// 1. 首先看这组companyID shopID里有没有这个personID的bitmap，bitmap里记录了一个值，当天这人有没有来过
	person, err := fetchFrequentCustomerPerson(&fcGroup, personID, today)
	if err != nil {
		return
	}

	// 1.1 有的话，这就是一个来过的人，记在bitmap中更新那一行
	bitMap, err := person.UpdateBitMap(today)
	if err != nil {
		return
	}

	// 1.2 person里的数据更新
	person.UpdateValueWithBitMap(&bitMap)

	// 2. report里的数据更新，记到当天的数据分布表中, 总人数，高频次数，低频次数，新客数，总到访间隔天数，总到访天数
	err = updateFrequentCustomerReport(&person, fcGroup.ID, today)
	if err != nil {
		return
	}

	// 3. 取一下频率规则，判断是不是高频次的人
	if person.IsHighFrequency() {
		// 3.1 是的话，根据来的captureAt时间，记到高频表里
		updateFrequentCustomerHighTimeTable(fcGroup.ID, today, captureAt)
	}
}
