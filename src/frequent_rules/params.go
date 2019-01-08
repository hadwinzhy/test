package frequent_rules

import (
	"encoding/json"

	"qiniupkg.com/x/log.v7"

	"github.com/jinzhu/gorm/dialects/postgres"
)

type PostRules struct {
	CompanyID uint      `form:"company_id" json:"company_id"`
	LowRule   []OneRule `form:"low_rule" json:"low_rule"`
	HighRule  []OneRule `form:"high_rule" json:"high_rule"`
	Limit     int       `form:"limit" json:"limit"`
}

type OneRule struct {
	From int    `form:"from" json:"from"`
	To   int    `form:"to" json:"to"`
	Type string `form:"type" json:"type"`
}

// to >= from
func (one OneRule) IsSuitable() bool {
	return one.To >= one.From
}

// lowRule and highRule should be suit
func (rule PostRules) IsSuitableParam() bool {
	if len(rule.LowRule) == 0 || len(rule.HighRule) == 0 {
		log.Println("error found step zero")
		return false
	}
	if rule.LowRule[len(rule.LowRule)-1].To > rule.Limit {
		log.Println("error found step one")
		return false
	}

	if rule.HighRule[0].From < rule.Limit {
		log.Println("error found step two")
		return false
	}

	for _, i := range rule.LowRule {
		if !i.IsSuitable() {
			log.Println("error found step three")
			return false
		}
	}
	if rule.LowRule[len(rule.LowRule)-1].To > rule.HighRule[0].From {
		log.Println("error found step four")
		return false
	}

	for _, j := range rule.HighRule {
		if !j.IsSuitable() {
			log.Println("error found step five")
			return false
		}
	}

	if len(rule.LowRule)+len(rule.HighRule) > 5 {
		log.Println("error found step six")
		return false
	}
	return true
}

func (rule PostRules) JsonbLowHandler() postgres.Jsonb {
	var values postgres.Jsonb
	values.RawMessage, _ = json.Marshal(rule.LowRule)
	return values
}

func (rule PostRules) JsonbHighHandler() postgres.Jsonb {
	var values postgres.Jsonb
	values.RawMessage, _ = json.Marshal(rule.HighRule)
	return values
}
