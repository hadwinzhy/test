package workers

import (
	"errors"
	"fmt"
	"net/http"
	"siren/configs"
	"siren/models"
	"siren/pkg/database"
	"siren/pkg/logger"
	"siren/pkg/utils"
	"strconv"
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

func fetchFrequentCustomerPerson(group *models.FrequentCustomerGroup, personID string, today time.Time, hour time.Time, eventID uint) (models.FrequentCustomerPeople, error) {
	var person models.FrequentCustomerPeople
	// var comer comerBasicInfo

	database.POSTGRES.Preload("FrequentCustomerGroup").
		Where("frequent_customer_group_id = ?", group.ID).
		Where("person_id = ?", personID).
		Where("hour >= ?", today).
		Where("hour < ?", today.AddDate(0, 0, 1)).
		First(&person)

	if person.ID == 0 { // 新建一个person
		person.FrequentCustomerGroupID = group.ID
		person.PersonID = personID
		person.Hour = hour
		person.EventID = eventID
		person.FrequentCustomerGroup = *group
		err := database.POSTGRES.Save(&person).Error
		if err != nil {
			return person, err
		}
	} else {
		logger.Error("store_worker", "fetch_person", "already_exists personID = ", personID, " groupID = ", group.ID, " captureTime = ", hour)
		return person, errors.New("今天已经来过，不能用作回头客")
	}
	return person, nil
}

func updateBitMap(frequentPerson *models.FrequentCustomerPeople, today time.Time) (models.FrequentCustomerPeopleBitMap, error) {
	frequentPersonID := frequentPerson.ID
	personID := frequentPerson.PersonID
	var bitMap models.FrequentCustomerPeopleBitMap
	database.POSTGRES.Preload("FrequentCustomerPeople").
		Where("person_id = ?", personID).
		Order("id desc").
		First(&bitMap)

	if bitMap.ID == 0 || len(bitMap.BitMap) != 32 { // 以前从来没来过
		bitMap.FrequentCustomerPeopleID = frequentPersonID
		bitMap.PersonID = personID
		bitMap.BitMap = "00000000000000000000000000000001"
		err := database.POSTGRES.Save(&bitMap).Error

		return bitMap, err
	} else {
		// 来过的话就重新计算一下bitMap保存下里
		frequentPerson.DefaultNumber = bitMap.FrequentCustomerPeople.DefaultNumber
		frequentPerson.LastCaptureAt = bitMap.FrequentCustomerPeople.Hour

		var newBitMap models.FrequentCustomerPeopleBitMap
		newBitMap.FrequentCustomerPeopleID = frequentPersonID
		newBitMap.PersonID = personID
		lastBitMapDate := utils.CurrentDate(bitMap.FrequentCustomerPeople.Hour) // 存储的时候是0时区
		fmt.Println(today, lastBitMapDate, today.Add(time.Second).Sub(lastBitMapDate))
		days := (today.Add(time.Second).Sub(lastBitMapDate)) / (86400 * time.Second) // +1s保证除尽

		if days > 30 || days <= 0 {
			newBitMap.BitMap = "00000000000000000000000000000001"
		} else {
			bitMapNum, err := strconv.ParseInt(bitMap.BitMap, 2, 64)
			fmt.Println("bitmap", bitMapNum)
			if err != nil {
				return newBitMap, err
			}

			bitMapNum = bitMapNum << uint(days)
			fmt.Println("bitmap", bitMapNum)
			bitMapNum += 4294967296
			newBitMapStr := strconv.FormatInt(bitMapNum, 2)

			fmt.Println("bitmap", newBitMapStr)

			if len(newBitMapStr) < 32 { // 没到32，往前补0
				return newBitMap, errors.New("bit error")
			}
			newBitMapStr = "00" + newBitMapStr[len(newBitMapStr)-30:len(newBitMapStr)-1] + "1" // 用不着32位，用30位，最后1位为1

			newBitMap.BitMap = newBitMapStr
		}

		err := database.POSTGRES.Save(&newBitMap).Error
		return newBitMap, err
	}
}

func updateFrequentCustomerReport(person *models.FrequentCustomerPeople, groupID uint, today time.Time, hour time.Time) error {
	var report models.FrequentCustomerReport
	database.POSTGRES.FirstOrCreate(
		&report,
		models.FrequentCustomerReport{
			FrequentCustomerGroupID: groupID,
			Date:                    today,
			Hour:                    hour,
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
			Date:                    today,
		},
	)

	table.AddCount(time.Unix(captureAt, 0))
}

func StoreFrequentCustomerHandler(companyID uint, shopID uint, personID string, captureAt int64, eventID uint) {
	// 来了个新客

	// 0. 看看有没有frequent customer group
	fcGroup, err := fetchFrequentCustomerGroup(companyID, shopID)
	if err != nil {
		return
	}

	captureTime := time.Unix(captureAt, 0)

	today := utils.CurrentDate(captureTime)
	thisHour := utils.CurrentTime(captureTime, "hour")
	// 1. 首先看这组companyID shopID里有没有这个personID的bitmap，bitmap里记录了一个值，当天这人有没有来过
	person, err := fetchFrequentCustomerPerson(&fcGroup, personID, today, captureTime, eventID)

	if err != nil {
		return
	}

	// 1.1 有的话，这就是一个来过的人，记在bitmap中更新那一行
	bitMap, err := updateBitMap(&person, today)
	if err != nil {
		return
	}

	// 1.2 person里的数据更新
	person.UpdateValueWithBitMap(&bitMap, &fcGroup)

	// 2. report里的数据更新，记到当天的数据分布表中, 总人数，高频次数，低频次数，新客数，总到访间隔天数，总到访天数
	err = updateFrequentCustomerReport(&person, fcGroup.ID, today, thisHour)
	if err != nil {
		return
	}

	if person.GetType() != models.FREQUENT_CUSTOMER_TYPE_NEW {
		go markNameNote(person.EventID, person.PersonID)
	}

	// 是回头客，触发 venus 回头客 消息推送
	if person.GetType() != models.FREQUENT_CUSTOMER_TYPE_NEW {
		go func(eventID uint) {
			if eventID != 0 {
				url := fmt.Sprintf(configs.FetchFieldValue("VENUSHOST")+"/v1/api/company/notification_frequent_person?event_id=%d&last_captured_at=%d&frequent_customer_id=%d", person.EventID, person.LastCaptureAt.Unix(), person.ID)
				request, _ := http.NewRequest(http.MethodGet, url, nil)
				client := http.DefaultClient
				response, err := client.Do(request)
				if err != nil {
					return
				}
				defer response.Body.Close()
			}
		}(person.EventID)
	}

	// 3. 取一下频率规则，判断是不是高频次的人
	if person.IsHighFrequency() {
		// 3.1 是的话，根据来的captureAt时间，记到高频表里
		updateFrequentCustomerHighTimeTable(fcGroup.ID, today, captureAt)
	}
}

func RemoveFrequentCustomerHandler(personID string) {
	var people []models.FrequentCustomerPeople

	if personID != "" {
		database.POSTGRES.Where("person_id = ?", personID).Delete(&people) // 查看回头客列表里就没有了
	}
}

func markNameNote(eventID uint, personID string) {
	if personID != "" && eventID != 0 {
		var event models.Event
		database.POSTGRES.First(&event, eventID)
		if event.ID == 0 || event.CustomerID == 0 {
			return
		}

		// eventid 和 customerid都不是0
		var mark models.FrequentCustomerMark
		database.POSTGRES.Where("person_id = ?", personID).First(&mark)
		if mark.ID > 0 {
			database.POSTGRES.Table("customers").
				Where("id = ?", event.CustomerID).
				Updates(map[string]interface{}{"name": mark.Name, "note": mark.Note})
		}
	}
}
