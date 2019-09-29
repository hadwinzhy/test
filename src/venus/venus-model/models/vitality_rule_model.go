package models

import (
	"time"
)

// setting type can be interval or frequency
type VitalityRule struct {
	BaseModel
	AdminID           uint   `gorm:"index" json:"admin_id"`
	RuleType          string `gorm:"type:varchar(10)" json:"rule_type"`
	Admin             Admin
	VitalityRuleItems []VitalityRuleItem
}

type VitalityRuleItem struct {
	BaseModel
	VitalityRuleID uint `gorm:"index" json:"vitality_rule_id"`
	From           int  `gorm:"type:integer" json:"from"`
	To             int  `gorm:"type:integer" json:"to"`
}

type VitalityRuleItemSerializer struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	From      int       `json:"from"`
	To        int       `json:"to"`
}

type VitalityRuleSerializer struct {
	ID                uint                         `json:"id"`
	CreatedAt         time.Time                    `json:"created_at"`
	AdminID           uint                         `json:"admin_id"`
	RuleType          string                       `json:"rule_type"`
	VitalityRuleItems []VitalityRuleItemSerializer `json:"vitality_rule_items"`
}

func SortRules(ruleItems []VitalityRuleItem) {
	for i := 0; i < len(ruleItems); i++ {
		for j := i + 1; j < len(ruleItems); j++ {
			if ruleItems[i].From > ruleItems[j].From {
				ruleItems[i], ruleItems[j] = ruleItems[j], ruleItems[i]
			}
		}
	}
}
