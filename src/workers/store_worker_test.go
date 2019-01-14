package workers

import (
	"siren/configs"
	"siren/initializers"
	"siren/pkg/database"
	"siren/pkg/utils"
	"testing"
	"time"
)

func TestWorker(t *testing.T) {
	configs.ENV = "dev"
	initializers.ViperDefaultConfig()
	database.DBinit()
	// redis.Init()
	defer database.POSTGRES.Close()
	StoreFrequentCustomerHandler(241, 682, "2b68d6f2f6e684f19b0aaf1a9b61a883", time.Now().Unix(), 12345)
}

func TestUpdateFrequentCustomerHighTimeTable(t *testing.T) {
	configs.ENV = "dev"
	initializers.ViperDefaultConfig()
	database.DBinit()
	// redis.Init()
	defer database.POSTGRES.Close()
	today := utils.CurrentDate(time.Now())
	captureAt := time.Now().Unix()
	updateFrequentCustomerHighTimeTable(2, today, captureAt)
}
