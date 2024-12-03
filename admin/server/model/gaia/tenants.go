// 自动生成模板Tenants
package gaia

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/gofrs/uuid/v5"
	"time"
)

// tenants表 结构体  Tenants
type Tenants struct {
	Id               string    `gorm:"column:id;type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`                     // 工作空间唯一标识
	Name             string    `gorm:"column:name;type:varchar(255);not null" json:"name"`                                       // 工作空间名称
	EncryptPublicKey string    `gorm:"column:encrypt_public_key;type:text" json:"encrypt_public_key"`                            // 加密公钥
	Plan             string    `gorm:"column:plan;type:varchar(255);not null;default:basic" json:"plan"`                         // 套餐类型
	Status           string    `gorm:"column:status;type:varchar(255);not null;default:normal" json:"status"`                    // 工作空间状态
	CreatedAt        time.Time `gorm:"column:created_at;type:timestamp(6);not null;default:CURRENT_TIMESTAMP" json:"created_at"` // 创建时间
	UpdatedAt        time.Time `gorm:"column:updated_at;type:timestamp(6);not null;default:CURRENT_TIMESTAMP" json:"updated_at"` // 更新时间
	CustomConfig     string    `gorm:"column:custom_config;type:text" json:"custom_config"`                                      // 自定义配置
}

func (*Tenants) TableName() string {
	return "tenants"
}

func (t *Tenants) GetSuperAdminTenantId() string {
	err := global.GVA_DB.Order("created_at ASC").First(&t).Error
	if err != nil {
		global.GVA_LOG.Error("GetSuperAdminTenantId gaia表查询失败,原因: " + err.Error())
		return ""
	}
	return t.Id
}

type TenantAccountJoins struct {
	ID        uuid.UUID `gorm:"column:id;type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`                     // 唯一标识符
	TenantID  uuid.UUID `gorm:"column:tenant_id;type:uuid;not null" json:"tenant_id"`                                     // 工作空间ID
	AccountID uuid.UUID `gorm:"column:account_id;type:uuid;not null" json:"account_id"`                                   // 账户ID
	Role      string    `gorm:"column:role;type:varchar(16);not null;default:normal" json:"role"`                         // 角色
	InvitedBy uuid.UUID `gorm:"column:invited_by;type:uuid" json:"invited_by"`                                            // 邀请人ID
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp(6);not null;default:CURRENT_TIMESTAMP" json:"created_at"` // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp(6);not null;default:CURRENT_TIMESTAMP" json:"updated_at"` // 更新时间
	Current   bool      `gorm:"column:current;type:bool;not null;default:false" json:"current"`                           // 是否当前
}

func (TenantAccountJoins) TableName() string {
	return "tenant_account_joins"
}
