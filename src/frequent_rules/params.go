package frequent_rules

import (
	"encoding/json"
	"siren/pkg/controllers/errors"

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
func (rule PostRules) IsSuitableParam() (bool, *errors.Error) {
	if len(rule.LowRule) == 0 || len(rule.HighRule) == 0 {
		log.Println("error found step zero")
		return false, &errors.Error{
			ErrorCode: errors.ErrorCode{
				HTTPStatus: 400,
				Code:       400,
				Title:      "the length of lowRules or highRules should be more than one",
				TitleZH:    "高低频规则至少需要设置一个",
			},
			Detail: "高低频规则至少需要设置一个",
		}
	}
	if rule.LowRule[len(rule.LowRule)-1].To > rule.Limit {
		log.Println("error found step one")
		return false, &errors.Error{
			ErrorCode: errors.ErrorCode{
				HTTPStatus: 400,
				Code:       400,
				Title:      "the low rule  in (to) is larger than rule limit",
				TitleZH:    "低频规则超过阈值",
			},
			Detail: "低频规则超过阈值",
		}
	}

	if rule.HighRule[0].From < rule.Limit {
		log.Println("error found step two")
		return false, &errors.Error{
			ErrorCode: errors.ErrorCode{
				HTTPStatus: 400,
				Code:       400,
				Title:      "the high rule in (from) is less than rule limit",
				TitleZH:    "高频规则小于阈值",
			},
			Detail: "高频规则小于阈值",
		}
	}

	for _, i := range rule.LowRule {
		if !i.IsSuitable() {
			log.Println("error found step three")
			return false, &errors.Error{
				ErrorCode: errors.ErrorCode{
					HTTPStatus: 400,
					Code:       400,
					Title:      "the low rule from larger than to",
					TitleZH:    "低频规则中From大于To",
				},
				Detail: "低频规则中From大于To",
			}
		}
	}
	if rule.LowRule[len(rule.LowRule)-1].To > rule.HighRule[0].From {
		log.Println("error found step four")
		return false, &errors.Error{
			ErrorCode: errors.ErrorCode{
				HTTPStatus: 400,
				Code:       400,
				Title:      "rule is not correct",
				TitleZH:    "低频规则的最大值大于高频规则的最小值",
			},
			Detail: "低频规则的最大值大于高频规则的最小值",
		}
	}

	for _, j := range rule.HighRule {
		if !j.IsSuitable() {
			log.Println("error found step five")
			return false, &errors.Error{
				ErrorCode: errors.ErrorCode{
					HTTPStatus: 400,
					Code:       400,
					Title:      "the high rule from larger than to",
					TitleZH:    "高频规则中From大于To",
				},
				Detail: "高频规则中From大于To",
			}
		}
	}

	if len(rule.LowRule)+len(rule.HighRule) > 5 {
		log.Println("error found step six")
		return false, &errors.Error{
			ErrorCode: errors.ErrorCode{
				HTTPStatus: 400,
				Code:       400,
				Title:      "the length of rules large than 5",
				TitleZH:    "总规则数大于5",
			},
			Detail: "总规则数大于5",
		}
	}
	return true, nil

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
