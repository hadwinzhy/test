package frequent_table

import (
	"log"
	"net/http"
	"siren/models"
	"siren/pkg/database"
	"siren/pkg/utils"

	"github.com/gin-gonic/gin"
)

func GetFrequentTableHandler(context *gin.Context) {
	var params GetFrequentTableParams

	if err := context.ShouldBindQuery(&params); err != nil {
		return
	}
	log.Println(params)

	var group models.FrequentCustomerGroup
	if dbError := database.POSTGRES.Where("company_id = ? AND shop_id = ?", params.CompanyID, params.ShopID).First(&group).Error; dbError != nil {
		MakeResponse(context, http.StatusBadRequest, dbError.Error())
		return
	}

	sql := `id, frequent_customer_group_id, sum(phase_one) as phase_one, sum(phase_two) as phase_two, sum(phase_three) as phase_three,
       sum(phase_four) as phase_four,sum(phase_five) as phase_five,sum(phase_six) as phase_six,sum(phase_seven) as phase_seven,
       sum(phase_eight) as phase_eight`

	var results []models.FrequentCustomerHighTimeTableSerializer
	for _, day := range weekDate() {
		day = utils.CurrentDate(day)
		log.Println(day)
		var data models.FrequentCustomerHighTimeTable
		query := database.POSTGRES.Model(&data).Where("frequent_customer_group_id = ?", group.ID)
		query = query.Select(sql).Where("date = ?", day).Group("id, frequent_customer_group_id")
		if dbError := query.First(&data).Error; dbError != nil {
			data = models.FrequentCustomerHighTimeTable{
				Date: day,
			}
		}
		data.Date = day
		log.Println(data)
		results = append(results, data.BasicSerializer())
	}

	MakeResponse(context, http.StatusOK, results)

}
