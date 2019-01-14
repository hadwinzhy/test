package models

import (
	"reflect"
	"siren/pkg/database"
	"siren/pkg/utils"
	"time"
)

type FrequentCustomerReport struct {
	BaseModel
	FrequentCustomerGroupID uint      `gorm:"index"`
	Date                    time.Time `gorm:"type:date"`
	Hour                    time.Time `gorm:"type:timestamp with time zone"`
	HighFrequency           uint      `gorm:"type:integer"`
	LowFrequency            uint      `gorm:"type:integer"`
	NewComer                uint      `gorm:"type:integer"`
	SumInterval             uint      `gorm:"type:integer"`
	SumTimes                uint      `gorm:"type:integer"`
}

// 总人数，高频次数，低频次数，新客数，总到访间隔天数，总到访天数

type FrequentCustomerHighTimeTable struct {
	BaseModel
	FrequentCustomerGroupID uint      `gorm:"index"`
	Date                    time.Time `gorm:"type:date"`
	PhaseOne                uint      `gorm:"type:integer"`
	PhaseTwo                uint      `gorm:"type:integer"`
	PhaseThree              uint      `gorm:"type:integer"`
	PhaseFour               uint      `gorm:"type:integer"`
	PhaseFive               uint      `gorm:"type:integer"`
	PhaseSix                uint      `gorm:"type:integer"`
	PhaseSeven              uint      `gorm:"type:integer"`
	PhaseEight              uint      `gorm:"type:integer"`
}

var hourPhaseMap = map[int]string{
	8:  "PhaseOne",
	9:  "PhaseOne",
	10: "PhaseTwo",
	11: "PhaseTwo",
	12: "PhaseThree",
	13: "PhaseThree",
	14: "PhaseFour",
	15: "PhaseFour",
	16: "PhaseFive",
	17: "PhaseFive",
	18: "PhaseSix",
	19: "PhaseSix",
	20: "PhaseSeven",
	21: "PhaseSeven",
	22: "PhaseEight",
	23: "PhaseEight",
}

func TimeToPhase(captureAt time.Time) string {
	return hourPhaseMap[captureAt.Hour()]
}

var phaseTitleMap = map[string]string{
	"PhaseOne":   "8:00-10:00",
	"PhaseTwo":   "10:00-12:00",
	"PhaseThree": "12:00-14:00",
	"PhaseFour":  "14:00-16:00",
	"PhaseFive":  "16:00-18:00",
	"PhaseSix":   "18:00-20:00",
	"PhaseSeven": "20:00-22:00",
	"PhaseEight": "22:00-24:00",
}

func (table *FrequentCustomerHighTimeTable) AddCount(captureAt time.Time) {
	phase := TimeToPhase(captureAt)
	if phase == "" {
		return
	}
	nowCount := reflect.ValueOf(table).Elem().FieldByName(phase).Uint()
	nowCount++

	reflect.ValueOf(table).Elem().FieldByName(phase).SetUint(nowCount)

	database.POSTGRES.Save(table)
}

type FrequentCustomerReports []FrequentCustomerReport

func (reports FrequentCustomerReports) InsertMissing(period string, fromTime time.Time, toTime time.Time, sortBy string) ([]FrequentCustomerReport, error) {
	duration := utils.GetDurationByPeriod(period)
	var total int
	if period == "month" {
		total = (toTime.Year()-fromTime.Year())*12 + (int(toTime.Month()) - int(fromTime.Month())) + 1
	} else {
		total = utils.GetSubByPeriod(period, fromTime, toTime, duration)
	}
	if total <= 0 {
		return []FrequentCustomerReport{}, nil
	}
	result := make([]FrequentCustomerReport, total, total+1)

	newTime := fromTime // init time
	for i := range result {
		result[i].Hour = utils.CurrentTime(newTime, period)
		for _, report := range reports {
			left := utils.CurrentDate(report.Hour).UTC().Unix()
			right := utils.CurrentTime(newTime, period).Unix()
			if left == right {
				result[i] = report
				break
			}
			if period == "year" {
				result[i] = report
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

	if sortBy == "desc" {
		anotherResult := make([]FrequentCustomerReport, len(result))
		for i := range result {
			anotherResult[i] = result[len(result)-1-i]
		}
		result = anotherResult
	}
	return result, nil
}
