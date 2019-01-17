package frequent_activities

import "github.com/gin-gonic/gin"

func Register(r *gin.RouterGroup) {
	r.GET("/frequent_activities", GetFrequentActivitiesHandler)
}
