package cmd

import (
	"os"
	"siren/initializers"
	"siren/pkg/database"
	"siren/pkg/logger"
	"siren/src/kafka"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(onlyConsumerCmd)
}

var onlyConsumerCmd = &cobra.Command{
	Use:   "onlyconsumer",
	Short: "Siren 0.0.1",
	Long:  "Web server based on Golang with simplicity and performance.",
	Run: func(cmd *cobra.Command, args []string) {
		// Step 1: init some data
		initializers.ViperDefaultConfig()
		// Step 2: init database
		database.DBinit()
		// redis.Init()
		defer database.POSTGRES.Close()

		// Step 4: init logger
		logger.Init()

		// Step 5: Sentry
		initializers.SentryConfig()

		if os.Getenv("KAFKA_CONSUMER_SWITCH") != "OFF" {
			go kafka.CountFrequentConsumer()
		}

		for {
		}
		// Step 5: init router
		// go func() {

		// }()

		// Step 6: init rpc
		// startRpc()

	},
}
