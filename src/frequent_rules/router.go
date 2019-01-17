package frequent_rules

import "github.com/gin-gonic/gin"

func Register(r *gin.RouterGroup) {
	r.POST("/frequent_rules", PostFrequentRuleHandler)
	r.GET("/frequent_rules", GetAllFrequentRulesHandler)
	r.DELETE("/frequent_rules/:company_id", DeleteOneFrequentRuleHandler)
}
