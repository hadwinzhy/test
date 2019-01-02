package initializers

import (
	"siren/configs"

	raven "github.com/getsentry/raven-go"
	"github.com/gin-contrib/sentry"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// SentryConfig ...
func SentryConfig() {
	if configs.ENV == "production" {
		raven.SetDSN(viper.GetString("sentry.dsn"))
	}
}

// AddSentryRecovery ...
func AddSentryRecovery(r *gin.Engine) {
	if configs.ENV == "production" {
		r.Use(sentry.Recovery(raven.DefaultClient, false))
	}
}
