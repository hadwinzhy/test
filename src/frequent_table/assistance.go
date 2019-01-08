package frequent_table

import (
	"time"

	"github.com/gin-gonic/gin"
)

func MakeResponse(context *gin.Context, code int, values interface{}) {
	context.JSON(code, values)
}

func weekDate() []time.Time {
	now := time.Now()
	day, _ := time.Parse("2006-01-02", now.Format("2006-01-02"))
	var week []time.Time
	week = append(week, day)
	for i := 1; i < 7; i++ {
		temp := now.AddDate(0, 0, -(i))
		tempDay, _ := time.Parse("2006-01-02", temp.Format("2006-01-02"))
		week = append(week, tempDay)
	}
	return week
}
