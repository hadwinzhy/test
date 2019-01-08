package frequent_rules

import "github.com/gin-gonic/gin"

func MakeResponse(context *gin.Context, code int, values interface{}) {
	context.JSON(
		code, values)
}
