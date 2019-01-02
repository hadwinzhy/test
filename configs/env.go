package configs

import (
	"os"

	"github.com/spf13/viper"
)

// 1. titan server
// 2. qiniu configs
// 3. db configs
// map[string]string : key stands for env var, value stands for setting key
func launchCollection() map[string]string {
	var collectionEnv = make(map[string]string)
	collectionEnv = map[string]string{
		"TitanHOST":        "titan.host",
		"TitanAPIID":       "titan.apiid",
		"TitanAPISecret":   "titan.apisecret",
		"titan_host":       "titan.host",      //compatibility
		"titan_api_id":     "titan.apiid",     //compatibility
		"titan_api_secret": "titan.apisecret", //compatibility
		"PGHOST":           "db.host",
		"PGSSLMODE":        "db.sslmode",
		"PGUSER":           "db.user",
		"PGPASSWORD":       "db.password",
		"PGPORT":           "db.port",
		"PGDBNAME":         "db.dbname",
	}
	return collectionEnv
}

func FetchFieldValue(field string) string {
	collections := launchCollection()
	var envResult string
	if value, ok := collections[field]; ok {
		if envValue := os.Getenv(field); envValue != "" {
			envResult = envValue
		} else {
			envResult = viper.GetString(ENV + "." + value)
		}
	}
	return envResult
}
