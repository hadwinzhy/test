package database

import (
	"fmt"
	"siren/configs"

	"github.com/jinzhu/gorm"
	// import postgres driver
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// POSTGRES is Singleton
var POSTGRES *gorm.DB

// DBinit connect the postgres database
func DBinit() *gorm.DB {
	conn, err := gorm.Open("postgres", getPGConnectionString())
	if err != nil {
		panic("failed to connect database" + err.Error())
	}

	if configs.ENV != "production" {
		conn.DB().SetMaxIdleConns(3)
		conn.LogMode(true)
	}

	POSTGRES = conn
	return POSTGRES
}

// getPGConnectionString by env
func getPGConnectionString() string {
	var host, sslmode, user, password, port, dbname string
	host = configs.FetchFieldValue("PGHOST")
	sslmode = configs.FetchFieldValue("PGSSLMODE")
	user = configs.FetchFieldValue("PGUSER")
	password = configs.FetchFieldValue("PGPASSWORD")
	port = configs.FetchFieldValue("PGPORT")
	dbname = configs.FetchFieldValue("PGDBNAME")

	fmt.Println("Host", host, "sslMode", sslmode, "user", user, "password", password, "port", port, "dbname", dbname)

	return "host=" + host + " port=" + port + " user=" + user + " password=" + password + " dbname=" + dbname + " sslmode=" + sslmode
}
