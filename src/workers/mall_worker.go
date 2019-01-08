package workers

import (
	"siren/models"
	"siren/pkg/utils"
	"time"
)

func MallCountFrequentCustomerHandler(person models.FrequentCustomerPeople, groupID uint, capturedAt int64) {
	today := utils.CurrentDate(time.Now())
	err := updateFrequentCustomerReport(&person, groupID, today)
	if err != nil {
		return
	}

	if person.IsHighFrequency() {
		updateFrequentCustomerHighTimeTable(groupID, today, capturedAt)
	}
}
