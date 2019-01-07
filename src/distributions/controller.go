package distributions

import (
	"siren/models"
	"siren/pkg/controllers"
	"siren/pkg/controllers/errors"
	"siren/pkg/database"

	"github.com/gin-gonic/gin"
)

func listDistributionProcessor(form ListDistributionParams) (DistributionOutput, *errors.Error) {
	fcGroups := models.FetchFrequentCustomerGroup(form.CompanyID, form.ShopID)

	var dataItems []models.FrequentCustomerReport
	if len(fcGroups) == 0 {

	} else {
		fromTime, toTime := form.GetFromAndToTime()

		database.POSTGRES.Model(&dataItems).
			Select("date, sum(high_frequency) AS high_frequency, sum(low_frequency) AS low_frequency, sum(new_comer) AS new_comer, sum(sum_interval) AS sum_interval, sum(sum_times) AS sum_times").
			Where("date >= ?", fromTime).
			Where("data < ?", toTime).
			Group("date").
			Find(&dataItems)
	}

	return DistributionOutput{}, nil
}

func ListDistributionHandler(c *gin.Context) {
	var form ListDistributionParams

	if err := controllers.CheckRequestQuery(c, &form); err != nil {
		return
	}

	form.Normalize()

	listDistributionProcessor(form)
}
