package frequent_table

import (
	"log"
	"net/http"
	"siren/models"
	"siren/pkg/database"
	"siren/pkg/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GetFrequentTableHandler(context *gin.Context) {
	var params GetFrequentTableParams

	if err := context.ShouldBindQuery(&params); err != nil {
		return
	}
	log.Println("siren", params)

	var group models.FrequentCustomerGroup
	if params.ShopID != 0 {
		if dbError := database.POSTGRES.Where("company_id =? AND shop_id = ?", params.CompanyID, params.ShopID).First(&group).Error; dbError != nil {
			//MakeResponse(context, http.StatusBadRequest, dbError.Error())
			//return
			// 未找到，返回空值
		}
	} else {
		if dbError := database.POSTGRES.Where("company_id = ?", params.CompanyID).First(&group).Error; dbError != nil {
			//MakeResponse(context, http.StatusBadRequest, dbError.Error())
			//return
			// 未找到，返回空值
		}
	}

	sql := `id, frequent_customer_group_id, sum(phase_one) as phase_one, sum(phase_two) as phase_two, sum(phase_three) as phase_three,
       sum(phase_four) as phase_four,sum(phase_five) as phase_five,sum(phase_six) as phase_six,sum(phase_seven) as phase_seven,
       sum(phase_eight) as phase_eight`

	var results []models.FrequentCustomerHighTimeTableSerializer

	var date string
	if params.Date == "" {
		date = strconv.Itoa(int(time.Now().Unix()))
	} else {
		date = params.Date
	}

	for _, day := range weekDate(date) {
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
