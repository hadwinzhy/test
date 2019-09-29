// Package models contains admin model associated bussiness logics.
package models

import (
	"crypto/rand"
	"fmt"
	"log"
	"strconv"
	"time"

	"siren/venus/venus-controller/controllers/errors"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

const (
	ADMIN_ROLE_SUPERADMIN   = "superadmin"
	ADMIN_ROLE_COMPANYGROUP = "company_group" // 集团 - 集团管理员
	ADMIN_ROLE_AREAMANAGER  = "areamanager"   // 集团 - 大区经理
	ADMIN_ROLE_DATAREADER   = "datareader"    // 集团 - 数据操作员
	ADMIN_ROLE_ADMIN        = "admin"
	ADMIN_ROLE_MANAGER      = "manager"
	ADMIN_ROLE_SHOPLEADER   = "shopleader"
	ADMIN_ROLE_EMPLOYEE     = "employee"
	ADMIN_STATE_PENDING     = "pending"
	ADMIN_STATE_APPROVED    = "approved"
	ADMIN_STATE_REJECTED    = "rejected"

	ADMIN_ROLE_IT        = "it"         // it 管理者
	ADMIN_ROLE_OPERATOR  = "operator"   // 运营人员 == manager
	ADMIN_ROLE_BRANDSHOP = "brandshops" // 品牌商铺 == shopleader
)

// 功能限制 限制功能的范围 shopping mall版的权限范围
// 前两位表示了一个大多数功能的等级概念，越高，则1的位数越多
// 后两位是功能FLAG，作为功能模块部分的预留
const (
	OP_SCOPE_ADMIN                 = 0xFF01
	OP_SCOPE_IT                    = 0xFF00
	OP_SCOPE_OPERATOR_ALL          = 0x7F01
	OP_SCOPE_OPERATOR_BUSINESSTYPE = 0x2F01
	OP_SCOPE_OPERATOR_SHOPS        = 0x1F01
	OP_SCOPE_BRANDSHOP             = 0x0700
)

var OP_SCOPE_LIST = []int{
	OP_SCOPE_ADMIN,
	OP_SCOPE_IT,
	OP_SCOPE_OPERATOR_ALL,
	OP_SCOPE_OPERATOR_BUSINESSTYPE,
	OP_SCOPE_OPERATOR_SHOPS,
	OP_SCOPE_BRANDSHOP,
}

var OP_SCOPE_MAP_PART = map[int]int{
	OP_SCOPE_ADMIN:        1,
	OP_SCOPE_IT:           2,
	OP_SCOPE_OPERATOR_ALL: 3,
}

const (
	LEVEL_SCOPE_ADMIN     = 0xFF00
	LEVEL_SCOPE_OPERATOR  = 0x0F00
	LEVEL_SCOPE_BRANDSHOP = 0x0700

	LEVEL_SCOPE_ALL = 0x7F00
)

var LEVEL_SCOPE_LIST = []int{
	LEVEL_SCOPE_ADMIN,
	LEVEL_SCOPE_OPERATOR,
	LEVEL_SCOPE_BRANDSHOP,
	LEVEL_SCOPE_ALL,
}

const (
	FUNC_ROLE_MANAGEMENT               = 0x0001
	FUNC_OPERATOR_BUSINESSTYPE_SPECIAL = 0x2000
	FUNC_OPERATOR_SHOPS_SPECIAL        = 0x1000
)

type ParsedAuthority struct {
	Scope                  int
	BusinessTypesResources []uint
	ShopsResources         []uint
	CompanyResources       []uint
}

// Admin represents object admin
type Admin struct {
	BaseModel
	Name              string         `gorm:"type:varchar(50)" json:"name"`
	Role              string         `gorm:"type:varchar(50);not null" json:"role"`
	AuthToken         string         `gorm:"type:varchar(100);not null;unique_index" json:"auth_token"`
	EncryptedPassword string         `gorm:"type:varchar(100);not null" json:"-"`
	Phone             string         `gorm:"type:varchar(20);not null;" json:"phone"`
	CompanyID         uint           `gorm:"index" json:"company_id"`
	CompanyGroupID    uint           `gorm:"index" json:"company_group_id"`
	ShopIDs           pq.Int64Array  `gorm:"type:integer[];" json:"shop_ids"`
	State             string         `gorm:"type:varchar(20)" json:"state"`
	Avatars           pq.StringArray `gorm:"type:varchar(255)[]" json:"avatars"`
	AuthorityID       uint           `gorm:"index" json:"authority_id"`
	Authority         ExtendedAuthority
	Company           Company
	CompanyGroup      CompanyGroup
	Shops             []Shop

	parsedAuthority ParsedAuthority
	GetuiFlag       string `gorm:"type:varchar;default:'on';column:getui_flag" json:"getui_flag"`
	WarningStatus   string `gorm:"type:varchar" json:"warning_status"`
}

// AdminBasicSerializer ...
type AdminBasicSerializer struct {
	ID                uint                          `json:"id"`
	Phone             string                        `json:"phone"`
	Role              string                        `json:"role"`
	AuthToken         string                        `json:"auth_token"`
	CreatedAt         time.Time                     `json:"created_at"`
	UpdatedAt         time.Time                     `json:"updated_at"`
	CompanyID         uint                          `json:"company_id"`
	CompanyName       string                        `json:"company_name"`
	CompanyBrandTitle string                        `json:"company_brand_title"`
	ShopIDs           []int64                       `json:"shop_ids"`
	Name              string                        `json:"name"`
	State             string                        `json:"state"`
	Avatars           []string                      `json:"avatars"`
	Shops             []ShopBasicSerializer         `json:"shops"`
	SailanAccount     string                        `json:"sailan_account"`
	SailanPassword    string                        `json:"sailan_password"`
	CompanyConfig     *CompanyConfigBasicSerializer `json:"company_config"`
	CompanyType       string                        `json:"company_type"`
	CompanyGroupID    uint                          `json:"company_group_id"`
	CompanyGroupName  string                        `json:"company_group_name"`
	// Profile   Profile        `json:"profile"`
}

func (admin *AdminBasicSerializer) APIBasePath() string {
	if admin.CompanyType == COMPANY_TYPE_SHOPPINGMALL || (admin.CompanyID == 0 && admin.CompanyGroupID != 0) {
		return "/v2/api/"
	} else {
		return "/v1/api/"
	}
}

type HomepageDefaultParam struct {
	ItemType string `json:"item_type"`
	Key      string `json:"key"`
	Value    string `json:"value"`
	Name     string `json:"name"`
	Title    string `json:"title"`
}

//
type AdminBasicSerializerV2 struct {
	AdminBasicSerializer
	APIBasePath      string                      `json:"api_base_path"`
	ExrendedAuhority ExtendedAuthoritySerializer `json:"extended_authority"`
	HomepageParam    HomepageDefaultParam        `json:"homepage_param"`
}

func (admin *Admin) HomepageParam(tx *gorm.DB) HomepageDefaultParam {
	if admin.parsedAuthority.Scope&LEVEL_SCOPE_ALL == LEVEL_SCOPE_ALL {
		return HomepageDefaultParam{
			Key:      "shortcut",
			Value:    "all",
			Name:     "all",
			ItemType: "shortcut",
			Title:    "全部",
		}
	} else if admin.parsedAuthority.Scope&FUNC_OPERATOR_BUSINESSTYPE_SPECIAL == FUNC_OPERATOR_BUSINESSTYPE_SPECIAL {
		return HomepageDefaultParam{
			Key:      "shortcut",
			Value:    "whole_business_type",
			Name:     "whole_business_type",
			ItemType: "shortcut",
			Title:    "全部",
		}
	} else if len(admin.parsedAuthority.ShopsResources) >= 1 {
		var shop SmShop
		tx.First(&shop, admin.parsedAuthority.ShopsResources[0])
		if shop.ID != 0 {
			shopIDStr := strconv.Itoa(int(admin.parsedAuthority.ShopsResources[0]))
			return HomepageDefaultParam{
				Key:      "shop_id",
				Value:    shopIDStr,
				Name:     "shop_" + shopIDStr,
				ItemType: "shop",
				Title:    shop.Name,
			}
		}
	}

	return HomepageDefaultParam{}
}

func (admin *Admin) basicSerializeWithShops(shops []ShopBasicSerializer) AdminBasicSerializer {
	companyType := admin.Company.Type
	if companyType == "" && admin.CompanyGroupID != 0 {
		companyType = "company_group"
	}

	result := AdminBasicSerializer{
		ID:                admin.ID,
		Phone:             admin.Phone,
		Role:              admin.Role,
		AuthToken:         admin.AuthToken,
		CreatedAt:         admin.CreatedAt,
		UpdatedAt:         admin.UpdatedAt,
		CompanyID:         admin.CompanyID,
		CompanyName:       admin.Company.Name,
		CompanyBrandTitle: admin.Company.BrandTitle,
		CompanyType:       companyType,
		ShopIDs:           admin.ShopIDs,
		Name:              admin.Name,
		State:             admin.State,
		Avatars:           admin.Avatars,
		Shops:             shops,
		SailanAccount:     "15692179960",
		SailanPassword:    "123456",
		CompanyConfig:     nil,
		CompanyGroupID:    admin.CompanyGroupID,
		CompanyGroupName:  admin.CompanyGroup.Name,
		// Profile:   profile,
	}

	var s CompanyConfigBasicSerializer
	if admin.Company.CompanyConfig.ID != 0 {
		s = admin.Company.CompanyConfig.BasicSerialize()
		result.CompanyConfig = &s
	}

	return result

}

// BasicSerialize ...
func (admin *Admin) BasicSerialize() AdminBasicSerializer {
	serializers := make([]ShopBasicSerializer, len(admin.Shops))
	for i, shop := range admin.Shops {
		serializers[i] = shop.BasicSerialize()
	}

	return admin.basicSerializeWithShops(serializers)
}

func (admin *Admin) CompanyConfigSerialize(companyConfig *CompanyConfig) AdminBasicSerializer {
	result := admin.BasicSerialize()
	if companyConfig == nil {
		return result
	}
	s := companyConfig.BasicSerialize()
	result.CompanyConfig = &s
	return result
}

// SerializeWithShops ...
func (admin *Admin) SerializeWithShops(tx *gorm.DB) AdminBasicSerializer {
	var shops []Shop
	if admin.RoleCheck(ADMIN_ROLE_ADMIN) {
		tx.Model(&admin.Company).Related(&shops)
	} else { // for other roles
		tx.Where("id in (?)", admin.ShopIDs).Find(&shops)
	}

	admin.Shops = shops
	return admin.BasicSerialize()
}

func tokenGenerator(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

// UpdatePassword will update password
func (self *Admin) UpdatePassword(tx *gorm.DB, newPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), 4)
	if err != nil {
		log.Fatal(err)
	}
	oldToken := self.AuthToken
	self.EncryptedPassword = string(hash)
	self.AuthToken = tokenGenerator(20)

	// when auth token refreshed, remove related getui accounts
	if oldToken != "" {
		var accounts []GetuiAccount
		tx.Where("auth_token = ?", oldToken).Delete(&accounts)
	}

	if dbs := tx.Save(self); dbs.Error != nil {
		return dbs.Error
	}

	return nil
}

// PasswordCheck only check password
func (admin *Admin) PasswordCheck(password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(admin.EncryptedPassword), []byte(password)); err == nil {
		return true
	} else {
		return false
	}
}

func role2level(role string) int {
	switch role {
	case ADMIN_ROLE_SUPERADMIN:
		return 5
	case ADMIN_ROLE_ADMIN, ADMIN_ROLE_IT:
		return 4
	case ADMIN_ROLE_MANAGER:
		return 3
	case ADMIN_ROLE_SHOPLEADER:
		return 2
	case ADMIN_ROLE_EMPLOYEE:
		return 1
	default:
		return 0
	}
}

// Check if the role is over a specific role
func (admin *Admin) RoleCheck(role string) bool {
	return role2level(admin.Role) >= role2level(role)
}

// RoleAbove can be used in creating accounts
func (admin *Admin) RoleAbove(role string) bool {
	if role == ADMIN_ROLE_IT && admin.Role == ADMIN_ROLE_ADMIN {
		return true
	}
	return role2level(admin.Role) > role2level(role)
}

func (admin *Admin) PendingErrCode() *errors.ErrorCode {
	switch admin.State {
	case ADMIN_STATE_APPROVED:
		return nil
	case ADMIN_STATE_PENDING:
		return &errors.ErrorAdminPending
	case ADMIN_STATE_REJECTED:
		return &errors.ErrorAdminRejected
	default:
		return nil
	}
}

func (admin *Admin) HasShop(shop Shop) bool {
	if admin.RoleAbove(ADMIN_ROLE_MANAGER) {
		return admin.CompanyID == shop.CompanyID
	} else {
		for _, shopID := range admin.ShopIDs {
			if uint(shopID) == shop.ID {
				return true
			}
		}
		return false
	}
}

// v2
type RegionAdminSerializer struct {
	ID          uint      `json:"id"`
	Phone       string    `json:"phone"`
	Role        string    `json:"role"`
	AuthToken   string    `json:"auth_token"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CompanyID   uint      `json:"company_id"`
	CompanyName string    `json:"company_name"`
	ShopIDs     []int64   `json:"shop_ids"`
	Name        string    `json:"name"`
}

func (a Admin) RegionAuthBasicSerializer() RegionAdminSerializer {
	return RegionAdminSerializer{
		ID:          a.ID,
		Phone:       a.Phone,
		Role:        a.Role,
		AuthToken:   a.AuthToken,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
		CompanyID:   a.CompanyID,
		CompanyName: a.Company.Name,
		ShopIDs:     a.ShopIDs,
		Name:        a.Name,
	}
}

func CheckRoleIt(role string) error {
	if role == ADMIN_ROLE_IT {
		return fmt.Errorf("role not allow")
	}
	return nil
}

// CalcParsedAuthorities 保存了权限scope，以及对应的资源id
func (admin *Admin) CalcParsedAuthorities() ParsedAuthority {
	if admin.Company.Type == "store" || admin.parsedAuthority.Scope != 0 {
		return admin.parsedAuthority
	}
	fmt.Println("calc")
	if admin.Role == "admin" || admin.Role == "company_group" || admin.Role == "areamanager" {
		admin.parsedAuthority.Scope = OP_SCOPE_ADMIN
	} else if admin.Role == "it" {
		admin.parsedAuthority.Scope = OP_SCOPE_IT
	} else if admin.Role == "operator" {
		if admin.Authority.PlaceType == AUTH_PLACE_TYPE_ALL {
			admin.parsedAuthority.Scope = OP_SCOPE_OPERATOR_ALL
		} else if admin.Authority.PlaceType == AUTH_PLACE_TYPE_BUSINESS_TYPE {
			admin.parsedAuthority.Scope = OP_SCOPE_OPERATOR_BUSINESSTYPE
			businessTypeIDs := make([]uint, len(admin.Authority.SmBusinessTypeIDs))
			for i := range businessTypeIDs {
				businessTypeIDs[i] = uint(admin.Authority.SmBusinessTypeIDs[i])
			}
			admin.parsedAuthority.BusinessTypesResources = businessTypeIDs
		} else {
			admin.parsedAuthority.Scope = OP_SCOPE_OPERATOR_SHOPS
			shopIDs := make([]uint, len(admin.Authority.SmShopIDs))
			for i := range shopIDs {
				shopIDs[i] = uint(admin.Authority.SmShopIDs[i])
			}
			admin.parsedAuthority.ShopsResources = shopIDs
		}
	} else {
		admin.parsedAuthority.Scope = OP_SCOPE_BRANDSHOP
		if admin.Authority.PlaceType == AUTH_PLACE_TYPE_BRANDSHOPS && len(admin.Authority.SmShopIDs) == 1 {
			admin.parsedAuthority.ShopsResources = []uint{uint(admin.Authority.SmShopIDs[0])}
		}
	}
	fmt.Println("parsed auth", admin.parsedAuthority)
	return admin.parsedAuthority
}

// GetSmBussinessTypeIDs 返回nil的时候，表示所有，返回[]时表示完全没有权限
func (admin *Admin) GetSmBussinessTypeIDs() []uint {
	if admin.parsedAuthority.Scope == OP_SCOPE_OPERATOR_BUSINESSTYPE {
		return admin.parsedAuthority.BusinessTypesResources
	} else if admin.parsedAuthority.Scope&LEVEL_SCOPE_ALL == LEVEL_SCOPE_ALL {
		return nil
	} else {
		return []uint{}
	}
}

// GetSmShopIDs 获取权限下的所有商铺,返回nil的时候，表示所有，返回[]时表示完全没有权限
func (admin *Admin) GetSmShopIDs(tx *gorm.DB) []uint {
	fmt.Println("scope", admin.parsedAuthority.Scope)
	if admin.parsedAuthority.Scope&LEVEL_SCOPE_ALL == LEVEL_SCOPE_ALL {
		return nil
	} else if admin.parsedAuthority.Scope == OP_SCOPE_OPERATOR_BUSINESSTYPE {
		businessTypeIDs := admin.parsedAuthority.BusinessTypesResources
		var shops []SmShop
		var shopIDs []uint
		tx.Where("business_type_id in (?)", businessTypeIDs).Find(&shops).Pluck("id", &shopIDs)
		if len(shops) == 0 {
			return []uint{}
		}
		return shopIDs
	} else if len(admin.parsedAuthority.ShopsResources) > 0 {
		return admin.parsedAuthority.ShopsResources
	} else {
		return []uint{}
	}
}

// GetSmRegionIDs 获取权限下的所有区域 返回nil的时候，表示所有，返回[]时表示完全没有权限
func (admin *Admin) GetSmRegionIDs(tx *gorm.DB) []uint {
	shopIDs := admin.GetSmShopIDs(tx)
	if shopIDs == nil {
		return nil
	} else if len(shopIDs) == 0 {
		return []uint{}
	} else {
		regionIDs := make([]uint, 0)
		var shops []SmShop
		var shopuvgroupIDs []uint
		tx.Where("id in (?)", shopIDs).Find(&shops).Pluck("entrance_uv_group_id", &shopuvgroupIDs)

		var uvGroups []SmUvGroup
		tx.Preload("SmRegions").Where("id in (?)", shopuvgroupIDs).Find(&uvGroups)

		for _, group := range uvGroups {
			for _, region := range group.SmRegions {
				regionIDs = append(regionIDs, region.ID)
			}
		}

		return regionIDs
	}
}
