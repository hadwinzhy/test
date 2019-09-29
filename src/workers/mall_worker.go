package workers

import (
	"siren/pkg/database"
	"siren/pkg/logger"
	"siren/pkg/utils"
	"time"

	"siren/venus/venus-model/models"
)

func MallCountFrequentCustomerHandler(person models.FrequentCustomerPeople, groupID uint, capturedAt int64) {
	today := utils.CurrentDate(time.Unix(capturedAt, 0))
	hour := time.Unix(capturedAt, 0).Truncate(time.Hour)
	nowUpdateFrequentCustomerReport := time.Now()
	logger.Info("statistic time", "update frequent customer report", "start")
	err := updateFrequentCustomerReport(&person, groupID, today, hour)
	if err != nil {
		return
	}
	logger.Info("statistic time", "update frequent customer report", "count time", time.Now().Sub(nowUpdateFrequentCustomerReport).Nanoseconds()/100000)

	nowHighFrequency := time.Now()
	logger.Info("statistic time", "update frequentCustomer High time table", "start")
	if person.IsHighFrequency(database.POSTGRES) {
		updateFrequentCustomerHighTimeTable(groupID, today, capturedAt)
	}
	logger.Info("statistic time", "update frequentCustomer High time table", "count time", time.Now().Sub(nowHighFrequency).Nanoseconds()/1000000)

}
