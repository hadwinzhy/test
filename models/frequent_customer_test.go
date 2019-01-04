package models

import (
	"siren/configs"
	"siren/initializers"
	"siren/pkg/database"
	"testing"
	"time"
)

func TestSave(t *testing.T) {
	configs.ENV = "dev"
	initializers.ViperDefaultConfig()
	database.DBinit()
	day, _ := time.Parse("2006-01-02 00:00:00", time.Now().Format("2006-01-02 00:00:00"))

	tt := []struct {
		values FrequentCustomerCount
	}{
		{
			values: FrequentCustomerCount{
				PersonUUID:      "0987654",
				Day:             day,
				CapturedAt:      time.Now(),
				EventVisitCount: 24,
			},
		},
	}
	for _, t := range tt {
		database.POSTGRES.FirstOrCreate(&t.values)
	}
}
