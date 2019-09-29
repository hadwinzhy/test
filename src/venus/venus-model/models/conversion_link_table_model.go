package models

import (
	"encoding/json"

	"github.com/jinzhu/gorm/dialects/postgres"
)

type ConversionLinkTable struct {
	BaseModel
	CompanyID               uint           `gorm:"index" json:"company_id"`
	LinkTable               postgres.Jsonb `json:"link_table"`
	regionConversionRuleMap map[uint][]ConversionLink
}

type ConversionLink struct {
	RegionType       string `json:"region_type"` // from 或者 to
	ConversionRuleID uint   `json:"conversion_rule_id"`
}

type RegionLinkPair struct {
	RegionID        uint             `json:"region_id"`
	ConversionLinks []ConversionLink `json:"conversion_links"`
}

// 把pair转成map，根据region_id查对应的是哪个rule
func (table *ConversionLinkTable) GenerateMap() map[uint][]ConversionLink {
	if table.regionConversionRuleMap != nil {
		return table.regionConversionRuleMap
	}
	var pairs []RegionLinkPair
	err := json.Unmarshal(table.LinkTable.RawMessage, &pairs)
	if err != nil {
		return nil
	}

	ruleMap := make(map[uint][]ConversionLink)
	for _, pair := range pairs {
		ruleMap[pair.RegionID] = pair.ConversionLinks
	}
	table.regionConversionRuleMap = ruleMap

	return ruleMap
}

// 把map转成序列，存储为jsonb
func (table *ConversionLinkTable) SaveLinkTable(ruleMap map[uint][]ConversionLink) {
	pairs := make([]RegionLinkPair, 0)

	for regionID, links := range ruleMap {
		pair := RegionLinkPair{
			RegionID:        regionID,
			ConversionLinks: links,
		}

		pairs = append(pairs, pair)
	}

	table.LinkTable.RawMessage, _ = json.Marshal(pairs)
}
