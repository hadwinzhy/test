package distributions

import (
	"fmt"
	"siren/models"
	"siren/pkg/controllers"
	"siren/pkg/controllers/errors"
	"siren/pkg/database"

	"github.com/gin-gonic/gin"
)

func proportionCounter(divider uint, dividedBy uint) string {
	if dividedBy == 0 {
		return "0.00%"
	} else {
		percentage := float64(divider) / float64(dividedBy)
		return fmt.Sprintf("%.2f%%", percentage*100)
	}
}

func averageCounter(divider uint, dividedBy uint) string {
	if dividedBy == 0 {
		return "0.0"
	} else {
		average := float64(divider) / float64(dividedBy)
		return fmt.Sprintf("%.1f", average)
	}
}

func reportOutputMapper(item models.FrequentCustomerReport) DistributionOutput {
	allCount := item.HighFrequency + item.LowFrequency + item.NewComer
	frequentCustomers := item.HighFrequency + item.LowFrequency
	return DistributionOutput{
		Date: item.Hour.Format("2006-01-02 15:04:05"),
		High: valueProportionPair{
			Count:      item.HighFrequency,
			Proportion: proportionCounter(item.HighFrequency, allCount),
		},
		Low: valueProportionPair{
			Count:      item.LowFrequency,
			Proportion: proportionCounter(item.LowFrequency, allCount),
		},
		New: valueProportionPair{
			Count:      item.NewComer,
			Proportion: proportionCounter(item.NewComer, allCount),
		},
		AverageFrequency: averageCounter(frequentCustomers, item.SumTimes),
		AverageInterval:  averageCounter(frequentCustomers, item.SumInterval),
	}
}

func listDistributionProcessor(form ListDistributionParams) ([]DistributionOutput, *errors.Error) {
	fcGroups := models.FetchFrequentCustomerGroup(form.CompanyID, form.ShopID)
	fromTime, toTime := form.GetFromAndToTime()

	var dataItems []models.FrequentCustomerReport
	var groupIDs []uint
	for _, group := range fcGroups {
		groupIDs = append(groupIDs, group.ID)
	}

	if len(groupIDs) == 0 {

	} else {
		database.POSTGRES.Model(&dataItems).
			Select("date_trunc('day', hour) AS hour, sum(high_frequency) AS high_frequency, sum(low_frequency) AS low_frequency, sum(new_comer) AS new_comer, sum(sum_interval) AS sum_interval, sum(sum_times) AS sum_times").
			Where("hour >= ?", fromTime).
			Where("hour <= ?", toTime).
			Where("frequent_customer_group_id in (?)", groupIDs).
			Group("1").
			Order("hour " + form.SortBy).
			Find(&dataItems)
	}

	// insert missing 来扩充dataItems
	dataItems, _ = models.FrequentCustomerReports(dataItems).InsertMissing(form.Period, fromTime, toTime, form.SortBy)

	// 每一个元素进行一波计算， 算点比例和平均值
	results := make([]DistributionOutput, len(dataItems))

	sumReport := models.FrequentCustomerReport{}

	for i, item := range dataItems {
		results[i] = reportOutputMapper(item)

		sumReport.HighFrequency += item.HighFrequency
		sumReport.LowFrequency += item.LowFrequency
		sumReport.NewComer += item.NewComer
		sumReport.SumInterval += item.SumInterval
		sumReport.SumTimes += item.SumTimes
	}

	// 有的话要计算平均值和总和值

	if form.ReturnALL == "all_list" {
		return results, nil
	}

	sumOutput := reportOutputMapper(sumReport)
	sumOutput.Date = "合计"

	if form.ReturnALL == "all_count" {
		return []DistributionOutput{sumOutput}, nil
	}

	if form.ReturnALL == "all_list_average_count" {
		itemCount := uint(len(dataItems))
		averageReport := models.FrequentCustomerReport{
			HighFrequency: sumReport.HighFrequency / itemCount,
			LowFrequency:  sumReport.LowFrequency / itemCount,
			NewComer:      sumReport.NewComer / itemCount,
			SumInterval:   sumReport.SumInterval / itemCount,
			SumTimes:      sumReport.SumTimes / itemCount,
		}
		averageOutput := reportOutputMapper(averageReport)
		averageOutput.Date = "平均"
		results = append(results, averageOutput, sumOutput)
	}
	return results, nil
}

func ListDistributionHandler(c *gin.Context) {
	var form ListDistributionParams

	if err := controllers.CheckRequestQuery(c, &form); err != nil {
		return
	}

	form.Normalize()

	results, _ := listDistributionProcessor(form)

	c.JSON(200, results)
}
