package distributions

import "github.com/gin-gonic/gin"

func Register(r *gin.RouterGroup) {
	r.GET("/distributions", ListDistributionHandler)
}
