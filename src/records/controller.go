package records

import (
	"siren/models"
	"siren/pkg/controllers"
	"siren/pkg/controllers/errors"
	"siren/pkg/database"

	"github.com/gin-gonic/gin"
)

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
	result := make([]FrequentCustomerRecord, len(peopleList))

	var shopIDSlice []uint

	for i, people := range peopleList {
		shopIDSlice = append(shopIDSlice, people.Event.ShopID)
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

	return result, paginations, nil
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
