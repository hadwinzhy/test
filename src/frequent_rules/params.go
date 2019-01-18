package frequent_rules

import (
	"encoding/json"
	"log"
	"siren/pkg/controllers/errors"

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

// count
func (one OneRule) Count() int {
	return one.To - one.From + 1
}

// postRules count should be  equal 30
func (rule PostRules) AllCount() bool {
	var count int
	for index, i := range rule.LowRule {
		if index == 0 {
			if i.From < 2 {
				return false
			}
			count += i.To
		} else {
			count += i.Count()
		}
	}
	for _, j := range rule.HighRule {
		count += j.Count()
	}
	log.Println("count", count)
	if count == 30 {
		return true
	}
	return false

}

// postRules
func (rule PostRules) InclusiveRange() bool {
	var numbers []int
	for _, i := range rule.LowRule {
		numbers = append(numbers, i.From)
		numbers = append(numbers, i.To)
	}
	for _, j := range rule.HighRule {
		numbers = append(numbers, j.From)
		numbers = append(numbers, j.To)
	}

	var numberCount = make(map[int]int)
	for _, k := range numbers {
		if numberCount[k] != 0 {
			numberCount[k]++
		} else {
			numberCount[k] = 1
		}
		log.Println("map", numberCount)
		if numberCount[k] > 2 {
			return false
		}
	}
	return true
}

// lowRule and highRule should be suit
func (rule PostRules) IsSuitableParam() (bool, *errors.Error) {
	if !rule.InclusiveRange() {
		return false, &errors.Error{
			ErrorCode: errors.ErrorCode{
				HTTPStatus: 400,
				Code:       400,
				Title:      "is not inclusive range",
				TitleZH:    "高低频中存在不是闭区间的集合",
			},
			Detail: "高低频规则中存在不是闭区间的集合",
		}
	}
	if !rule.AllCount() {
		return false, &errors.Error{
			ErrorCode: errors.ErrorCode{
				HTTPStatus: 400,
				Code:       400,
				Title:      "all count of rules should be equal 30",
				TitleZH:    "高低频规则总天数不等于30、低频最小从2开始或者区间之间不得有重复时间",
			},
			Detail: "高低频规则总天数需满30天, 或者低频最小从2开始",
		}
	}
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
