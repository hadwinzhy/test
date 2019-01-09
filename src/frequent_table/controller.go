package frequent_table

import (
	"fmt"
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

	// max
	sql2 := fmt.Sprintf(`select max(phase_one) as phase_one, max(phase_two) as phase_two, max(phase_three) as phase_three, max(phase_four) as phase_four,
			max(phase_five) as phase_five, max(phase_six) as phase_six, max(phase_seven) as phase_seven, max(phase_eight) as phase_eight from
			frequent_customer_high_time_tables where frequent_customer_group_id = %d and date between %s and %s
`, group.ID, weekDate()[6], weekDate()[0])
	var maxCount maxNumber
	database.POSTGRES.Raw(sql2).Scan(&maxCount)
	MakeResponse(context, http.StatusOK, results)

}

type maxNumber struct {
	PhaseOne   int `json:"phase_one"`
	PhaseTwo   int `json:"phase_two"`
	PhaseThree int `json:"phase_three"`
	PhaseFour  int `json:"phase_four"`
	PhaseFive  int `json:"phase_five"`
	PhaseSix   int `json:"phase_six"`
	PhaseSeven int `json:"phase_seven"`
	PhaseEight int `json:"phase_eight"`
}
