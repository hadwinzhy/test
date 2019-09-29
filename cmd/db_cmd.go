package cmd

import (
	"errors"
	"fmt"
	"reflect"
	"siren/initializers"
	"siren/venus/venus-model/models"
	"siren/pkg/database"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(dbCmd)
}

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Database modifications.",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires at least one arguments")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

		// Step 1: init some data
		initializers.ViperDefaultConfig()
		// Step 2: init database
		database.DBinit()
		defer database.POSTGRES.Close()

		switch args[0] {
		case "create":
			CreateTables()
			fmt.Println("Create Database Done!")
		case "migrate":
			MigrateTables()
			fmt.Println("Migration Done!")

		case "bitmap":
			createBitMaps()
		case "addIndexForPeople":
			addIndexForPeople()
		}
	},
}

// CreateTables will create Table if needed
func CreateTables() {
	for _, model := range models.GetModels() {
		if !database.POSTGRES.HasTable(model) {
			fmt.Print("--- Create Table ")
			fmt.Println(reflect.TypeOf(model))
			database.POSTGRES.CreateTable(model)
		}
	}
}

// MigrateTables will migrate Table if needed
func MigrateTables() {
	for _, model := range models.GetModels() {
		fmt.Print("--- Migrate Table")
		database.POSTGRES.AutoMigrate(model)
	}
}

func createBitMaps() {
	// newBitMap := models.FrequentCustomerPeopleBitMap{
	// 	PersonID: "a random id",
	// 	BitMap:   "00000000000000000000000000000010",
	// }

	// database.POSTGRES.Save(&newBitMap)
}

func addIndexForPeople() {
	database.POSTGRES.Model(&models.FrequentCustomerPeople{}).AddIndex("idx_group_frequent_customer", "frequent_customer_group_id", "is_frequent_customer")
}
