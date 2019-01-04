package models

import (
	"github.com/jinzhu/gorm/dialects/postgres"
)

type FrequentCustomerRule struct {
	BaseModel
	CompanyID     uint           `gorm:"index"`
	LowFrequency  postgres.Jsonb `gorm:"jsonb"`
	HighFrequency postgres.Jsonb `gorm:"jsonb"`
	readableRules ReadableFrequencyRule
}
type ReadableFrequencyRule struct {
	LowFrequency  []rulePair
	HighFrequency []rulePair
	Limit         uint
}

type rulePair struct {
	From uint   `json:"from"`
	To   uint   `json:"to"`
	Type string `json:"type"`
}

type FrequentCustomerRuleSerializer []rulePair

func (rule *FrequentCustomerRule) ReadableRule() ReadableFrequencyRule {
	if rule.readableRules.Limit != 0 {
		return rule.readableRules
	}

	// if rule.ID == 0 {
	// TODO: 从数据中解析出readableFrequencyRule的方法
	lf := []rulePair{
		rulePair{
			From: 1,
			To:   2,
			Type: "low",
		},
	}

	hf := []rulePair{
		rulePair{
			From: 3,
			To:   3,
			Type: "high",
		},
		rulePair{
			From: 4,
			To:   4,
			Type: "high",
		},
		rulePair{
			From: 5,
			To:   5,
			Type: "high",
		},
		rulePair{
			From: 6,
			To:   30,
			Type: "high",
		},
	}
	return ReadableFrequencyRule{
		LowFrequency:  lf,
		HighFrequency: hf,
		Limit:         3,
	}
	// }
}
