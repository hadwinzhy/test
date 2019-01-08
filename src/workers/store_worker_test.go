package workers

import (
	"siren/configs"
	"siren/initializers"
	"siren/pkg/database"
	"testing"
	"time"
)

func TestWorker(t *testing.T) {
	configs.ENV = "dev"
	initializers.ViperDefaultConfig()
	database.DBinit()
	// redis.Init()
	defer database.POSTGRES.Close()
	StoreFrequentCustomerHandler(241, 682, "2b68d6f2f6e684f19b0aaf1a9b61a883", time.Now().Unix())
}
