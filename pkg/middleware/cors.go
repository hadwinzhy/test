package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// UseCORS will set cors config
func UseCORS(r *gin.Engine) {
	// r.Use(cors.Default())
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AddAllowHeaders("Authorization")
	config.AddAllowHeaders("X-Requested-With")

	config.AddAllowMethods("OPTIONS")
	config.AddAllowMethods("DELETE")
	config.AddAllowMethods("PATCH")

	config.AddExposeHeaders("X-Total-Count")
	config.AddExposeHeaders("X-Current-Page")
	config.AddExposeHeaders("X-Per-Page")

	r.Use(cors.New(config))
}
