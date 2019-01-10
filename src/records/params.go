package records

import (
	"siren/pkg/controllers"
	"time"
)

type FrequentCustomerRecordParams struct {
	CompanyID uint `form:"company_id" binding:"required"`
	ShopID    uint `form:"shop_id"`
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
}

func (form *FrequentCustomerRecordParams) Normalize() {
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
