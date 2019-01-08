package routers

import (
	"siren/pkg/logger"
	"siren/pkg/middleware"
	"siren/src/frequent_rules"

	raven "github.com/getsentry/raven-go"
	"github.com/gin-contrib/sentry"

	"github.com/gin-gonic/gin"
)

// InitRouters ...
func InitRouters(r *gin.Engine) {
	middleware.UseCORS(r)

	// logger and sentry
	r.Use(logger.GinLogger())
	r.Use(sentry.Recovery(raven.DefaultClient, false))

	router404Handler(r)

	router500Handler(r)

	frequentGroup := r.Group("/v1/api")
	{
		frequent_rules.Register(frequentGroup)
	}

}
