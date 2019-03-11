package frequent_person

type JudgeParams struct {
	EventID   uint   `form:"event_id" json:"event_id"`
	CompanyID uint   `form:"company_id" json:"company_id"`
	ShopID    uint   `form:"shop_id" json:"shop_id"`
	PersonID  string `form:"person_id" json:"person_id"`
	CreatedAt int64  `form:"created_at" json:"created_at"`
}
