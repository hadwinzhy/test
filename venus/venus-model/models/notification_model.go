package models

import (
	"encoding/json"
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
)

type Notification struct {
	BaseModel
	Name           string    `gorm:"type:varchar" json:"name"`
	Group          string    `gorm:"type:varchar" json:"group"`
	ShopName       string    `gorm:"type:varchar" json:"shop_name"`
	LastCapturedAt time.Time `gorm:"type:timestamp with time zone" json:"last_captured_at"`
	Gender         int       `gorm:"type:integer" json:"gender"` // 男 1 女 0
	Age            int       `gorm:"type:integer" json:"age"`
	Telephone      string    `gorm:"type:varchar(11)" json:"telephone"`
	Birthday       time.Time `gorm:"type:timestamp with time zone" json:"birthday"`
	Type           string    `gorm:"type:varchar" json:"type"`
	URL            string    `gorm:"type:varchar" json:"url"`
	Comments       string    `gorm:"type:varchar" json:"Comments"`
	CustomerID     uint      `json:"customer_id"`
	EventID        uint      `json:"event_id"`
	CompanyID      uint      `json:"company_id"`

	// 以下字段适用于异常预警
	WarningText string         `gorm:"type:text;column:warning_text" json:"warning_text"` // 提示语
	PayLoad     postgres.Jsonb `gorm:"type:jsonb";column:payload" json:"payload"`         // 透传给前端，使其方便调用参数进行页面跳转
	WarningType int            `gor:"type:integer" json:"warning_type"`
	Day         time.Time      `json:"day"`
}

type OnePayLoad struct {
	Name        int    `json:"name"`         // 跳转的类型: 7 个
	Param       string `json:"params"`       // 请求参数 key: 多个用`,`隔开
	Value       string `json:"value"`        //  请求参数 value: 多个用`,`隔开
	CompareTime int    `json:"compare_time"` // Y/W/M
}

func (n Notification) ToPayLoad() OnePayLoad {
	var onePayLoad OnePayLoad
	if err := json.Unmarshal(n.PayLoad.RawMessage, &onePayLoad); err != nil {
		return onePayLoad
	}
	return onePayLoad
}

func (n Notification) ToJsonb(one OnePayLoad) postgres.Jsonb {
	var value postgres.Jsonb
	value.RawMessage, _ = json.Marshal(&one)
	return value
}

type NotificationRelation struct {
	BaseModel
	AdminID        uint         `json:"admin_id"`
	Status         bool         `gorm:"type:boolean" json:"status"`
	NotificationID uint         `json:"notification_id"`
	Notification   Notification `json:"notification"`

	// 是否删除: 异常警告某账号删除信息
	IsDeleted bool `json:"is_deleted"`
}

type NotificationRelationSerializer struct {
	ID             uint         `json:"id"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
	DeletedAt      *time.Time   `json:"deleted_at"`
	AdminID        uint         `json:"admin_id"`
	Status         bool         `json:"status"`
	NotificationID uint         `json:"notification_id"`
	Notification   Notification `json:"notification"`
	IsDeleted      bool         `json:"is_deleted"`
}

type NotificationForWarningSerializer struct {
	ID          uint       `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	WarningText string     `json:"warning_text"`
	PayLoad     OnePayLoad `json:"payload"`
	WarningType int        `json:"warning_type"`
	Day         time.Time  `json:"day"`
}

func (n Notification) SerializerForWarning() NotificationForWarningSerializer {
	return NotificationForWarningSerializer{
		ID:          n.ID,
		CreatedAt:   n.CreatedAt.Round(time.Second),
		UpdatedAt:   n.UpdatedAt.Round(time.Second),
		WarningText: n.WarningText + "%",
		PayLoad:     n.ToPayLoad(),
		WarningType: n.WarningType,
		Day:         n.Day.Round(time.Second),
	}
}

func (n NotificationRelation) Serializer() *NotificationRelationSerializer {
	notification := n.Notification
	notification.CreatedAt = timeHandler(notification.CreatedAt)
	notification.UpdatedAt = timeHandler(notification.UpdatedAt)
	notification.Birthday = timeHandler(notification.Birthday)

	return &NotificationRelationSerializer{
		ID:             n.ID,
		CreatedAt:      timeHandler(n.CreatedAt),
		UpdatedAt:      timeHandler(n.UpdatedAt),
		DeletedAt:      n.DeletedAt,
		AdminID:        n.AdminID,
		Status:         n.Status,
		NotificationID: n.NotificationID,
		Notification:   n.Notification,
	}
}

// VIP 会员
type VIPNotificationDetail struct {
	ID            uint      `json:"id"`
	Name          string    `json:"name"`
	Group         string    `json:"group"`
	ShopName      string    `json:"shop_name"`
	CreatedAt     time.Time `json:"created_at"`
	LastCaptureAt time.Time `json:"last_capture_at"`
	Gender        int       `json:"gender"` // 男 1 女 0
	Age           int       `json:"age"`
	Telephone     string    `json:"telephone"`
	Birthday      time.Time `gorm:"type:timestamp with time zone" json:"birthday"`
	URL           string    `json:"url"`
	Comments      string    `json:"comments"`
	Type          string    `json:"type"`
	CustomerID    uint      `json:"customer_id"`
	EventID       uint      `json:"event_id"`
	CompanyID     uint      `json:"company_id"`
}

var MessageForVipOrCustomer = "%s %s %s, 性别：%s，年龄：%d，上次来访时间：%s"

var MessageForFrequentCustomer = "您有一名回头客进入%s，性别：%s"

func (n Notification) SerializerForVip() *VIPNotificationDetail {
	return &VIPNotificationDetail{
		ID:            n.ID,
		Name:          n.Name,
		Group:         n.Group,
		ShopName:      n.ShopName,
		CreatedAt:     timeHandler(n.CreatedAt),
		LastCaptureAt: timeHandler(n.LastCapturedAt),
		Gender:        n.Gender,
		Age:           n.Age,
		Telephone:     n.Telephone,
		Type:          n.Type,
		URL:           n.URL,
		Comments:      n.Comments,
		CustomerID:    n.CustomerID,
		EventID:       n.EventID,
		CompanyID:     n.CompanyID,
		Birthday:      timeHandler(n.Birthday),
	}
}

// 老客
type PotentialCustomerNotificationDetail struct {
	ID             uint      `json:"id"`
	Name           string    `json:"name"`
	Group          string    `json:"group"`
	ShopName       string    `json:"shop_name"`
	CreatedAt      time.Time `json:"created_at"`
	LastCapturedAt time.Time `json:"last_capture_at"`
	Gender         int       `json:"gender"` //男 1 女 0
	Age            int       `json:"age"`
	URL            string    `json:"url"`
	Type           string    `json:"type"`
	CustomerID     uint      `json:"customer_id"`
	EventID        uint      `json:"event_id"`
	CompanyID      uint      `json:"company_id"`
}

func (n Notification) SerializerForPotentialCustomer() *PotentialCustomerNotificationDetail {
	return &PotentialCustomerNotificationDetail{
		ID:             n.ID,
		Name:           n.Name,
		Group:          n.Group,
		ShopName:       n.ShopName,
		CreatedAt:      timeHandler(n.CreatedAt),
		LastCapturedAt: timeHandler(n.LastCapturedAt),
		Gender:         n.Gender,
		Age:            n.Age,
		Type:           n.Type,
		URL:            n.URL,
		CustomerID:     n.CustomerID,
		EventID:        n.EventID,
		CompanyID:      n.CompanyID,
	}
}

// 回头客
type FrequentCustomerNotificationDetail struct {
	ID             uint      `json:"id"`
	Group          string    `json:"group"`
	ShopName       string    `json:"shop_name"`
	CreatedAt      time.Time `json:"created_at"`
	LastCapturedAt time.Time `json:"last_capture_at"`
	Gender         int       `json:"gender"`
	Age            int       `json:"age"`
	URL            string    `json:"url"`
	Type           string    `json:"type"`
	CustomerID     uint      `json:"customer_id"`
	EventID        uint      `json:"event_id"`
	CompanyID      uint      `json:"company_id"`
}

func (n Notification) SerializerForFrequentCustomer() *FrequentCustomerNotificationDetail {
	return &FrequentCustomerNotificationDetail{
		ID:             n.ID,
		ShopName:       n.ShopName,
		Group:          n.Group,
		CreatedAt:      timeHandler(n.CreatedAt),
		LastCapturedAt: timeHandler(n.LastCapturedAt),
		Gender:         n.Gender,
		Age:            n.Age,
		Type:           n.Type,
		URL:            n.URL,
		CustomerID:     n.CustomerID,
		EventID:        n.EventID,
		CompanyID:      n.CompanyID,
	}
}

func timeHandler(value time.Time) time.Time {
	valueString := value.Format("2006-01-02 15:04:05")
	result, _ := time.ParseInLocation("2006-01-02 15:04:05", valueString, time.Local)
	return result
}
