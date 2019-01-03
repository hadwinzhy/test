package kafka

type producerParams struct {
	Topic  string `form:"topic" json:"topic"`
	Key    string `form:"key" json:"key"`
	Values string `form:"values" json:"values"`
}

//{
//	ID         uint `form:"id" json:"id"`
//	CustomerID uint `form:"customer_id" json:"customer_id"`
//	ShopID     uint `form:"shop_id" json:"shop_id"`
//	GroupID    uint `form:"group_id" json:"group_id"`
//	GroupUUID  uint `form:"group_uuid" json:"group_uuid"`
//}
