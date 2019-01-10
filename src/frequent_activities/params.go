package frequent_activities

type GetFrequentActivitiesParams struct {
	CompanyID uint   `form:"company_id" json:"company_id"`
	ShopID    uint   `form:"shop_id" json:"shop_id"`
	Date      string `form:"date" json:"date"`
	Type      string `form:"type" json:"type" binding:"required,eq=store|eq=mall"` // 用于区分是门店还是购物中心
}
