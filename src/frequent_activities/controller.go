package frequent_activities

import (
	"log"
	"net/http"
	"siren/models"
	"siren/pkg/database"

	"github.com/gin-gonic/gin"
)

func GetFrequentActivitiesHandler(context *gin.Context) {
	/*
	 */
	var params GetFrequentActivitiesParams
	if err := context.ShouldBindQuery(&params); err != nil {
		return
	}

	log.Println("params", params)
	var groups []models.FrequentCustomerGroup
	if params.Type == "store" {
		if params.ShopID == 0 {
			if dbError := database.POSTGRES.Where("company_id = ?", params.CompanyID).Find(&groups).Error; dbError != nil {
				// 未找到，也返回值，只不过是空值
			}
		} else {
			if dbError := database.POSTGRES.Where("company_id = ? AND shop_id = ?", params.CompanyID, params.ShopID).Find(&groups).Error; dbError != nil {
				// 未找到，也返回值，只不过是空值
			}
		}
	} else {
		if dbError := database.POSTGRES.Where("company_id = ?", params.CompanyID).Find(&groups).Error; dbError != nil {
			// 未找到，也返回值，只不过是空值
		}
	}

	var groupIDs []uint
	if len(groups) != 0 {
		for _, group := range groups {
			groupIDs = append(groupIDs, group.ID)
		}
	}

	left, right := dateHandler(params.Date)
	// 日期，不传则为当天时间
	var results models.FrequentCustomerPeoples
	query := database.POSTGRES.
		Where("frequent_customer_group_id in (?)", groupIDs).
		Where("date BETWEEN ? AND ?", left, right).Where("is_frequent_customer = ?", "true")
	query.Find(&results)

	resultsReport := results.Activities()
	log.Println("report", resultsReport)

	beforeMonth, day := monthHandler(params.Date)
	queryLowHigh := database.POSTGRES.
		Where("frequent_customer_group_id in (?)", groupIDs).
		Where("date BETWEEN ? AND ?", beforeMonth, day).Where("is_frequent_customer = ?", "true")

	var rules models.FrequentCustomerRule
	if dbError := database.POSTGRES.Where("company_id = ?", params.CompanyID).First(&rules).Error; dbError != nil {

	}
	var lowHigh models.FrequentCustomerPeoples
	queryLowHigh.Find(&lowHigh)
	var lowHighResult []*models.OneStatic
	lowHighResult = lowHigh.FrequentMonthStatic(rules)
	log.Println("lowHighResult", lowHighResult)

	var allResult models.FrequentCount
	allResult.Vitality = make(map[string]interface{})
	allResult.Vitality["frequency"] = resultsReport
	allResult.Vitality["interval"] = lowHighResult
	MakeResponse(context, http.StatusOK, allResult)

}
