package records

import (
	"fmt"
	"siren/models"
	"siren/pkg/controllers"
	"siren/pkg/controllers/errors"
	"siren/pkg/database"
	"strconv"

	"github.com/gin-gonic/gin"
)

func recordListMaker(peopleList []models.FrequentCustomerPeople) []FrequentCustomerRecord {
	result := make([]FrequentCustomerRecord, len(peopleList))
	var shopIDSlice []uint
	var personIDSlice []string

	for i, people := range peopleList {
		shopIDSlice = append(shopIDSlice, people.Event.ShopID)
		personIDSlice = append(personIDSlice, people.PersonID)
		fmt.Println(people.DefaultNumber)
		result[i] = FrequentCustomerRecord{
			FrequentCustomerPersonID: people.ID,
			FirstCaptureURL:          people.Event.OriginalFace,
			Name:                     fmt.Sprintf("回头客%d", people.DefaultNumber), // 要根据person_id对应的去取，作标记的时候再做
			CaptureAt:                people.Event.CaptureAt,
			LastCaptureAt:            people.LastCaptureAt,
			Age:                      people.Event.Age,
			Gender:                   people.Event.Gender,
			ShopID:                   people.Event.ShopID,
			ShopName:                 "", // 根据shopID去取
			DeviceID:                 people.Event.DeviceID,
			DeviceName:               people.Event.DeviceName,
			Frequency:                people.Frequency,
			Note:                     "", // 要根据person_id对应的去取，作标记的时候再做
			personID:                 people.PersonID,
		}
	}

	firstPicMap := make(map[string]string)
	if len(personIDSlice) > 0 {
		// manually load first pic
		var persons []models.FrequentCustomerPeople
		database.POSTGRES.Model(&models.FrequentCustomerPeople{}).
			Preload("Event").
			Select("person_id, min(event_id) AS event_id").
			Where("person_id in (?)", personIDSlice).
			Group("person_id").
			Find(&persons)

		for _, person := range persons {
			firstPicMap[person.PersonID] = person.Event.OriginalFace
		}

		// manually  load  mark
		personIDMap := make(map[string]models.FrequentCustomerMark)

		var marks []models.FrequentCustomerMark

		database.POSTGRES.Where("person_id in (?)", personIDSlice).Find(&marks)

		for i := range marks {
			personIDMap[marks[i].PersonID] = marks[i]
		}

		for i := range result {
			result[i].Note = personIDMap[result[i].personID].Note
			if personIDMap[result[i].personID].Name != "" {
				result[i].Name = personIDMap[result[i].personID].Name
			}

			result[i].FirstCaptureURL = firstPicMap[result[i].personID]
		}
	}

	// manually load map
	if len(shopIDSlice) > 0 {
		shopIDMap := make(map[uint]string)
		var shopIDName []struct {
			ID   uint   `json:"id"`
			Name string `json:"name"`
		}

		database.POSTGRES.Table("shops").Select("id, name").Where("id in (?)", shopIDSlice).Where("deleted_at is NULL").Find(&shopIDName)
		for _, pair := range shopIDName {
			shopIDMap[pair.ID] = pair.Name
		}

		for i := range result {
			result[i].ShopName = shopIDMap[result[i].ShopID]
		}
	}

	return result
}

func eventListMaker(allPeople []models.FrequentCustomerPeople) []SingleEventRecord {
	result := make([]SingleEventRecord, len(allPeople))

	var shopIDSlice []uint

	for i, people := range allPeople {
		event := people.Event
		shopIDSlice = append(shopIDSlice, event.ShopID)
		result[i] = SingleEventRecord{
			ID:              event.ID,
			OriginalFaceURL: event.OriginalFace,
			CaptureAt:       event.CaptureAt,
			DeviceName:      event.DeviceName,
			ShopID:          event.ShopID,
			DeviceID:        event.DeviceID,
		}
	}

	if len(shopIDSlice) > 0 {
		shopIDMap := make(map[uint]string)
		var shopIDName []struct {
			ID   uint   `json:"id"`
			Name string `json:"name"`
		}

		database.POSTGRES.Table("shops").Select("id, name").Where("id in (?)", shopIDSlice).Where("deleted_at is NULL").Find(&shopIDName)
		for _, pair := range shopIDName {
			shopIDMap[pair.ID] = pair.Name
		}

		for i := range result {
			result[i].ShopName = shopIDMap[result[i].ShopID]
		}
	}

	return result

}

func RecordListProcessor(form FrequentCustomerRecordParams) ([]FrequentCustomerRecord, controllers.PaginationResponse, *errors.Error) {
	fcGroups := models.FetchFrequentCustomerGroup(form.CompanyID, form.shopIDs)

	var groupIDs []uint
	for _, group := range fcGroups {
		groupIDs = append(groupIDs, group.ID)
	}

	fromTime, toTime := form.GetFromAndToTime()

	// 回头group范围，根据personid聚合取最新的
	query := database.POSTGRES.Model(&models.FrequentCustomerPeople{}).
		Preload("Event").
		Select("person_id, frequent_customer_group_id, max(id) AS id, max(event_id) AS event_id, max(frequency) AS frequency, max(hour) AS hour, max(last_capture_at) AS last_capture_at, max(default_number) AS default_number")

	query = query.Where("frequent_customer_group_id in (?)", groupIDs)

	query = query.Where("is_frequent_customer = ? ", true)

	// 时间范围
	query = query.Where("hour >= ?", fromTime).Where("hour < ?", toTime)

	query = query.Group("person_id, frequent_customer_group_id")

	order := form.OrderBy + " " + form.SortBy

	var total int
	query.Count(&total)

	var peopleList []models.FrequentCustomerPeople
	query.Order(order).Offset(form.PerPage * (form.Page - 1)).Limit(form.PerPage).Find(&peopleList)

	// 组装pagination
	paginations := controllers.PaginationResponse{
		Page:  form.Page,
		Per:   form.PerPage,
		Total: total,
	}
	// 组装结果
	result := recordListMaker(peopleList)

	return result, paginations, nil
}

func recordDetailListProcessor(
	people models.FrequentCustomerPeople,
	form FrequentCustomerRecordDetailParams,
) ([]SingleEventRecord, controllers.PaginationResponse, *errors.Error) {
	var result []SingleEventRecord
	var paginations controllers.PaginationResponse
	if people.PersonID == "" {
		structedErr := errors.MakeNotFoundError("回头客没有对应的personid")
		return result, paginations, &structedErr
	}

	toTime := people.Hour                 // 回头客最后一次抓到的时间
	fromTime := toTime.AddDate(0, 0, -30) // 30天前

	var allPeople []models.FrequentCustomerPeople
	query := database.POSTGRES.Model(&allPeople).Preload("Event").
		Where("hour >= ?", fromTime).
		Where("hour <= ?", toTime).
		Where("person_id = ?", people.PersonID)

	var total int
	query.Count(&total)

	query.Order("hour desc").
		Limit(form.PerPage).
		Offset(form.PerPage * (form.Page - 1)).
		Find(&allPeople)

	paginations = controllers.PaginationResponse{
		Page:  form.Page,
		Per:   form.PerPage,
		Total: total,
	}

	results := eventListMaker(allPeople)

	return results, paginations, nil
}

func recordDetailMarkProcessor(
	people models.FrequentCustomerPeople,
	form FrequentCustomerRecordMarkParams,
) (models.FrequentCustomerMark, *errors.Error) {
	var mark models.FrequentCustomerMark
	database.POSTGRES.FirstOrInit(&mark, models.FrequentCustomerMark{PersonID: people.PersonID})

	mark.Name = form.Name
	mark.Note = form.Note

	database.POSTGRES.Save(&mark)

	return mark, nil
}

func RecordsListHandler(c *gin.Context) {
	var form FrequentCustomerRecordParams

	if err := controllers.CheckRequestQuery(c, &form); err != nil {
		return
	}

	form.Normalize()

	list, paginations, errPtr := RecordListProcessor(form)

	if errPtr != nil {
		errors.ResponseError(c, *errPtr)
		return
	}

	controllers.SetPaginationToHeaderByStruct(c, paginations)

	c.JSON(200, list)
}

func fetchFreuqentCustomerPerson(c *gin.Context, form CompanyShopParams) (models.FrequentCustomerPeople, *errors.Error) {
	var result models.FrequentCustomerPeople
	fcID := c.Param("id")
	fcIDInt, err := strconv.Atoi(fcID)
	if err != nil {
		structedErr := errors.MakeInvalidaParamsError("id不为数字")
		return result, &structedErr
	}

	// 已经删除的也找出来
	database.POSTGRES.Unscoped().Preload("FrequentCustomerGroup").Preload("Event").First(&result, fcIDInt)

	if result.ID == 0 {
		structedErr := errors.MakeNotFoundError("未找到对应回头客")
		return result, &structedErr
	}

	if result.FrequentCustomerGroup.CompanyID != form.CompanyID {
		structedErr := errors.MakeInvalidaParamsError("没有操作回头客的权限")
		return result, &structedErr
	}

	return result, nil
}

func fetchFrequentCustomerEvent(c *gin.Context, form CompanyShopParams) (models.Event, *errors.Error) {
	var result models.Event
	evID := c.Param("event_id")
	evIDInt, err := strconv.Atoi(evID)
	if err != nil {
		structedErr := errors.MakeInvalidaParamsError("id不为数字")
		return result, &structedErr
	}

	database.POSTGRES.First(&result, evIDInt)

	if result.ID == 0 {
		structedErr := errors.MakeNotFoundError("未找到事件")
		return result, &structedErr
	}

	return result, nil
}

func RecordDetailHandler(c *gin.Context) {
	var form CompanyShopParams
	if err := controllers.CheckRequestQuery(c, &form); err != nil {
		return
	}

	person, errPtr := fetchFreuqentCustomerPerson(c, form)
	if errPtr != nil {
		errors.ResponseError(c, *errPtr)
		return
	}

	resultList := recordListMaker([]models.FrequentCustomerPeople{person})

	if len(resultList) > 0 {
		c.JSON(200, resultList[0])
	} else {
		errors.ResponseUnexpected(c, "转换数据格式错误")
	}
}

func RecordDetailListHandler(c *gin.Context) {
	var form FrequentCustomerRecordDetailParams
	if err := controllers.CheckRequestQuery(c, &form); err != nil {
		return
	}

	person, errPtr := fetchFreuqentCustomerPerson(c, form.CompanyShopParams)
	if errPtr != nil {
		errors.ResponseError(c, *errPtr)
		return
	}

	results, paginations, errPtr := recordDetailListProcessor(person, form)

	if errPtr != nil {
		errors.ResponseError(c, *errPtr)
		return
	}

	controllers.SetPaginationToHeaderByStruct(c, paginations)
	c.JSON(200, results)
}

func RecordDetailMarkHandler(c *gin.Context) {
	var form FrequentCustomerRecordMarkParams
	if err := controllers.CheckRequestBody(c, &form); err != nil {
		return
	}

	person, errPtr := fetchFreuqentCustomerPerson(c, form.CompanyShopParams)
	if errPtr != nil {
		errors.ResponseError(c, *errPtr)
		return
	}

	mark, errPtr := recordDetailMarkProcessor(person, form)
	if errPtr != nil {
		errors.ResponseError(c, *errPtr)
		return
	}

	if mark.PersonID != "" {
		database.POSTGRES.Table("customers").Where("person_id = ?", mark.PersonID).Updates(map[string]interface{}{"name": mark.Name, "note": mark.Note})
	}

	c.JSON(200, mark)
}

// RecordEventRemoveHandler 标记“不是TA”
func RecordEventRemoveHandler(c *gin.Context) {
	var form CompanyShopParams
	if err := controllers.CheckRequestBody(c, &form); err != nil {
		return
	}

	person, errPtr := fetchFreuqentCustomerPerson(c, form)
	if errPtr != nil {
		errors.ResponseError(c, *errPtr)
		return
	}

	event, errPtr := fetchFrequentCustomerEvent(c, form)
	if errPtr != nil {
		errors.ResponseError(c, *errPtr)
		return
	}

	if person.PersonID != event.PersonID {
		errors.ResponseNotPermitted(c)
		return
	}

	// 删除那一天的记录，就是找到frequent_customer_people进行删除
	var tobeDeleted models.FrequentCustomerPeople
	database.POSTGRES.Where("event_id = ?").First(&tobeDeleted)

	if tobeDeleted.ID != 0 {
		database.POSTGRES.Delete(&tobeDeleted)
	} else {
		errors.ResponseNotFound(c, "没有这条记录")
		return
	}

	c.JSON(200, event)
}
