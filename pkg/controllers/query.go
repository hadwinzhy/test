package controllers

import (
	"strconv"
	"strings"
	"time"
)

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
