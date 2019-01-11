package records

import "github.com/gin-gonic/gin"

func Register(r *gin.RouterGroup) {
	r.GET("/records", RecordsListHandler)
	r.GET("/records/:id", RecordDetailHandler)

	r.GET("/records/:id/events", RecordDetailListHandler)
}
