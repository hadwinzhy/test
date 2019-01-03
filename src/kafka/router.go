package kafka

import "github.com/gin-gonic/gin"

func Register(r *gin.RouterGroup) {

	r.POST("/head_count", PostProducerDataHandler)
}
