package utils

import (
	"strconv"
	"strings"
	"time"
)

func CurrentDate(captureAt time.Time) time.Time {
	return time.Date(captureAt.Year(), captureAt.Month(), captureAt.Day(), 0, 0, 0, 0, time.Local)
}

func GetDurationByPeriod(period string) time.Duration {
	duration := time.Duration(time.Hour)
	if period == "hour" {
		duration = time.Duration(time.Hour)
	} else if period == "day" {
		duration = time.Duration(time.Hour * 24)
	} else if period == "week" || period == "offset_week" {
		duration = time.Duration(time.Hour * 24 * 7)
	} else if period == "month" || period == "offset_month" {
		duration = time.Duration(time.Hour * 24 * 30) // ignore 闰年
	} else if period == "year" || period == "offset_year" {
		duration = time.Duration(time.Hour * 24 * 365)
	}
	return duration
}

func GetSubByPeriod(period string, fromTime time.Time, toTime time.Time, duration time.Duration) int {
	sub := int((toTime.Sub(fromTime)-time.Millisecond)/duration) + 1
	// if period == "hour" && fromTime.Hour() != toTime.Hour() {
	// 	sub++
	// }
	// if period == "day" && fromTime.Day() != toTime.Day() {
	// 	sub++
	// }
	if period == "week" {
		fw := fromTime.Weekday()
		tw := toTime.Weekday()
		if fw == 0 {
			fw = 7
		}
		if tw == 7 {
			tw = 7
		}
		if tw < fw {
			sub++
		}
	}
	if period == "month" {
		sub = int(toTime.Month()+12-fromTime.Month())%12 + 1
	}

	return sub
}

func CurrentTime(originTime time.Time, period string) time.Time {
	year, month, day := originTime.Date()
	if period == "hour" {
		return time.Date(year, month, day, originTime.Hour(), 0, 0, 0, time.Local)
	} else if period == "day" {
		return time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	} else if period == "week" {
		day = day - int(originTime.Weekday()) + 1
		return time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	} else if period == "month" {
		return time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	}
	return originTime
}

func TimestampToTime(from string) time.Time {
	if id := strings.Index(from, "."); id != -1 {
		from = from[:id]
	}
	i, _ := strconv.ParseInt(from, 10, 64)
	value := time.Unix(i, 0)
	return value
}
