// Package main initialize the app
package main

import (
	"fmt"
	"siren/cmd"
	"siren/configs"
)

// 设置全局环境变量
var (
	Env string
)

// @title 智能零售 API
// @version 1.0
// @description This is a server for retail.
// @termsOfService http://readsense.cn/term

// @contact.name API Support
// @contact.email huiyun.zheng@readsense.cn

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host lightyear.readsense.cn:8080
// @BasePath /v1/api

// @securityDefinitions.apiKey Bearer
// @type apiKey
// @in header
// @name Authorization

// main initialize the app with necessary commands.

func main() {

	if Env != "" {
		configs.ENV = Env
	} else {
		configs.ENV = "dev"
	}

	fmt.Println("Running on " + configs.ENV)
	cmd.Execute()
}
