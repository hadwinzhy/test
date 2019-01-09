package frequent_table

type GetFrequentTableParams struct {
	CompanyID uint   `form:"company_id" json:"company_id"`
	ShopID    uint   `form:"shop_id" json:"shop_id"`
	Date      string `form:"date" json:"date"`
}
