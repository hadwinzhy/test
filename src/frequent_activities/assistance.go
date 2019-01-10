package frequent_activities

import (
	"siren/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func MakeResponse(context *gin.Context, code int, values interface{}) {
	context.JSON(code, values)
}

func dateHandler(date string) (string, string) {
	day := utils.TimestampToTime(date)
	left, _ := time.ParseInLocation("2006-01-02 00:00:00", day.Format("2006-01-02 15:04:05"), time.Local)
	return left.Format("2006-01-02 15:04:05"), day.Format("2006-01-02 15:04:05")
}

func monthHandler(date string) (string, string) {
	day := utils.TimestampToTime(date)
	left := day.AddDate(0, -1, 0).Format("2006-01-02 15:04:05")
	return left, day.Format("2006-01-02 15:04:05")

}
