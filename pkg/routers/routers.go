package routers

import (
	"net/http"
	"siren/pkg/logger"
	"siren/pkg/middleware"
	"siren/src/frequent_rules"
	"siren/src/frequent_table"

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

	r.GET("/heart_beat", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"ping": "pong",
		})
	})
	frequentGroup := r.Group("/v1/api")
	{
		frequent_rules.Register(frequentGroup)
		frequent_table.Register(frequentGroup)
	}

}
