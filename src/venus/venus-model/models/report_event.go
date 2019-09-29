package models

import (
	"encoding/json"
	"fmt"
	"time"
	"siren/venus/venus-model/models/logger"

	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
)

// AgeMap ..
type AgeMap map[string]int

// ReportEvent ...
type ReportEvent struct {
	ID                  uint           `gorm:"primary_key" json:"id"`
	Hour                time.Time      `gorm:"type:timestamp with time zone;index;" json:"hour"`
	MaleCount           uint           `gorm:"default:0" json:"male_count"`
	FemaleCount         uint           `gorm:"default:0" json:"female_count"`
	EventInCount        uint           `gorm:"default:0" json:"event_in_count"`
	CustomerInCount     uint           `gorm:"default:0" json:"customer_in_count"`
	VIPInCount          uint           `gorm:"default:0" json:"vip_in_count"`
	MaleVIPTimesCount   uint           `gorm:"default:0" json:"male_vip_times_count"`
	FemaleVIPTimesCount uint           `gorm:"default:0" json:"female_vip_times_count"`
	ShopID              uint           `gorm:"index;not null" json:"shop_id"`
	DeviceID            uint           `gorm:"index;not null" json:"device_id"`
	CompanyID           uint           `gorm:"index;not null" json:"company_id"`
	SmRegionID          uint           `gorm:"index" json:"sm_region_id"`
	AgeCount            postgres.Jsonb `gorm:"jsonb;" json:"age_count"`
	MaleAgeCount        postgres.Jsonb `gorm:"jsonb;" json:"male_age_count"`
	FemaleAgeCount      postgres.Jsonb `gorm:"jsonb;" json:"female_age_count"`
}

// ReportEventSerializer ...
type ReportEventSerializer struct {
	Time                time.Time `json:"time"`
	EventCount          uint      `json:"event_count"`
	CustomerCount       uint      `json:"customer_count"`
	MaleCount           uint      `json:"male_count"`
	FemaleCount         uint      `json:"female_count"`
	VIPCount            uint      `json:"vip_count"`
	MaleVIPTimesCount   uint      `json:"male_vip_times_count"`
	FemaleVIPTimesCount uint      `json:"female_vip_times_count"`
	AgeCount            AgeMap    `json:"age_count"`
	FemaleAgeCount      AgeMap    `json:"female_age_count"`
	MaleAgeCount        AgeMap    `json:"male_age_count"`
	AverageCount        uint      `json:"average_count"`
	FrequentCount       uint      `json:"frequent_count"`
}

// BaseSerialize ...
func (r *ReportEvent) BaseSerialize() ReportEventSerializer {
	return ReportEventSerializer{
		Time:                r.Hour.Local(),
		EventCount:          r.EventInCount,
		CustomerCount:       r.CustomerInCount,
		MaleCount:           r.MaleCount,
		FemaleCount:         r.FemaleCount,
		VIPCount:            r.VIPInCount,
		MaleVIPTimesCount:   r.MaleVIPTimesCount,
		FemaleVIPTimesCount: r.FemaleVIPTimesCount,
		AgeCount:            GetAgeStruct(r.AgeCount),
		FemaleAgeCount:      GetAgeStruct(r.FemaleAgeCount),
		MaleAgeCount:        GetAgeStruct(r.MaleAgeCount),
	}
}

// GetAgeStruct ...
func GetAgeStruct(ageString postgres.Jsonb) AgeMap {
	var ages AgeMap

	json.Unmarshal(ageString.RawMessage, &ages)
	return ages
}

// GetSubByPeriod ...
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
		if tw == 0 {
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

func getDurationByPeriod(period string) time.Duration {
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

func currentTime(originTime time.Time, period string) time.Time {
	year, month, day := originTime.Date()
	if period == "hour" {
		return time.Date(year, month, day, originTime.Hour(), 0, 0, 0, time.Local)
	} else if period == "day" {
		return time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	} else if period == "week" {
		weekday := originTime.Weekday()
		if weekday == 0 {
			weekday = 7
		}
		return time.Date(year, month, day, 0, 0, 0, 0, time.Local).AddDate(0, 0, 1-int(weekday))
	} else if period == "month" {
		return time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	}
	return originTime
}

// InsertMissing will insert blank for period
func InsertMissing(period string, fromTime time.Time, toTime time.Time, reports []ReportEvent, ageRerpots []ReportEvent) ([]ReportEventSerializer, int) {
	duration := getDurationByPeriod(period)
	var total int
	if period == "month" {
		total = (toTime.Year()-fromTime.Year())*12 + (int(toTime.Month()) - int(fromTime.Month())) + 1
	} else {
		total = GetSubByPeriod(period, fromTime, toTime, duration)
	}

	fmt.Println(total, duration)

	if total <= 0 {
		return []ReportEventSerializer{}, 0
	}
	result := make([]ReportEventSerializer, total, total+1)

	newTime := fromTime // init time
	for i := range result {
		result[i].Time = currentTime(newTime, period)
		for _, report := range reports {
			left := report.Hour.UTC().Unix()
			right := currentTime(newTime, period).Unix()
			if left == right {
				result[i] = report.BaseSerialize()
				for _, ageReport := range ageRerpots {
					if report.Hour.UTC().Unix() == ageReport.Hour.UTC().Unix() {
						result[i].AgeCount = GetAgeStruct(ageReport.AgeCount)
						result[i].FemaleAgeCount = GetAgeStruct(ageReport.FemaleAgeCount)
						result[i].MaleAgeCount = GetAgeStruct(ageReport.MaleAgeCount)
					}
				}
				break
			}
			if period == "year" {
				result[i] = report.BaseSerialize()
				for _, ageReport := range ageRerpots {
					if report.Hour.UTC().Unix() == ageReport.Hour.UTC().Unix() {
						result[i].AgeCount = GetAgeStruct(ageReport.AgeCount)
						result[i].FemaleAgeCount = GetAgeStruct(ageReport.FemaleAgeCount)
						result[i].MaleAgeCount = GetAgeStruct(ageReport.MaleAgeCount)
					}
				}
				break
			}
		}
		if period == "month" {
			month := int(newTime.Month())
			day := newTime.Day()
			nextYear := newTime.Year()
			if month == 12 {
				nextYear++
			}
			newTime = time.Date(nextYear, time.Month(month%12+1), day, 0, 0, 0, 0, time.Local)
		} else {
			newTime = newTime.Add(duration)
		}
	}
	return result, len(result)
}

func AgeCountAccumulate(targetJson *postgres.Jsonb, targetAge uint) {
	ages := GetAgeStruct(*targetJson)
	if ages == nil {
		ages = make(AgeMap, 1)
	}
	ages[fmt.Sprint(targetAge)]++
	ageStr, _ := json.Marshal(ages)
	if error := targetJson.Scan(ageStr); error != nil {
		logger.Error(nil, "grpc_event", "send", "fail in scan json.", error)
	}
}

func (eventReport *ReportEvent) UpdateReportCount(tx *gorm.DB, event *Event, hasCustomerGroup bool, isPotentialCustomer bool, isNewComer bool) {
	eventReport.EventInCount++

	// 性别信息变为非去重的
	if event.Gender == Male {
		eventReport.MaleCount++
	} else if event.Gender == Female {
		eventReport.FemaleCount++
	}

	// 年龄也是
	if event.Age != 0 {
		AgeCountAccumulate(&eventReport.AgeCount, event.Age)
		if event.Gender == Male {
			AgeCountAccumulate(&eventReport.MaleAgeCount, event.Age)
		} else {
			AgeCountAccumulate(&eventReport.FemaleAgeCount, event.Age)
		}
	}

	if hasCustomerGroup && !isPotentialCustomer {
		if event.Gender == Male {
			eventReport.MaleVIPTimesCount++
		} else if event.Gender == Female {
			eventReport.FemaleVIPTimesCount++
		}
	}

	if isNewComer { // new customer
		eventReport.CustomerInCount++

		if hasCustomerGroup && !isPotentialCustomer {
			eventReport.VIPInCount++
			// go notification.ShopVIPComeNotification(&event)
		}
	}

	if err := tx.Save(&eventReport).Error; err != nil {
		panic(err.Error())
	}
}
