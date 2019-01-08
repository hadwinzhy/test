package frequent_rules

import (
	"log"
	"net/http"
	"siren/models"
	"siren/pkg/database"

	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
)

// createFrequentRule by params
func PostFrequentRuleHandler(context *gin.Context) {
	var params PostRules
	if err := context.ShouldBindJSON(&params); err != nil {
		return
	}

	if !params.IsSuitableParam() {
		MakeResponse(context, http.StatusBadRequest, errors.New("params are not correct").Error())
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
	var results models.FrequentCustomerRules
	if dbError := database.POSTGRES.Where("company_id = ?", params).Find(&results).Error; dbError != nil {
		MakeResponse(context, http.StatusBadRequest, errors.New("records not found").Error())
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
		MakeResponse(context, http.StatusBadRequest, errors.New("records not found").Error())
		return
	}
	database.POSTGRES.Delete(&result)
	MakeResponse(context, http.StatusOK, result.BasicSerializer())

}
