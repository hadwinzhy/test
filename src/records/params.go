package records

import (
	"siren/pkg/controllers"
	"siren/pkg/utils"
	"time"
)

type FrequentCustomerRecordParams struct {
	CompanyID uint   `form:"company_id" binding:"required"`
	ShopID    string `form:"shop_id"`
	shopIDs   []uint
	controllers.PaginationParam
	controllers.FromToParam
}

type FrequentCustomerRecord struct {
	FrequentCustomerPersonID uint      `json:"frequent_customer_person_id"`
	FirstCaptureURL          string    `json:"first_capture_url"`
	Name                     string    `json:"name"`
	CaptureAt                time.Time `json:"capture_at"`
	Age                      uint      `json:"age"`
	Gender                   uint      `json:"gender"`
	ShopID                   uint      `json:"shop_id"`
	ShopName                 string    `json:"shop_name"`
	DeviceID                 uint      `json:"device_id"`
	DeviceName               string    `json:"device_name"`
	Frequency                uint      `json:"frequency"`
	Note                     string    `json:"note"`
	personID                 string
}

func (form *FrequentCustomerRecordParams) Normalize() {
	form.ShopID, form.shopIDs = utils.NumberGroupStringNormalize(form.ShopID)

	if form.Page == 0 {
		form.Page = 1
	}

	if form.PerPage == 0 {
		form.PerPage = 5
	}

	if form.OrderBy == "capture_at" {
		form.OrderBy = "hour" // 回头客用的字段是hour
	}

	form.FromToParam.Normalize()
}

// 需要传company shop params 保证权限
type CompanyShopParams struct {
	CompanyID uint `form:"company_id" json:"company_id"`
	ShopID    uint `form:"shop_id" json:"shop_id"`
}

type FrequentCustomerRecordDetailParams struct {
	controllers.PaginationParam
	CompanyShopParams
}

type SingleEventRecord struct {
	ID              uint      `json:"id"`
	OriginalFaceURL string    `json:"original_face_url"`
	CaptureAt       time.Time `json:"capture_at"`
	ShopID          uint      `json:"shop_id"`
	ShopName        string    `json:"shop_name"`
	DeviceID        uint      `json:"device_id"`
	DeviceName      string    `json:"device_name"`
}

type FrequentCustomerRecordMarkParams struct {
	Name string `json:"name"`
	Note string `json:"note"`
	CompanyShopParams
}
