package distributions

import (
	"siren/pkg/controllers"
	"siren/pkg/utils"
)

type ListDistributionParams struct {
	ReturnALL string `form:"return" binding:"eq=all_list|eq=all_count|eq=all_list_average_count"`
	CompanyID uint   `form:"company_id" binding:"required"`
	ShopID    string `form:"shop_id"`
	SortBy    string `form:"sort_by"`
	controllers.FromToParam
	controllers.PeriodParam
	shopIDs []uint
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

func (form *ListDistributionParams) Normalize() {
	form.ShopID, form.shopIDs = utils.NumberGroupStringNormalize(form.ShopID)

	form.FromToParam.Normalize()
	if form.Period == "" {
		form.Period = "day"
	}

	if form.SortBy != "desc" {
		form.SortBy = "asc"
	}
}
