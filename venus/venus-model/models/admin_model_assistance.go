package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type AdminForRoleSerializer struct {
	ID          uint                               `json:"id"`
	CreatedAt   time.Time                          `json:"created_at"`
	UpdatedAt   time.Time                          `json:"updated_at"`
	Name        string                             `json:"name"`
	Role        string                             `json:"role"`
	AuthToken   string                             `json:"auth_token"`
	Phone       string                             `json:"phone"`
	CompanyID   uint                               `json:"company_id"`
	State       string                             `json:"state"`
	AuthorityID uint                               `json:"authority_id"`
	Authority   ExtendedAuthoritySerializerForRole `json:"authority"`
}

func (admin Admin) SerializerForRole(tx *gorm.DB) AdminForRoleSerializer {
	return AdminForRoleSerializer{
		ID:          admin.ID,
		CreatedAt:   admin.CreatedAt.Round(time.Second),
		UpdatedAt:   admin.UpdatedAt.Round(time.Second),
		Name:        admin.Name,
		Role:        admin.Role,
		AuthToken:   admin.AuthToken,
		Phone:       admin.Phone,
		CompanyID:   admin.CompanyID,
		State:       admin.State,
		AuthorityID: admin.AuthorityID,
		Authority:   admin.Authority.SerializeForRole(tx, admin.Company.Name),
	}
}
