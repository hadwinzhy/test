package initializers

import (
	"fmt"

	"github.com/spf13/viper"
)

// ViperDefaultConfig ...
func ViperDefaultConfig() {
	viper.SetConfigName("settings")                   // name of config file (without extension)
	viper.AddConfigPath("$GOPATH/src/siren/configs/") // path to look for the config file in
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig() // Find and read the config file

	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s ", err))
	}
}

// ViperReadFile will read file
func ViperReadFile(viper *viper.Viper, fileName string, fileType string) {
	viper.SetConfigName(fileName)     // name of config file (without extension)
	viper.AddConfigPath("./configs/") // path to look for the config file in
	viper.SetConfigType(fileType)
	err := viper.ReadInConfig() // Find and read the config file

	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s ", err))
	}
}
