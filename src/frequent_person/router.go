package frequent_person

import "github.com/gin-gonic/gin"

func RegisterJudgeFrequentPerson(r *gin.RouterGroup) {
	r.GET("/is_frequent_person", JudgeFrequentPersonHandler)
}
