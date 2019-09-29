package models

import (
	"encoding/json"
	"fmt"

	"github.com/jinzhu/gorm/dialects/postgres"
)

type FrequentCustomerRule struct {
	BaseModel
	CompanyID     uint           `gorm:"index"`
	LowFrequency  postgres.Jsonb `gorm:"jsonb"`
	HighFrequency postgres.Jsonb `gorm:"jsonb"`
	Limit         int            `gorm:"type:integer"`
	readableRules ReadableFrequencyRule
}

type FrequentCustomerRules []FrequentCustomerRule

type ReadableFrequencyRule struct {
	LowFrequency  []rulePair `json:"low_frequency"`
	HighFrequency []rulePair `json:"high_frequency"`
	Limit         uint       `json:"limit"`
}

type rulePair struct {
	From uint   `json:"from"`
	To   uint   `json:"to"`
	Type string `json:"type"`
}

type rulePairs []rulePair

type FrequentCustomerRuleSerializer []rulePair

var lf = []rulePair{
	rulePair{
		From: 2,
		To:   3,
		Type: "low",
	},
}

var hf = []rulePair{
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
		To:   6,
		Type: "high",
	},
	rulePair{
		From: 7,
		To:   30,
		Type: "high",
	},
}

func (rule *FrequentCustomerRule) ReadableRule() ReadableFrequencyRule {
	if rule.readableRules.Limit != 0 {
		return rule.readableRules
	}

	// 从数据中解析出readableFrequentcyRule
	if rule.ID != 0 {
		var lowRulePair []rulePair
		var highRulePair []rulePair
		if err := json.Unmarshal(rule.LowFrequency.RawMessage, &lowRulePair); err != nil {
			fmt.Println("parse error", string(rule.LowFrequency.RawMessage), err)
			return ReadableFrequencyRule{
				LowFrequency:  lf,
				HighFrequency: hf,
				Limit:         3,
			}
		}

		if err := json.Unmarshal(rule.HighFrequency.RawMessage, &highRulePair); err != nil {
			fmt.Println("parse error", string(rule.HighFrequency.RawMessage), err)
			return ReadableFrequencyRule{
				LowFrequency:  lf,
				HighFrequency: hf,
				Limit:         3,
			}
		}

		limit := lowRulePair[len(lowRulePair)-1].To
		return ReadableFrequencyRule{
			LowFrequency:  lowRulePair,
			HighFrequency: highRulePair,
			Limit:         limit,
		}
	}

	return ReadableFrequencyRule{
		LowFrequency:  lf,
		HighFrequency: hf,
		Limit:         3,
	}

}

type FrequentCustomerRuleBasicSerializer struct {
	ID            uint       `json:"id"`
	CompanyID     uint       `json:"company_id"`
	LowFrequency  []rulePair `json:"low_frequency"`
	HighFrequency []rulePair `json:"high_frequency"`
	Limit         int        `json:"limit"`
}

func (rule FrequentCustomerRule) BasicSerializer() FrequentCustomerRuleBasicSerializer {
	return FrequentCustomerRuleBasicSerializer{
		ID:            rule.ID,
		CompanyID:     rule.CompanyID,
		LowFrequency:  rule.GetLowFrequency(),
		HighFrequency: rule.GetHighFrequency(),
		Limit:         rule.Limit,
	}
}

func (rule FrequentCustomerRule) GetLowFrequency() rulePairs {
	var lowRulePairs rulePairs
	if err := json.Unmarshal(rule.LowFrequency.RawMessage, &lowRulePairs); err != nil {
		return nil
	}
	return lowRulePairs
}

func (rule FrequentCustomerRule) GetHighFrequency() rulePairs {
	var highRulePairs rulePairs
	if err := json.Unmarshal(rule.HighFrequency.RawMessage, &highRulePairs); err != nil {
		return nil
	}
	return highRulePairs

}