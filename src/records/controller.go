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

	// 回头group范围
	query := database.POSTGRES.Model(&models.FrequentCustomerPeople{}).
		Preload("Event").
		Select("person_id, max(id) AS id, max(event_id) AS event_id, max(frequency) AS frequency, max(capture_at) AS capture_at")

	query = query.Where("frequent_customer_group_id in (?)", groupIDs)

	// 时间范围
	query = query.Where("hour >= ?", fromTime).Where("hour < ?", toTime)

	query = query.Group("person_id")

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

	return []FrequentCustomerRecord{}, controllers.PaginationResponse{}, nil
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
