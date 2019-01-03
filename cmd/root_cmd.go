// Package cmd implement helper functions
package cmd

import (
	"fmt"
	"log"
	"os"
	"siren/initializers"
	"siren/pkg/database"
	"siren/pkg/logger"
	"siren/pkg/routers"
	"siren/src/kafka"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "siren",
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

		server := kafka.HeadCountProducer()
		defer func() {
			if err := server.Close(); err != nil {
				log.Println("Failed to close server", err)
			}
		}()

		go kafka.HeadCountCustomer()
		// Step 5: init router
		// go func() {
		r := gin.Default()
		routers.InitRouters(r)
		r.Run(":8088")
		// }()

		// Step 6: init rpc
		// startRpc()

	},
}

// Execute ...
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
