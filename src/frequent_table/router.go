package frequent_table

import "github.com/gin-gonic/gin"

func Register(r *gin.RouterGroup) {
	r.GET("/frequent_table", GetFrequentTableHandler)
}
