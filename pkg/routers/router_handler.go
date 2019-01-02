package routers

import (
	"fmt"
	"net/http"
	"siren/configs"
	"siren/pkg/controllers/errors"
	"siren/pkg/logger"

	nativeErrors "errors"

	raven "github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
)

// router404Handler ...
func router404Handler(r *gin.Engine) {
	r.NoRoute(func(c *gin.Context) {
		errorCode := errors.ErrorRouterNotExist
		err := errors.Error{
			ErrorCode: errorCode,
			Detail:    "No Router for: " + c.Request.Method + " " + c.Request.URL.Path,
		}
		c.AbortWithStatusJSON(404, errors.ErrorResponse{
			Errors: []*errors.Error{&err},
		})
	})
}

// router500Handler ...
func router500Handler(r *gin.Engine) {
	if configs.ENV == "production" {
		r.Use(func(c *gin.Context) {
			defer func(c *gin.Context) {
				if rec := recover(); rec != nil {
					s := fmt.Sprintf("%v", rec)
					logger.Error("response", "fatal", s)
					raven.CaptureError(nativeErrors.New(s), nil)
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": rec,
					})
				}
			}(c)
			c.Next()
		})
	}
}
