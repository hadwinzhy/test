package workers

import (
	"siren/models"
	"siren/pkg/utils"
	"time"
)

func MallCountFrequentCustomerHandler(person models.FrequentCustomerPeople, groupID uint, capturedAt int64) {
	today := utils.CurrentDate(time.Now())
	hour := time.Unix(capturedAt, 0)
	err := updateFrequentCustomerReport(&person, groupID, today, hour)
	if err != nil {
		return
	}

	if person.IsHighFrequency() {
		updateFrequentCustomerHighTimeTable(groupID, today, capturedAt)
	}
}
