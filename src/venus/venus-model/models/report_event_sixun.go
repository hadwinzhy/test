package models

import (
	"strconv"
	"time"
)

type AgePair struct {
	Age   int `json:"age"`
	Count int `json:"count"`
}
type AgeArray []AgePair

type ReportEventSixunSerializer struct {
	Time           time.Time `json:"time"`
	EventCount     uint      `json:"event_count"`
	CustomerCount  uint      `json:"customer_count"`
	MaleCount      uint      `json:"male_count"`
	FemaleCount    uint      `json:"female_count"`
	VIPCount       uint      `json:"vip_count"`
	AgeCount       AgeArray  `json:"age_count"`
	FemaleAgeCount AgeArray  `json:"female_age_count"`
	MaleAgeCount   AgeArray  `json:"male_age_count"`
}

func GetAgeArrayFromMap(ageCount AgeMap) (result AgeArray) {
	for age, count := range ageCount {
		ageNum, err := strconv.Atoi(age)
		if err != nil {
			continue
		}
		result = append(result, AgePair{
			Age:   ageNum,
			Count: count,
		})
	}

	return
}

func ConvertToSixunSerialize(r ReportEventSerializer) ReportEventSixunSerializer {
	return ReportEventSixunSerializer{
		Time:           r.Time,
		EventCount:     r.EventCount,
		CustomerCount:  r.CustomerCount,
		MaleCount:      r.MaleCount,
		FemaleCount:    r.FemaleCount,
		VIPCount:       r.VIPCount,
		AgeCount:       GetAgeArrayFromMap(r.AgeCount),
		FemaleAgeCount: GetAgeArrayFromMap(r.FemaleAgeCount),
		MaleAgeCount:   GetAgeArrayFromMap(r.MaleAgeCount),
	}
}
