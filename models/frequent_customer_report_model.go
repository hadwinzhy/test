package models

import (
	"reflect"
	"time"
)

type FrequentCustomerReport struct {
	BaseModel
	FrequentCustomerGroupID uint      `gorm:"index"`
	Date                    time.Time `gorm:"type:date"`
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

type FrequentCustomerHighTimeTables []FrequentCustomerHighTimeTable

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
	nowCount := reflect.ValueOf(table).Elem().FieldByName(phase).Int()
	nowCount++

	reflect.ValueOf(table).Elem().FieldByName(phase).SetInt(nowCount)
}

type FrequentCustomerHighTimeTableSerializer struct {
	ID         uint      `json:"id"`
	GroupID    uint      `json:"group_id"`
	Date       time.Time `json:"date"`
	PhaseOne   uint      `json:"phase_one"`
	PhaseTwo   uint      `json:"phase_two"`
	PhaseThree uint      `json:"phase_three"`
	PhaseFour  uint      `json:"phase_four"`
	PhaseFive  uint      `json:"phase_five"`
	PhaseSix   uint      `json:"phase_six"`
	PhaseSeven uint      `json:"phase_seven"`
	PhaseEight uint      `json:"phase_eight"`
}

func (table FrequentCustomerHighTimeTable) BasicSerializer() FrequentCustomerHighTimeTableSerializer {
	return FrequentCustomerHighTimeTableSerializer{
		ID:         table.ID,
		GroupID:    table.FrequentCustomerGroupID,
		Date:       table.Date,
		PhaseOne:   table.PhaseOne,
		PhaseTwo:   table.PhaseTwo,
		PhaseThree: table.PhaseThree,
		PhaseFour:  table.PhaseFour,
		PhaseFive:  table.PhaseFive,
		PhaseSix:   table.PhaseSix,
		PhaseSeven: table.PhaseSeven,
		PhaseEight: table.PhaseEight,
	}
}
