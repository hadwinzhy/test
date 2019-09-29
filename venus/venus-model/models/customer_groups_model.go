package models

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"siren/venus/venus-model/models/connectors"
	"siren/venus/venus-model/models/logger"

	"github.com/jinzhu/gorm"
)

var (
	CUSTOMER_GROUP_TYPE_POTENTIAL = "potential"
	CUSTOMER_GROUP_TYPE_NORMAL    = "normal"
	CUSTOMER_GROUP_TYPE_BLACKLIST = "blacklist"
)

// CustomerGroup ...
type CustomerGroup struct {
	BaseModel
	Name           string `gorm:"type:varchar(50);not null" json:"name"`
	CompanyID      uint   `gorm:"index;not null" json:"company_id"`
	GroupID        string `gorm:"index" json:"group_id"`
	CreatedBy      string `gorm:"index;not null" json:"created_by"`
	CustomerCount  uint   `json:"customer_count"`
	DeleteOk       int    `gorm:"default:0" json:"delete_ok"`
	GroupType      string `gorm:"default:'normal'" json:"group_type"`
	Customers      []Customer
	Company        Company
	CompanyGroupID uint
}

// CustomerGroupSerializer ...
type CustomerGroupBasicSerializer struct {
	ID             uint      `json:"id"`
	Name           string    `json:"name"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	CompanyID      uint      `json:"company_id"`
	CompanyName    string    `json:"company_name"`
	CreatedBy      string    `json:"created_by"`
	CustomerCount  uint      `json:"customer_count"`
	GroupType      string    `json:"group_type"`
	CompanyGroupID uint      `json:"company_group_id"`
}

// BasicSerialize ...
func (group *CustomerGroup) BasicSerialize() CustomerGroupBasicSerializer {

	return CustomerGroupBasicSerializer{
		ID:             group.ID,
		Name:           group.Name,
		CreatedAt:      group.CreatedAt,
		UpdatedAt:      group.UpdatedAt,
		CreatedBy:      group.CreatedBy,
		CompanyID:      group.CompanyID,
		CompanyName:    group.Company.Name,
		CustomerCount:  group.CustomerCount,
		GroupType:      group.GroupType,
		CompanyGroupID: group.CompanyGroupID,
	}
}

// IsVipGroup ...
func (group *CustomerGroup) IsVipGroup() bool {
	if group.CreatedBy == "user" && group.GroupType == CUSTOMER_GROUP_TYPE_NORMAL {
		return true
	}
	return false
}

type groupCreateUUIDReq struct {
	Name         string   `json:"name"`
	MacAddresses []string `json:"mac_addresses"`
	Extension    string   `json:"extension"`
	HostName     string   `json:"host_name"`
}

type groupCreateUUIDResp struct {
	ID        int        `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"DeletedAt"`
	Name      string     `json:"name"`
	GroupUUID string     `json:"group_uuid"`
	IsUser    bool       `json:"is_user"`
	Date      time.Time  `json:"date"`
}

// Delete customer group
func (group *CustomerGroup) Delete(tx *gorm.DB) (err error) {
	gid := group.ID
	tx.Where("customer_group_id = ?", gid).Delete(Customer{})
	tx.Where("customer_group_id = ?", gid).Delete(VipRecord{})

	if err = tx.Delete(group).Error; err != nil {
		return
	}
	return
}

// CreatePandoraGroup ...
func CreatePandoraGroup(name string, macAddresses []string, extension string, hostName string) (groupUUID string, err error) {
	req := groupCreateUUIDReq{
		Name:         name,
		MacAddresses: macAddresses,
		Extension:    extension,
		HostName:     hostName,
	}

	var reqJSON []byte

	if reqJSON, err = json.Marshal(req); err != nil {
		return
	}

	// call pandora to create bindings
	reqReader := bytes.NewBuffer(reqJSON)
	var response http.Response
	response, err = connectors.HTTPRequest("POST", "/v1/api/groups", reqReader)

	if err != nil {
		return
	}

	var bodyContent []byte
	var x groupCreateUUIDResp
	defer response.Body.Close()

	bodyContent, _ = ioutil.ReadAll(response.Body)
	json.Unmarshal(bodyContent, &x)

	if x.GroupUUID != "" {
		groupUUID = x.GroupUUID
		err = nil
		return
	}

	err = errors.New("fail generating group_uuid. response:" + string(bodyContent))
	return
}

func DeletePandoraGroup(groupUUID string) (err error) {
	response, err := connectors.HTTPRequest("DELETE", "/v1/api/groups?group_uuid="+groupUUID, nil)

	defer response.Body.Close()
	if response.StatusCode != 200 {
		bodyContent, _ := ioutil.ReadAll(response.Body)
		if err == nil {
			err = errors.New("pandora group failed")
		}
		logger.Error(nil, "customer_group", "delete", "DELETE_GROUP_IN_PANDORA_ERROR: ", err, string(bodyContent))
		// raven.CaptureError(err, map[string]string{
		// 	"action":   "DELETE_GROUP_IN_PANDORA_ERROR",
		// 	"group_id": groupUUID,
		// 	"response": string(bodyContent),
		// })
		return
	}

	return
}

type groupSyncReq struct {
	GroupUUIDs         []string `json:"group_uuids"`
	DeviceMacAddresses []string `json:"device_mac_addresses" binding:"required"`
}

type groupUnsyncReq struct {
	GroupUUIDs   []string `json:"group_uuids"`
	MacAddresses []string `json"mac_addresses"`
}

// UnsyncPandoraGroups 解绑部分的关系
func UnsyncPandoraGroups(groupUUIDs []string, macAddresses []string) (err error) {
	req := groupUnsyncReq{
		GroupUUIDs:   groupUUIDs,
		MacAddresses: macAddresses,
	}

	var reqJSON []byte

	if reqJSON, err = json.Marshal(req); err != nil {
		return
	}

	// call pandora to create bindings
	reqReader := bytes.NewBuffer(reqJSON)
	var response http.Response
	response, err = connectors.HTTPRequest("POST", "/v1/api/groups/remove_device", reqReader)

	defer response.Body.Close()
	if response.StatusCode != 200 {
		bodyContent, _ := ioutil.ReadAll(response.Body)
		if err == nil {
			err = errors.New("pandora group failed")
		}
		logger.Error(nil, "customer_group", "sync", "UNSYNC_GROUP_IN_PANDORA_ERROR: ", err, string(bodyContent))
		// raven.CaptureError(err, map[string]string{
		// 	"action":        "DELETE_GROUP_IN_PANDORA_ERROR",
		// 	"group_uuids":   fmt.Sprintf("%v", groupUUIDs),
		// 	"mac_addresses": fmt.Sprintf("%v", macAddresses),
		// 	"response":      string(bodyContent),
		// })
		return
	}

	return
}

// SyncPandoraGroups ...
func SyncPandoraGroups(groupUUIDs []string, macAddresses []string) (err error) {
	req := groupSyncReq{
		GroupUUIDs:         groupUUIDs,
		DeviceMacAddresses: macAddresses,
	}

	var reqJSON []byte

	if reqJSON, err = json.Marshal(req); err != nil {
		return
	}

	// call pandora to create bindings
	reqReader := bytes.NewBuffer(reqJSON)
	var response http.Response
	response, err = connectors.HTTPRequest("POST", "/v1/api/groups_sync", reqReader)

	defer response.Body.Close()
	if response.StatusCode != 200 {
		bodyContent, _ := ioutil.ReadAll(response.Body)
		if err == nil {
			err = errors.New("pandora group failed")
		}
		logger.Error(nil, "customer_group", "sync", "SYNC_GROUP_IN_PANDORA_ERROR: ", err, string(bodyContent))
		// raven.CaptureError(err, map[string]string{
		// 	"action":        "DELETE_GROUP_IN_PANDORA_ERROR",
		// 	"group_uuids":   fmt.Sprintf("%v", groupUUIDs),
		// 	"mac_addresses": fmt.Sprintf("%v", macAddresses),
		// 	"response":      string(bodyContent),
		// })
		return
	}

	return
}

// BindUUID ...
func (group *CustomerGroup) BindUUID(tx *gorm.DB, company Company) (err error) {
	var devices []Device
	if group.CompanyGroupID != 0 {
		var companies []Company
		tx.Where("company_group_id = ?", group.CompanyGroupID).Find(&companies)
		var companyIds []uint
		for _, i := range companies {
			companyIds = append(companyIds, i.ID)
		}
		tx.Where("company_id in (?)", companyIds).Find(&devices)
	} else {
		tx.Where("company_id = ?", company.ID).Find(&devices)
	}

	macAddresses := make([]string, len(devices))
	for i, device := range devices {
		macAddresses[i] = device.MacAddress
	}

	hostName := ""
	var companyConfig CompanyConfig
	if group.CompanyGroupID == 0 {
		tx.Where("company_id = ?", company.ID).First(&companyConfig)
	} else {
		tx.First(&companyConfig)
	}
	if companyConfig.IDShow && company.ID > 300 {
		hostName = "rt"
	} else {
		hostName = companyConfig.TitanHostName
	}

	var groupUUID string
	groupUUID, err = CreatePandoraGroup(group.Name, macAddresses, "", hostName)

	if err != nil {
		logger.Error(nil, "customer_group", "create", "CREATE_GROUP_IN_PANDORA_ERROR: ", err)
		// raven.CaptureError(err, map[string]string{
		// 	"action":     "CREATE_GROUP_IN_PANDORA_ERROR",
		// 	"company_id": strconv.Itoa(int(company.ID)),
		// 	"group_name": group.Name,
		// 	"detail":     err.Error(),
		// })
		return
	}

	group.GroupID = groupUUID

	return
}

// UnBindUUID ...
func (group *CustomerGroup) UnBindUUID() (err error) {
	DeletePandoraGroup(group.GroupID)
	// not regard of response
	return
}

// AfterDelete callback
func (group *CustomerGroup) AfterDelete() (err error) {
	err = group.UnBindUUID()
	if err != nil {
		return
	}

	return
}
