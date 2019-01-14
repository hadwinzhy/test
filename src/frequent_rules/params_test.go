package frequent_rules

import (
	"fmt"
	"testing"
)

func TestPostRulesIsSuitable(test *testing.T) {
	tt := []struct {
		params PostRules
	}{
		{
			params: PostRules{
				LowRule: []OneRule{
					{
						From: 2,
						To:   3,
					},
				},
				HighRule: []OneRule{
					{
						From: 3,
						To:   30,
					},
				},
			},
		},
	}
	for _, t := range tt {
		fmt.Println(t.params.AllCount())
		fmt.Println(t.params.InclusiveRange())
	}
}
