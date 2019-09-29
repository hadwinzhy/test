package frequent_rules

import (
	"log"
	"net/http"
	"siren/venus/venus-model/models"
	"siren/pkg/controllers/errors"
	"siren/pkg/database"

	"github.com/gin-gonic/gin"
)

// createFrequentRule by params
func PostFrequentRuleHandler(context *gin.Context) {
	var params PostRules
	if err := context.ShouldBindJSON(&params); err != nil {
		return
	}

	var (
		ok  bool
		err *errors.Error
	)
	if ok, err = params.IsSuitableParam(); !ok {
		errors.ResponseError(context, *err)
		return
	}
	log.Println(params, "siren")

	low := params.JsonbLowHandler()
	high := params.JsonbHighHandler()

	var rules models.FrequentCustomerRule
	if dbError := database.POSTGRES.Where("company_id = ?", params.CompanyID).First(&rules).Error; dbError != nil {

		rules = models.FrequentCustomerRule{
			CompanyID:     params.CompanyID,
			LowFrequency:  low,
			HighFrequency: high,
			Limit:         params.Limit,
		}
		database.POSTGRES.Save(&rules)
	} else {
		if rules.ID != 0 {
			database.POSTGRES.Model(&rules).Updates(map[string]interface{}{
				"low_frequency":  low,
				"high_frequency": high,
				"limit":          params.Limit,
			})

		}
	}
	MakeResponse(context, http.StatusOK, rules.BasicSerializer())

}

// getOneFrequentRule by company id
func GetAllFrequentRulesHandler(context *gin.Context) {
	var params string
	params = context.Query("company_id")
	log.Println("params in siren", params)
	var results models.FrequentCustomerRules
	if dbError := database.POSTGRES.Where("company_id = ?", params).Find(&results).Error; len(results) == 0 || dbError != nil {
		var frequency models.FrequentCustomerRule
		MakeResponse(context, http.StatusOK, frequency.ReadableRule())
		return
	}
	var resultsSerializer []models.FrequentCustomerRuleBasicSerializer
	for _, result := range results {
		resultsSerializer = append(resultsSerializer, result.BasicSerializer())
	}
	MakeResponse(context, http.StatusOK, resultsSerializer)

}

// deleteOneFrequentRule by company_id
func DeleteOneFrequentRuleHandler(context *gin.Context) {
	var params string
	params = context.Param("company_id")
	log.Println(params)
	var result models.FrequentCustomerRule
	if dbError := database.POSTGRES.Where("company_id = ?", params).First(&result).Error; dbError != nil {
		err := &errors.Error{
			ErrorCode: errors.ErrorCode{
				HTTPStatus: 400,
				Code:       400,
				Title:      "record not found",
				TitleZH:    "记录未找到",
			},
			Detail: "记录未找到",
		}
		errors.ResponseError(context, *err)
		return
	}
	database.POSTGRES.Delete(&result)
	MakeResponse(context, http.StatusOK, result.BasicSerializer())

}
