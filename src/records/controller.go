package records

import (
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
		result[i] = FrequentCustomerRecord{
			FrequentCustomerPersonID: people.ID,
			FirstCaptureURL:          people.Event.OriginalFace,
			Name:                     "", // 要根据person_id对应的去取，作标记的时候再做
			CaptureAt:                people.Event.CaptureAt,
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

	// manually load mark
	if len(personIDSlice) > 0 {
		personIDMap := make(map[string]models.FrequentCustomerMark)

		var marks []models.FrequentCustomerMark

		database.POSTGRES.Where("person_id in (?)", personIDSlice).Find(&marks)

		for i := range marks {
			personIDMap[marks[i].PersonID] = marks[i]
		}

		for i := range result {
			result[i].Note = personIDMap[result[i].personID].Note
			result[i].Name = personIDMap[result[i].personID].Name
		}
	}

	return result
}

func eventListMaker(events []models.Event) []SingleEventRecord {
	result := make([]SingleEventRecord, len(events))

	var shopIDSlice []uint

	for i, event := range events {
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
	fcGroups := models.FetchFrequentCustomerGroup(form.CompanyID, form.ShopID)

	var groupIDs []uint
	for _, group := range fcGroups {
		groupIDs = append(groupIDs, group.ID)
	}

	fromTime, toTime := form.GetFromAndToTime()

	// 回头group范围，根据personid聚合取最新的
	query := database.POSTGRES.Model(&models.FrequentCustomerPeople{}).
		Preload("Event").
		Select("person_id, frequent_customer_group_id, max(id) AS id, max(event_id) AS event_id, max(frequency) AS frequency, max(hour) AS hour")

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
	var events []models.Event
	query := database.POSTGRES.Model(&models.Event{}).
		Where("capture_at >= ?", fromTime).
		Where("capture_at <= ?", toTime).
		Where("person_id = ?", people.PersonID)

	if form.ShopID != 0 {
		query = query.Where("shop_id = ?", form.ShopID)
	}

	var total int
	query.Count(&total)

	query.Order("capture_at desc").
		Limit(form.PerPage).
		Offset(form.PerPage * (form.Page - 1)).
		Find(&events)

	paginations = controllers.PaginationResponse{
		Page:  form.Page,
		Per:   form.PerPage,
		Total: total,
	}

	results := eventListMaker(events)

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

	database.POSTGRES.Preload("FrequentCustomerGroup").Preload("Event").First(&result, fcIDInt)

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
	if err := controllers.CheckRequestQuery(c, &form); err != nil {
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

	c.JSON(200, mark)
}
