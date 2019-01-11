package distributions

import "siren/pkg/controllers"

type ListDistributionParams struct {
	ReturnALL string `form:"return" binding:"eq=all_list|eq=all_count|eq=all_list_average_count"`
	CompanyID uint   `form:"company_id" binding:"required"`
	ShopID    uint   `form:"shop_id"`
	controllers.FromToParam
}

type valueProportionPair struct {
	Count      uint   `json:"count"`
	Proportion string `json:"proportion"`
}

type DistributionOutput struct {
	Date             string              `json:"date"`
	High             valueProportionPair `json:"high"`
	Low              valueProportionPair `json:"low"`
	New              valueProportionPair `json:"new"`
	AverageInterval  string              `json:"average_interval"`
	AverageFrequency string              `json:"average_frequency"`
}
