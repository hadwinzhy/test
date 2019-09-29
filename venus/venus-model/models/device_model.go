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
)

var (
	TYPE_YUEKE          = "yueke"
	TYPE_LIGHTYEAR_FLOW = "lightyear_flow"
	TYPE_LIGHTYEAR_SNAP = "lightyear_snap"
	TYPE_YUEKE2_FLOW    = "yueke2_flow"
	TYPE_USB_SNAP       = "usb_snap"

	DEVICE_STATUS_ONLINE  = 1
	DEVICE_STATUS_OFFLINE = -1
)

var httpClient = &http.Client{}

// VideoPushAccountModel ...
type VideoPushAccountModel interface {
	BasicSerialize(interface{})
	RefSerialize(interface{})
}

// Device ...
type Device struct {
	BaseModel
	Name        string `gorm:"type:varchar(50);not null" json:"name"`
	DeviceType  string `gorm:"type:varchar(50);not null" json:"device_type"`
	CompanyID   uint   `gorm:"index;not null" json:"company_id"`
	MacAddress  string `gorm:"type:varchar(50);index;not null" json:"mac_address"`
	VersionCode string `gorm:"type:varchar(30);" json:"version_code"`
	Location    string `gorm:"type:varchar(50);" json:"location"`
	DeleteToken uint   `gorm:"type:integer;" json:"delete_token"`
	DeviceLot   string `gorm:"type:varchar" json:"device_lot"`

	NetworkStatus     string    `gorm:"type:varchar(20);" json:"network_status"`
	NetWorkName       string    `gorm:"type:varchar;" json:"network_name"`
	IPAddr            string    `gorm:"type:varchar;" json:"ip_addr"`
	LastHeartBeatTime time.Time `gorm:"type:timestamp with time zone" json:"last_heart_beat_time"`
	YueKeMeta         string    `gorm:"type:varchar" json:"yueke_meta"`
	Status            int       `gorm:"type:integer;default:0" json:"status"`
	Versions          string    `gorm:"type:varchar(255)" json:"versions"`

	Company       Company
	XcloudAccount XcloudAccount
	SailanAccount SailanAccount

	// shopID FloorID shopType专用于旧版逻辑
	ShopID   uint   `gorm:"not null" json:"shop_id"`
	ShopType string `gorm:"type:varchar(30)" json:"shop_type"`
	FloorID  uint   `gorm:"index" json:"floor_id"`
	Shop     Shop
	Floor    Floor

	// SmRegionID 用于区域多对多逻辑
	SmRegionID uint `gorm:"index" json:"sm_region_id"`
	SmRegion   SmRegion
}

type DeviceVersionInfoSerializer struct {
	RsHisiVer  string `json:"rs_hisi_ver"`
	FxDevFwVer string `json:"fx_dev_fw_ver"`
	HtFwVer    string `json:"ht_fw_ver"`
}

// DeviceBasicSerializer ...
type DeviceBasicSerializer struct {
	ID                uint                        `json:"id"`
	CreatedAt         time.Time                   `json:"created_at"`
	UpdateAt          time.Time                   `json:"updated_at"`
	Name              string                      `json:"name"`
	ShopID            uint                        `json:"shop_id"`
	ShopType          string                      `json:"shop_type"`
	ShopName          string                      `json:"shop_name"`
	MacAddress        string                      `json:"mac_address"`
	XcloudAccount     *XcloudAccountSerializer    `json:"xcloud_account"`
	SailanAccount     *SailanAccountSerializer    `json:"sailan_account"`
	DeviceType        string                      `json:"device_type"`
	DeviceLot         string                      `json:"device_lot"`
	VersionCode       string                      `json:"version_code"`
	Location          string                      `json:"location"`
	FloorID           uint                        `json:"floor_id"`
	FloorName         string                      `json:"floor_name"`
	LastHeartBeatTime time.Time                   `json:"last_heart_beat_time"`
	YuekeMeta         map[string]string           `json:"yueke_meta"`
	Online            bool                        `json:"online"`
	Versions          DeviceVersionInfoSerializer `json:"versions"`
}

func (d *Device) getOnlineStatus() bool {
	if d.Status == DEVICE_STATUS_ONLINE {
		return true
	}

	now := time.Now()
	difference := now.Sub(d.LastHeartBeatTime)

	if difference.Minutes() > 20 {
		return false
	}
	return true

}

func handleCallBacks(handlerList []func() error) (err error) {
	for _, handler := range handlerList {
		err = handler()
		if err != nil {
			break
		}
	}

	return
}

func detailYuekeMeta(value string) map[string]string {
	if value == "" || !strings.Contains(value, "{") {
		return nil
	}
	newValueReplacer := strings.NewReplacer("{", "", "}", "", " ", "")
	newValue := newValueReplacer.Replace(value)
	newValueList := strings.Split(newValue, ",")

	var returnValue = make(map[string]string)
	for _, key := range newValueList {
		keyList := strings.Split(strings.TrimSpace(key), ":")
		returnValue[keyList[0]] = keyList[1]
	}
	return returnValue
}

// DeviceInfoSerializer ... used by app
type DeviceInfoSerializer struct {
	NetworkStatus string `json:"network_status"`
	NetWorkName   string `json:"network_name"`
	IPAddress     string `json:"ip_addr"`
}

func (d *Device) versionDetail() DeviceVersionInfoSerializer {
	var version DeviceVersionInfoSerializer
	json.Unmarshal([]byte(d.Versions), &version)
	return version
}

// BasicSerialize ...
func (d *Device) BasicSerialize() DeviceBasicSerializer {

	return DeviceBasicSerializer{
		ID:                d.ID,
		Name:              d.Name,
		ShopID:            d.ShopID,
		ShopType:          d.ShopType,
		ShopName:          d.Shop.Name,
		MacAddress:        d.MacAddress,
		CreatedAt:         d.CreatedAt.Round(time.Second),
		UpdateAt:          d.UpdatedAt.Round(time.Second),
		LastHeartBeatTime: d.LastHeartBeatTime.Round(time.Second),
		XcloudAccount:     d.XcloudAccount.RefSerialize(),
		SailanAccount:     d.SailanAccount.RefSerialize(),
		DeviceType:        d.DeviceType,
		DeviceLot:         d.DeviceLot,
		VersionCode:       d.VersionCode,
		Location:          d.Location,
		FloorID:           d.FloorID,
		FloorName:         d.Floor.Name,
		Online:            d.getOnlineStatus(),
		YuekeMeta:         detailYuekeMeta(d.YueKeMeta),
		Versions:          d.versionDetail(),
	}
}

// DeviceFullSerializer ...
type DeviceFullSerializer struct {
	DeviceBasicSerializer
	CompanyID   uint   `json:"company_id"`
	CompanyName string `json:"company_name"`
}

// FullSerialize ...
func (d *Device) FullSerialize() DeviceFullSerializer {
	return DeviceFullSerializer{
		DeviceBasicSerializer: d.BasicSerialize(),
		CompanyName:           d.Company.Name,
	}
}

// DeviceInfoSerialize ...
func (d *Device) DeviceInfoSerialize() DeviceInfoSerializer {
	state := d.NetworkStatus
	if state == "" {
		state = "network_offline"
	}
	return DeviceInfoSerializer{
		NetworkStatus: state,
		NetWorkName:   d.NetWorkName,
		IPAddress:     d.IPAddr,
	}
}

type DeviceShoppingMallSerializer struct {
	BaseSerializer
	Name              string                      `json:"name"`
	DeviceLot         string                      `json:"device_lot"`
	DeviceType        string                      `json:"device_type"`
	FloorID           uint                        `json:"floor_id"`
	RegionID          uint                        `json:"region_id"`
	Region            SmRegionBasicSerializer     `json:"region"`
	LastHeartBeatTime time.Time                   `json:"last_heart_beat_time"`
	Online            bool                        `json:"online"`
	Location          string                      `json:"location"`
	MacAddress        string                      `json:"mac_address"`
	SailanAccount     *SailanAccountSerializer    `json:"sailan_account"`
	VersionCode       string                      `json:"version_code"`
	YuekeMeta         map[string]string           `json:"yueke_meta"`
	Versions          DeviceVersionInfoSerializer `json:"versions"`
	CompanyName       string                      `json:"company_name"`
}

func (d *Device) ShoppingMallSerialize() DeviceShoppingMallSerializer {
	return DeviceShoppingMallSerializer{
		BaseSerializer:    d.BaseModel.Serialize(),
		CompanyName:       d.Company.Name,
		Name:              d.Name,
		DeviceLot:         d.DeviceLot,
		DeviceType:        d.DeviceType,
		FloorID:           d.SmRegion.FloorID,
		RegionID:          d.SmRegionID,
		Region:            d.SmRegion.BasicSerialize(),
		LastHeartBeatTime: d.LastHeartBeatTime,
		Online:            d.getOnlineStatus(),
		YuekeMeta:         detailYuekeMeta(d.YueKeMeta),
		Location:          d.Location,
		MacAddress:        d.MacAddress,
		SailanAccount:     d.SailanAccount.RefSerialize(),
		VersionCode:       d.VersionCode,
		Versions:          d.versionDetail(),
	}
}

func (d *Device) handleSailanAccountAfterCreate(tx *gorm.DB) func() error {
	return func() (err error) {
		if !strings.Contains(d.DeviceType, "yueke") && !strings.Contains(d.DeviceType, "lightyear") {
			return
		}

		var account SailanAccount
		searchErr := tx.Where("mac_address = ?", d.MacAddress).First(&account).Error
		if searchErr != nil {
			searchErr = tx.Where("mac_address is NULL").First(&account).Error
		}

		if account.ID > 0 {
			account.IsUsed = true
			account.DeviceID = d.ID
			account.MacAddress = d.MacAddress

			// update its device id and used state
			err = tx.Save(&account).Error
			d.SailanAccount = account
		} else {
			logger.Error(nil, "device", "create", "no video push accounts left")
			err = errors.New("no video push accounts left")
		}
		return
	}
}

func (d *Device) handleShopDeviceAfterCreate(tx *gorm.DB, platform string) func() error {
	return func() (err error) {
		var shop Shop
		tx.Model(d).Related(&shop)

		if shop.ID > 0 {
			if platform == "pandora" {
				shop.PandoraDeviceCount++
			} else {
				shop.DeviceCount++
			}

			err = tx.Save(&shop).Error
		}

		return
	}
}

type bindReqBody struct {
	DeviceType      string   `json:"device_type"`
	MacAddress      string   `json:"mac_address"`
	GroupUUIDs      []string `json:"group_uuids"`
	DevicePackUUIDs []string `json:"device_pack_uuids"`
	ShopUUID        string   `json:"shop_uuid"`
	HostName        string   `json:"host_name"`
	NormalThreshold uint     `json:"normal_threshold"`
	VIPThreshold    uint     `json:"vip_threshold"`
}

func (d *Device) handleUUIDBindAfterCreate(tx *gorm.DB) func() error {
	return func() (err error) {
		// 判断是新逻辑还是老逻辑，新逻辑给他绑device_pack_uuid，老逻辑给它绑shop_uuid
		var company Company
		var shop Shop
		var devicePackUUIDs []string

		if err = tx.Model(d).Related(&company).Error; err != nil {
			return
		}

		tx.Model(d).Related(&shop)

		if shop.ID == 0 { // 没有shop，则说明使用的是新逻辑
			devicePackUUIDs = []string{
				company.DevicePackUUID,
			}
		}

		var groups []CustomerGroup
		err = tx.Model(&company).Related(&groups).Error
		if err != nil {
			return
		}

		if company.CompanyGroupID != 0 { // 加上集团组
			var companyGroupGroups []CustomerGroup
			err = tx.Where("company_group_id = ?", company.CompanyGroupID).Find(&companyGroupGroups).Error
			if err != nil {
				return
			}

			groups = append(groups, companyGroupGroups...)
		}

		hostName := ""
		var companyConfig CompanyConfig
		tx.Where("company_id = ?", company.ID).First(&companyConfig)

		if companyConfig.IDShow && company.ID > 300 {
			hostName = "rt" // 阅小客有实时要求
		} else {
			hostName = companyConfig.TitanHostName
		}

		uuids := make([]string, len(groups))

		for _, group := range groups {
			uuids = append(uuids, group.GroupID)
		}

		req := bindReqBody{
			DeviceType:      d.DeviceType,
			MacAddress:      d.MacAddress,
			GroupUUIDs:      uuids,
			HostName:        hostName,
			NormalThreshold: companyConfig.NormalThreshold,
			VIPThreshold:    companyConfig.VIPThreshold,
		}

		if shop.ID == 0 { // 新逻辑
			req.DevicePackUUIDs = devicePackUUIDs
		} else { // 老逻辑
			req.ShopUUID = shop.ShopUUID
		}

		var reqJSON []byte

		if reqJSON, err = json.Marshal(req); err != nil {
			return
		}
		fmt.Println(string(reqJSON))

		// call pandora to create bindings
		reqReader := bytes.NewBuffer(reqJSON)
		var response http.Response

		response, err = connectors.HTTPRequest("POST", "/v1/api/devices", reqReader)
		if err != nil {
			logger.Error(nil, "device", "create", "CREATE_DEVICE_IN_PANDORA_ERROR, device_macaddress= ", d.MacAddress, " error: ", err)
			raven.CaptureError(err, map[string]string{"action": "CREATE_DEVICE_IN_PANDORA_ERROR", "customer": d.MacAddress})
			return
		}

		// call post to pandora failed
		if response.StatusCode != 200 {
			buf := new(bytes.Buffer)
			buf.ReadFrom(response.Body)
			newStr := buf.String()
			err = errors.New(newStr)
		}

		return
	}
}

func (d *Device) handleFloorDeviceAfterCreate(tx *gorm.DB, platform string) func() error {
	return func() (err error) {
		var tempFloor Floor

		if d.FloorID == 0 {
			return
		}
		err = tx.First(&tempFloor, d.FloorID).Error

		fmt.Println(tempFloor)
		if err != nil {
			return
		}

		if platform == "pandora" {
			tempFloor.PandoraDeviceCount++
		} else {
			tempFloor.DeviceCount++
		}

		fmt.Println("before save", tempFloor.DeviceCount)
		err = tx.Save(&tempFloor).Error

		return
	}
}

func (d *Device) getRelatedSmRegion(tx *gorm.DB) {
	if d.SmRegion.ID == 0 {
		var region SmRegion
		tx.First(&region, d.SmRegionID)
		d.SmRegion = region
	}
}

func (d *Device) handleSmShopsDevice(tx *gorm.DB, plus int) (err error) {
	for _, group := range d.SmRegion.SmUvGroups {
		if group.RegionType == UVGROUP_REGION_TYPE_SHOPENTRANCE {
			var shop SmShop
			tx.First(&shop, group.RelatedID)
			if shop.ID == 0 {
				continue
			}
			if plus > 0 {
				shop.DeviceCount++
			}
			if plus < 0 {
				shop.DeviceCount--
			}
			if err = tx.Save(&shop).Error; err != nil {
				return
			}
		}
	}
	return
}

func (d *Device) handleSmRegionDeviceAfterCreate(tx *gorm.DB) func() error {
	return func() (err error) {
		d.getRelatedSmRegion(tx)
		if d.SmRegion.ID == 0 {
			logger.Error(nil, "device", "create", "device_count error region not existed. ", d.MacAddress, d.SmRegionID)
			return
		}
		d.SmRegion.DeviceCount++

		// 给region加上devicecount
		if err = tx.Save(&d.SmRegion).Error; err != nil {
			return
		}

		// 取对应的去重组
		var uvGroups []SmUvGroup
		if err = tx.Model(&d.SmRegion).Association("SmUvGroups").Find(&uvGroups).Error; err != nil {
			return
		}
		d.SmRegion.SmUvGroups = uvGroups

		err = d.handleSmShopsDevice(tx, 1)
		return
	}
}

func (d *Device) handleSmRegionDeviceAfterDelete(tx *gorm.DB) func() error {
	return func() (err error) {
		d.getRelatedSmRegion(tx)
		if d.SmRegion.ID == 0 {
			logger.Error(nil, "device", "create", "device_count error region not existed. ", d.MacAddress, d.SmRegionID)
			return
		}
		d.SmRegion.DeviceCount--

		// 给region加上devicecount
		if err = tx.Save(&d.SmRegion).Error; err != nil {
			return
		}

		// 取对应的去重组
		var uvGroups []SmUvGroup
		if err = tx.Model(&d.SmRegion).Association("SmUvGroups").Find(&uvGroups).Error; err != nil {
			return
		}
		d.SmRegion.SmUvGroups = uvGroups

		err = d.handleSmShopsDevice(tx, -1)
		return
	}
}

func (d *Device) Save(tx *gorm.DB) error {
	createFlag := false
	if d.ID == 0 {
		createFlag = true
	}

	if err := tx.Save(d).Error; err != nil {
		return err
	}

	if createFlag {
		return d.afterCreate(tx)
	}

	return nil
}

// AfterCreate is after create callback:
// after creating the device, add an account on it
func (d *Device) afterCreate(tx *gorm.DB) (err error) {
	var funcs [](func() error)
	if d.SmRegionID == 0 { // 原先逻辑
		if d.DeviceType == "pandora" { // TODO pandora有了店铺之后，可以删除
			funcs = []func() error{
				d.handleShopDeviceAfterCreate(tx, "pandora"),
				d.handleFloorDeviceAfterCreate(tx, "pandora"),
			}
		} else {
			funcs = []func() error{
				d.handleSailanAccountAfterCreate(tx),
				d.handleFloorDeviceAfterCreate(tx, "normal"),
				d.handleShopDeviceAfterCreate(tx, "normal"),
				d.handleUUIDBindAfterCreate(tx),
			}
		}
	} else { // 新的shopping mall的hook
		funcs = []func() error{
			d.handleSailanAccountAfterCreate(tx),
			d.handleSmRegionDeviceAfterCreate(tx),
			d.handleUUIDBindAfterCreate(tx),
		}
	}

	err = handleCallBacks(funcs)
	fmt.Println(err)
	if err == nil {
		d.createEventCallback(tx)
	}
	return
}

func (d *Device) handleShopDeviceAfterDelete(tx *gorm.DB, platform string) func() error {
	return func() (err error) {
		var shop Shop
		tx.Model(d).Related(&shop)

		if shop.ID > 0 {
			if platform == "normal" && shop.DeviceCount > 0 {
				shop.DeviceCount--
			}
			if platform == "pandora" && shop.PandoraDeviceCount > 0 {
				shop.PandoraDeviceCount--
			}
			err = tx.Save(&shop).Error
		}

		return
	}
}

func (d *Device) handleUUIDUnBindAfterDelete(tx *gorm.DB) func() error {
	return func() (err error) {
		// ignore error in delete pandora device
		response, err := connectors.HTTPRequest("DELETE", "/v1/api/devices/"+d.MacAddress, nil)
		if err != nil || response.StatusCode != 200 {
			logger.Error(nil, "device", "delete", "DELETE_DEVICE_IN_PANDORA_ERROR, device_macaddress= ", d.MacAddress, " error: ", err)
			raven.CaptureError(err, map[string]string{"action": "DELETE_DEVICE_IN_PANDORA_ERROR", "customer": d.MacAddress})
			return
		}

		return
	}
}

func (d *Device) handleFloorDeviceAfterDelete(tx *gorm.DB, platform string) func() error {
	return func() (err error) {
		var tempFloor Floor

		if d.FloorID == 0 {
			return
		}
		query := tx.Where("id = ?", d.FloorID).First(&tempFloor)

		if platform == "normal" && tempFloor.DeviceCount > 0 {
			err = query.Update("device_count", tempFloor.DeviceCount-1).Error
		}
		if platform == "pandora" && tempFloor.PandoraDeviceCount > 0 {
			err = query.Update("device_count", tempFloor.PandoraDeviceCount-1).Error
		}

		return
	}
}

func (d *Device) Delete(tx *gorm.DB) error {
	d.DeleteToken = d.ID

	if d.ID == 0 {
		return nil
	}

	if err := tx.Save(d).Error; err != nil {
		return err
	}

	if err := tx.Delete(d).Error; err != nil {
		return err
	}

	return d.afterDelete(tx)
}

// AfterDelete callback: after delete device, release the account
func (d *Device) afterDelete(tx *gorm.DB) (err error) {
	var deviceType string
	if d.DeviceType == "pandora" {
		deviceType = "pandora"
	} else {
		deviceType = "normal"
	}
	funcs := []func() error{
		// d.handleSailanAccountAfterDelete,
		d.handleShopDeviceAfterDelete(tx, deviceType),
		d.handleUUIDUnBindAfterDelete(tx),
		d.handleFloorDeviceAfterDelete(tx, deviceType),
		d.handleSmRegionDeviceAfterDelete(tx), // 新逻辑
	}
	err = handleCallBacks(funcs)

	if err == nil {
		d.deleteEventCallback(tx)
	}

	return
}

func (d *Device) manuallyLoadCompanyAndShop(tx *gorm.DB) {
	if d.Shop.ID == 0 {
		var shop Shop
		tx.First(&shop, d.ShopID)
		d.Shop = shop
	}

	if d.Company.ID == 0 {
		var company Company
		tx.First(&company, d.CompanyID)
		d.Company = company
	}

}

func (d *Device) sendCallback(tx *gorm.DB, companyConfig CompanyConfig, action string) {
	d.manuallyLoadCompanyAndShop(tx)
	eventTime := time.Now().Format("2006-01-02 15:04:05")
	data := string(`{
        "action": "` + action + `",
        "event_time": "` + eventTime + `",
        "device": {
            "id": ` + strconv.Itoa(int(d.ID)) + `,
            "shop_id": ` + strconv.Itoa(int(d.ShopID)) + `,
            "shop_name": "` + d.Shop.Name + `",
            "company_id": ` + strconv.Itoa(int(d.CompanyID)) + `,
            "company_name": "` + d.Company.Name + `",
            "name": "` + d.Name + `",
            "device_type": "` + d.DeviceType + `",
            "mac_address": "` + d.MacAddress + `"
        }
    }`)

	fmt.Println(data)
	for i := 0; i < 5; i++ {

		reqReader := bytes.NewBuffer([]byte(data))

		request, _ := http.NewRequest("POST", companyConfig.DeviceCallback, reqReader)
		headers := map[string]*string(companyConfig.Headers)
		request.Header = http.Header{
			"Content-Type": []string{"application/json"},
		}

		for headerKey, headerValue := range headers {
			request.Header.Add(headerKey, *headerValue)
		}

		responsePtr, err := httpClient.Do(request)

		if err != nil {
			fmt.Println(err)
			logger.Error(nil, "device", "callback", err)
			return
		}

		defer responsePtr.Body.Close()

		bodyContent, _ := ioutil.ReadAll(responsePtr.Body)
		fmt.Println(string(bodyContent))
		fmt.Println(err)
		if err == nil {
			return
		} else {
			fmt.Println("failed, retrying...")
			time.Sleep(time.Second * time.Duration(i))
		}
	}
}

func (d *Device) eventCallback(tx *gorm.DB, action string) {
	var companyConfig CompanyConfig
	tx.Where("company_id = ?", d.CompanyID).First(&companyConfig)
	fmt.Println(companyConfig.ID, companyConfig.DeviceCallback)
	if companyConfig.ID != 0 && companyConfig.DeviceCallback != "" {
		go d.sendCallback(tx, companyConfig, action)
	}
}

func (d *Device) createEventCallback(tx *gorm.DB) {
	d.eventCallback(tx, "create")
}

func (d *Device) deleteEventCallback(tx *gorm.DB) {
	d.eventCallback(tx, "delete")
}

func GetStatusByString(statusStr string) int {
	if statusStr == "online" {
		return DEVICE_STATUS_ONLINE
	} else {
		return DEVICE_STATUS_OFFLINE
	}
}
