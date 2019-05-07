package workers

import (
	"siren/pkg/database"
	"siren/pkg/utils"
	"time"

	"bitbucket.org/readsense/venus-model/models"
)

func MallCountFrequentCustomerHandler(person models.FrequentCustomerPeople, groupID uint, capturedAt int64) {
	today := utils.CurrentDate(time.Unix(capturedAt, 0))
	hour := time.Unix(capturedAt, 0).Truncate(time.Hour)
	err := updateFrequentCustomerReport(&person, groupID, today, hour)
	if err != nil {
		return
	}

	if person.IsHighFrequency(database.POSTGRES) {
		updateFrequentCustomerHighTimeTable(groupID, today, capturedAt)
	}
}
