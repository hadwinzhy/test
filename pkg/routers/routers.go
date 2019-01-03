package routers

import (
	"siren/pkg/logger"
	"siren/pkg/middleware"
	"siren/src/kafka"

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
	group := r.Group("/v1/api")
	{
		kafka.Register(group)
	}

}
