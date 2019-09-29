package controllers

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type TimeValue struct {
	Time  string `json:"time"`
	Value int    `json:"value"`
}

type FromToParam struct {
	FromDate string `form:"from_date" json:"from_date"`
	ToDate   string `form:"to_date" json:"to_date"`
}

func (form *FromToParam) Normalize() {
	today := time.Now().Format("2006-01-02")
	timestamp, _ := time.ParseInLocation("2006-01-02", today, time.Local)
	if form.FromDate == "" {
		form.FromDate = strconv.FormatInt(timestamp.Unix(), 10)
	}
	if form.ToDate == "" {
		form.ToDate = strconv.FormatInt(timestamp.AddDate(0, 0, 1).Unix()-1, 10)
	}
}

func (form *FromToParam) GetFromAndToTime() (fromTime time.Time, toTime time.Time) {
	if form.FromDate != "" {
		fromTime = TimestampToTime(form.FromDate)
		year, month, day := fromTime.Date()
		fromTime = time.Date(year, month, day, fromTime.Hour(), 0, 0, 0, time.Local)

		if form.ToDate == "" {
			toTime = time.Now()
		} else {
			toTime = TimestampToTime(form.ToDate)
		}
	} else {
		todayStr := time.Now().Format("2006-01-02")
		fromTime, _ = time.ParseInLocation("2006-01-02", todayStr, time.Local)
		toTime = time.Now()
	}

	return fromTime, toTime
}

type PeriodParam struct {
	Period string `form:"period" binding:"omitempty,eq=month|eq=day|eq=hour|eq=week|eq=year"`
}

func TimestampToTime(from string) time.Time {
	if id := strings.Index(from, "."); id != -1 {
		from = from[:id]
	}
	i, _ := strconv.ParseInt(from, 10, 64)
	value := time.Unix(i, 0)
	return value
}

func QueryPeriod(query *gorm.DB, fromDate string, toDate string, dbKey string) (*gorm.DB, time.Time, time.Time) {
	fromTime := TimestampToTime(fromDate)
	var toTime time.Time
	if toDate == "" {
		toTime = time.Now()
	} else {
		toTime = TimestampToTime(toDate)
	}

	return query.Where((dbKey + " > ? AND " + dbKey + " < ?"), fromTime, toTime), fromTime, toTime

}

func QueryToTimeValue(baseQuery *gorm.DB, period string, sortBy string, orderBy string, table string) []TimeValue {
	periodQuery := ""
	switch period {
	case "month":
		periodQuery = fmt.Sprintf("to_char(%s.%s, '%s')", table, orderBy, "YYYY/MM")
	case "week":
		periodQuery = fmt.Sprintf("to_char(%s.%s, '%s')", table, orderBy, "YY/WW周")
	case "hour":
		periodQuery = fmt.Sprintf("to_char(%s.%s, '%s')", table, orderBy, "HH24:00")
	case "minute":
		periodQuery = fmt.Sprintf("to_char(%s.%s, '%s')", table, orderBy, "HH24:MI")
	default:
		periodQuery = fmt.Sprintf("to_char(%s.%s, '%s')", table, orderBy, "MM/DD")
	}
	order := fmt.Sprintf("time %s", sortBy)
	rows, err := baseQuery.Select(fmt.Sprintf("%s AS time, count(*) AS value", periodQuery)).Group(periodQuery).Order(order).Rows()
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	response := []TimeValue{}
	for rows.Next() {
		var timeValue TimeValue
		if err := rows.Scan(&timeValue.Time, &timeValue.Value); err != nil {
			log.Fatal(err)
		}
		response = append(response, timeValue)
	}
	return response
}
