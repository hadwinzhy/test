package models

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"siren/venus/venus-model/models/connectors"
	"siren/venus/venus-model/models/logger"

	raven "github.com/getsentry/raven-go"
	"github.com/jinzhu/gorm"

	"github.com/lib/pq"
)

// Customer ...
type Customer struct {
	BaseModel
	CompanyID           uint           `gorm:"index;not null" json:"company_id"`
	PersonID            string         `gorm:"type:varchar(100);index;not null" json:"person_id"`
	Name                string         `gorm:"type:varchar(50)" json:"name"`
	Phone               string         `gorm:"type:varchar(50)" json:"phone"`
	CaptureAt           time.Time      `gorm:"type:timestamp with time zone" json:"capture_at"`
	LastCaptureAt       time.Time      `gorm:"type:timestamp with time zone" json:"last_capture_at"`
	Avatars             pq.StringArray `gorm:"type:varchar(255)[]" json:"avatars"`
	OriginalFaceUrl     string         `gorm:"type:varchar(255)" json:"original_face_url"`
	LastEventDeviceName string         `gorm:"type:varchar(100)" json:"last_event_device_name"`
	LastEventShopName   string         `gorm:"type:varchar(100)" json:"last_event_shop_name"` // 在shopping mall逻辑中取的是region
	Status              string         `gorm:"type:varchar(100);index" json:"status"`
	EventsCount         uint           `json:"events_count"`
	CustomerGroupID     uint           `gorm:"index" json:"customer_group_id"`
	Age                 uint           `json:"age"`
	Gender              uint           `json:"gender"`
	Birthday            time.Time      `gorm:"type:date" json:"birthday"`
	Note                string         `json:"note"`
	Company             Company
	CustomerGroup       CustomerGroup
}

// CustomerBasicSerializer ...
type CustomerBasicSerializer struct {
	ID                  uint                         `json:"id"`
	CreatedAt           time.Time                    `json:"created_at"`
	UpdatedAt           time.Time                    `json:"updated_at"`
	CaptureAt           time.Time                    `json:"capture_at"`
	LastCaptureAt       time.Time                    `json:"last_capture_at"`
	Name                string                       `json:"name"`
	Phone               string                       `json:"phone"`
	Avatars             []ImageSerializer            `json:"avatars"`
	OriginalFace        ImageSerializer              `json:"original_face"`
	OriginalFaceURL     string                       `json:"original_face_url"`
	PersonID            string                       `json:"person_id"`
	LastEventShopName   string                       `json:"last_event_shop_name"`
	LastEventDeviceName string                       `json:"last_event_device_name"`
	Status              string                       `json:"status"`
	CustomerGroupID     uint                         `json:"customer_group_id"`
	CustomerGroup       CustomerGroupBasicSerializer `json:"customer_group"`
	EventsCount         uint                         `json:"events_count"`
	Age                 uint                         `json:"age"`
	Gender              uint                         `json:"gender"`
	Birthday            time.Time                    `json:"birthday"`
	Note                string                       `json:"note"`
}

type CustomerEventSerializer struct {
	CustomerBasicSerializer
	Event EventBasicSerializer `json:"event_info"`
}

type CustomerEventMallSerializer struct {
	CustomerBasicSerializer
	Event EventMallSerializer `json:"event_info"`
}

// BasicSerializer ...
func (c *Customer) BasicSerializer(customerGroup *CustomerGroup) CustomerBasicSerializer {
	avatars := make([]ImageSerializer, len(c.Avatars))
	for i, avatarURL := range c.Avatars {
		avatars[i] = ImageSerializer{
			URL: avatarURL,
		}
	}

	return CustomerBasicSerializer{
		ID:                  c.ID,
		CreatedAt:           c.CreatedAt.Round(time.Second),
		UpdatedAt:           c.UpdatedAt.Round(time.Second),
		CaptureAt:           c.CaptureAt.Round(time.Second),
		LastCaptureAt:       c.LastCaptureAt.Round(time.Second),
		Name:                c.Name,
		Phone:               c.Phone,
		Avatars:             avatars,
		PersonID:            c.PersonID,
		LastEventShopName:   c.LastEventShopName,
		LastEventDeviceName: c.LastEventDeviceName,
		Status:              c.Status,
		CustomerGroupID:     c.CustomerGroupID,
		CustomerGroup:       c.CustomerGroup.BasicSerialize(),
		EventsCount:         c.EventsCount,
		OriginalFaceURL:     c.OriginalFaceUrl,
		OriginalFace: ImageSerializer{
			URL: c.OriginalFaceUrl,
		},
		Age:      c.Age,
		Gender:   c.Gender,
		Birthday: c.Birthday,
		Note:     c.Note,
	}
}

// CustomerEventSerialize ...
func (c *Customer) CustomerEventSerialize(eventTime *time.Time, lastEventTime *time.Time, eventDetail *Event) CustomerEventSerializer {
	var result CustomerEventSerializer
	customerResult := c.BasicSerializer(&c.CustomerGroup)

	eventResult := eventDetail.BasicSerialize()

	// round time
	if eventTime != nil {
		customerResult.CaptureAt = *eventTime
	}
	if lastEventTime != nil {
		customerResult.LastCaptureAt = *lastEventTime
	}

	customerResult.CaptureAt = customerResult.CaptureAt.Round(time.Second)
	customerResult.LastCaptureAt = customerResult.LastCaptureAt.Round(time.Second)
	customerResult.OriginalFace = eventResult.OriginalFace
	customerResult.OriginalFaceURL = eventResult.OriginalFaceURL

	result.CustomerBasicSerializer = customerResult

	// event_info
	result.Event = eventResult

	return result
}

func (c *Customer) CustomerEventMallSerialize(eventTime *time.Time, lastEventTime *time.Time, eventDetail *Event) CustomerEventMallSerializer {
	var result CustomerEventMallSerializer
	customerResult := c.BasicSerializer(&c.CustomerGroup)

	eventResult := eventDetail.MallSerialize()

	// round time
	if eventTime != nil {
		customerResult.CaptureAt = *eventTime
	}
	if lastEventTime != nil {
		customerResult.LastCaptureAt = *lastEventTime
	}

	customerResult.CaptureAt = customerResult.CaptureAt.Round(time.Second)
	customerResult.LastCaptureAt = customerResult.LastCaptureAt.Round(time.Second)
	customerResult.OriginalFace = eventResult.OriginalFace
	customerResult.OriginalFaceURL = eventResult.OriginalFaceURL

	result.CustomerBasicSerializer = customerResult

	// event_info
	result.Event = eventResult

	return result
}

func (c *Customer) CustomerEventSerializeWithVipRecordCount(eventTime *time.Time, lastEventTime *time.Time, eventDetail *Event, count uint) CustomerEventSerializer {
	result := c.CustomerEventSerialize(eventTime, lastEventTime, eventDetail)
	// when output, modify its event counts
	result.EventsCount = count
	result.Event.CustomerEventsCount = count
	return result
}

func (c *Customer) CustomerEventMallSerializeWithVipRecordCount(eventTime *time.Time, lastEventTime *time.Time, eventDetail *Event, count uint) CustomerEventMallSerializer {
	result := c.CustomerEventMallSerialize(eventTime, lastEventTime, eventDetail)
	// when output, modify its event counts
	result.EventsCount = count
	result.Event.CustomerEventsCount = count
	return result
}

func (c *Customer) getGroup(tx *gorm.DB, group *CustomerGroup) (err error) {
	tx.Model(c).Related(&group)

	if group.ID == 0 { // some person doesn't belong to any customer group
		return
	}

	return
}

type regPersonReq struct {
	Name      string   `json:"name"`
	Avatars   []string `json:"avatars"`
	GroupUUID string   `json:"group_uuid"`
}

type regPersonResp struct {
	Status    string `json:"status"`
	PersonID  string `json:"person_id"`
	Name      string `json:"name"`
	FaceCount int    `json:"face_count"`
}

type addFaceReq struct {
	PersonUUID string   `json:"person_uuid"`
	FaceUUIDs  []string `json:"face_uuids"`
}

type addFaceResp struct {
	Status    string `json:"status"`
	PersonID  string `json:"person_id"`
	Name      string `json:"name"`
	FaceCount int    `json:"face_count"`
}

// AddFaces will patch face of customer
func (c *Customer) AddFaces(faceIDs []string) (err error) {
	req := addFaceReq{
		PersonUUID: c.PersonID,
		FaceUUIDs:  faceIDs,
	}

	var reqJSON []byte

	if reqJSON, err = json.Marshal(req); err != nil {
		return
	}

	// call pandora to create bindings
	reqReader := bytes.NewBuffer(reqJSON)
	var response http.Response

	response, err = connectors.HTTPRequest("POST", "/v1/api/people/add_face", reqReader)

	defer response.Body.Close()

	var bodyContent []byte
	var x addFaceResp
	bodyContent, _ = ioutil.ReadAll(response.Body)

	if err != nil || response.StatusCode != 200 {
		if err == nil {
			err = errors.New("add face failed")
		}
		logger.Error(nil, "customer", "add_face", "ADD_PHOTO_IN_PANDORA_ERROR, customer_id = ", c.ID, " face_ids: ", faceIDs, " error: ", err)
		raven.CaptureError(err, map[string]string{
			"action":   "ADD_PHOTO_IN_PANDORA_ERROR",
			"customer": strconv.Itoa(int(c.ID)),
			"face_ids": fmt.Sprintf("%v", faceIDs),
			"resp":     string(bodyContent),
		})
		return
	}

	json.Unmarshal(bodyContent, &x)

	logger.Info(nil, "customer", "add face", req, x)

	return
}

func common(c *Customer, request interface{}, url string) (err error) {
	var reqJSON []byte

	if reqJSON, err = json.Marshal(request); err != nil {
		return
	}

	// call pandora to create bindings
	reqReader := bytes.NewBuffer(reqJSON)
	var response http.Response
	fmt.Println(reqReader)
	response, err = connectors.HTTPRequest("POST", url, reqReader)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	var bodyContent []byte
	var x regPersonResp
	bodyContent, _ = ioutil.ReadAll(response.Body)
	fmt.Println(string(bodyContent), "888")
	if strings.Contains(string(bodyContent), "Code") {
		if strings.Contains(string(bodyContent), "401") {
			err = errors.New("注册图片未检测到人脸，请重新上传")
			return
		} else if strings.Contains(string(bodyContent), "402") {
			err = errors.New("注册图片大小不符合规范，请重新上传")
			return
		} else if strings.Contains(string(bodyContent), "403") {
			err = errors.New("注册图片质量不符合规范，请重新上传")
			return
		}
	}

	json.Unmarshal(bodyContent, &x)
	c.PersonID = x.PersonID

	if err == nil && x.PersonID == "" {
		err = errors.New("no person")
	}

	if err != nil || response.StatusCode != 200 {
		if err == nil {
			err = errors.New("register person fail")
		}
		logger.Error(nil, "customer", "create", "CREATE_PEOPLE_IN_PANDORA_ERROR, customer_id = ", c.ID, " error: ", err)
		raven.CaptureError(err, map[string]string{"action": "CREATE_PEOPLE_IN_PANDORA_ERROR", "customer": strconv.Itoa(int(c.ID))})
		return
	}
	return
}

// RegisterPerson ...
func (c *Customer) RegisterPerson(group *CustomerGroup) (err error) {
	req := regPersonReq{
		Name:      c.Name,
		Avatars:   c.Avatars,
		GroupUUID: group.GroupID,
	}

	if group.GroupType == CUSTOMER_GROUP_TYPE_POTENTIAL && req.Name == "" {
		req.Name = "某老客"
	}

	err = common(c, req, "/v1/api/people")
	return

}

// PatchRegisterPerson
func (c *Customer) PatchRegisterPerson(group *CustomerGroup, remove []string, add []string) (err error) {
	request := struct {
		PersonID  string   `json:"people_id"`
		GroupUUID string   `json:"group_id"`
		RemoveURL []string `json:"remove_url"`
		AddURL    []string `json:"add_url"`
	}{
		PersonID:  c.PersonID,
		GroupUUID: group.GroupID,
		RemoveURL: remove,
		AddURL:    add,
	}
	if group.GroupType == CUSTOMER_GROUP_TYPE_POTENTIAL && c.Name == "" {
		c.Name = "某老客"
	}

	err = common(c, request, "/v1/api/people/register")
	return err

}

type unregPersonReq struct {
	PersonUUIDs []string `json:"person_uuids"`
	GroupUUID   string   `json:"group_uuid"`
}

// UnregisterPerson ...
func (c *Customer) UnregisterPerson(tx *gorm.DB) (err error) {
	if c.CustomerGroupID > 0 {
		var customerGroup CustomerGroup
		if err = tx.Model(c).Related(&customerGroup).Error; err != nil {
			return
		}

		req := unregPersonReq{
			PersonUUIDs: []string{c.PersonID},
			GroupUUID:   customerGroup.GroupID,
		}

		var reqJSON []byte

		if reqJSON, err = json.Marshal(req); err != nil {
			return
		}

		// call pandora to removebindings
		reqReader := bytes.NewBuffer(reqJSON)
		var response http.Response

		response, err = connectors.HTTPRequest("POST", "/v1/api/groups/remove_person", reqReader)
		if err != nil {
			logger.Error(nil, "customer", "delete", "DELETE_PEOPLE_IN_PANDORA_ERROR, customer_id = ", c.ID, " error: ", err)
			raven.CaptureError(err, map[string]string{"action": "DELETE_PEOPLE_IN_PANDORA_ERROR", "customer": strconv.Itoa(int(c.ID))})
			return
		}

		defer response.Body.Close()
	}
	return
}

func todayZeroTime() time.Time {
	year, month, d := time.Now().Date()
	return time.Date(year, month, d, 0, 0, 0, 0, time.Local)
}

// IsNewCustomerToday ...
func (c *Customer) IsNewCustomerToday() bool {
	today := todayZeroTime()

	if c.ID != 0 && c.CaptureAt.After(today) {
		return false
	}

	return true
}

// IsNewCustomerTodayInShop customer是会员的话有相应会员记录，如果没有会员记录就肯定不是了，因为会在前面被过滤掉（仅限门店）
func (c *Customer) IsNewCustomerTodayInShop(tx *gorm.DB, shop *Shop) bool {
	var viprecord VipRecord

	today := todayZeroTime()

	tx.Where(&VipRecord{CustomerID: c.ID, ShopID: shop.ID, Date: today}).First(&viprecord)

	if viprecord.ID != 0 {
		return false
	}

	return true
}

func (c *Customer) IsVIP() bool {
	return c.CustomerGroupID != 0
}

// AfterCreate ...
func (c *Customer) AfterCreate(tx *gorm.DB) (err error) {
	var group CustomerGroup

	err = c.getGroup(tx, &group)
	if err != nil {
		return
	}

	if group.ID > 0 {
		group.CustomerCount++
		err = tx.Save(&group).Error
	}

	return
}

// AfterDelete ...
func (c *Customer) AfterDelete(tx *gorm.DB) (err error) {
	var group CustomerGroup

	err = c.getGroup(tx, &group)
	if err != nil {
		return
	}

	if group.ID > 0 {
		if group.CustomerCount > 0 {
			group.CustomerCount--
		}

		err = tx.Save(&group).Error
		if err != nil {
			return
		}
	}

	// delete relative vip_records
	err = tx.Where("customer_id = ?", c.ID).Delete(VipRecord{}).Error

	return
}

// 这个函数指定了customer碰到一个新event的时候的数据更新方式
func (customer *Customer) UpdateWithNewEvent(tx *gorm.DB, newEvent *Event, device *Device) error {
	if customer.ID == 0 {
		return errors.New("顾客id异常为0")
	}

	var columns Customer
	columns.EventsCount = customer.EventsCount + 1
	columns.LastCaptureAt = customer.CaptureAt
	columns.CaptureAt = newEvent.CaptureAt
	columns.OriginalFaceUrl = newEvent.OriginalFace

	columns.LastEventDeviceName = device.Name

	if columns.EventsCount == 1 {
		columns.Gender = newEvent.Gender
	}

	if device.ShopID != 0 { // 区分设备是shop架构还是区域架构
		columns.LastEventShopName = device.Shop.Name
	} else {
		columns.LastEventShopName = device.SmRegion.Name
	}

	if newEvent.Age != 0 {
		columns.Age = (customer.Age*(columns.EventsCount-1) + newEvent.Age) / columns.EventsCount
	}

	// 使用updateColumns更新有好处，不会更新对应的CustomerGroup
	if err := tx.Model(customer).UpdateColumns(columns).Error; err != nil {
		logger.Error(nil, "grpc", "event", "save customer failed", err.Error())
		return err
	}

	return nil
}
