package frequent_table

import (
	"siren/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func MakeResponse(context *gin.Context, code int, values interface{}) {
	context.JSON(code, values)
}

func weekDate(date string) []time.Time {

	now := utils.TimestampToTime(date)
	day, _ := time.ParseInLocation("2006-01-02", now.Format("2006-01-02"), time.Local)
	var week []time.Time
	week = append(week, day)
	for i := 1; i < 7; i++ {
		temp := now.AddDate(0, 0, -(i))
		tempDay, _ := time.ParseInLocation("2006-01-02", temp.Format("2006-01-02"), time.Local)
		week = append(week, tempDay)
	}
	return week
}
