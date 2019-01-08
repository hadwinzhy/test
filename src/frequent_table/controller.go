package frequent_table

import (
	"log"
	"net/http"
	"siren/models"
	"siren/pkg/database"

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

	var table models.FrequentCustomerHighTimeTables
	query := database.POSTGRES.Model(&table).Where("frequent_customer_group_id = ?", group.ID)

	sql := `select id, sum(phase_one) as phase_one, sum(phase_two) as phase_two, sum(phase_three) as phase_three,
       sum(phase_four) as phase_four,sum(phase_five) as phase_five,sum(phase_six) as phase_six,sum(phase_seven) as phase_seven,
       sum(phase_eight) as phase_eight
       from frequent_customer_high_time_tables`

	var results models.FrequentCustomerHighTimeTables
	for _, day := range weekDate() {
		query = query.Raw(sql).Where("date = ?", day)
		var data models.FrequentCustomerHighTimeTable
		if dbError := query.First(&data).Error; dbError != nil {
			MakeResponse(context, http.StatusBadRequest, dbError.Error())
			return
		}
		results = append(results, data)
	}
	MakeResponse(context, http.StatusOK, results)

}
