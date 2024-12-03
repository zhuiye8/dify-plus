// 自动生成模板Providers
package gaia

import (
	"github.com/gofrs/uuid/v5"
	"time"
)

// providers表 结构体  Providers
type Providers struct {
	Id              string     `json:"id" form:"id" gorm:"primarykey;column:id;comment:;"`                                    //id字段
	TenantId        string     `json:"tenantId" form:"tenantId" gorm:"index;column:tenant_id;comment:;"`                      //tenantId字段
	ProviderName    string     `json:"providerName" form:"providerName" gorm:"index;column:provider_name;comment:;size:255;"` //providerName字段
	ProviderType    string     `json:"providerType" form:"providerType" gorm:"column:provider_type;comment:;size:40;"`        //providerType字段
	EncryptedConfig string     `json:"encryptedConfig" form:"encryptedConfig" gorm:"column:encrypted_config;comment:;"`       //encryptedConfig字段
	IsValid         bool       `json:"isValid" form:"isValid" gorm:"column:is_valid;comment:;"`                               //isValid字段
	LastUsed        *time.Time `json:"lastUsed" form:"lastUsed" gorm:"column:last_used;comment:;size:6;"`                     //lastUsed字段
	QuotaType       string     `json:"quotaType" form:"quotaType" gorm:"column:quota_type;comment:;size:40;"`                 //quotaType字段
	QuotaLimit      int        `json:"quotaLimit" form:"quotaLimit" gorm:"column:quota_limit;comment:;size:64;"`              //quotaLimit字段
	QuotaUsed       int        `json:"quotaUsed" form:"quotaUsed" gorm:"column:quota_used;comment:;size:64;"`                 //quotaUsed字段
	CreatedAt       *time.Time `json:"createdAt" form:"createdAt" gorm:"column:created_at;comment:;size:6;"`                  //createdAt字段
	UpdatedAt       *time.Time `json:"updatedAt" form:"updatedAt" gorm:"column:updated_at;comment:;size:6;"`                  //updatedAt字段
	ModelName       string     `gorm:"column:model_name;type:varchar(255);not null" json:"model_name"`                        // 模型名称
}

// TableName providers表 Providers自定义表名 providers
func (Providers) TableName() string {
	return "providers"
}

type ProviderModels struct {
	ID              string    `gorm:"column:id;type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`                     // 唯一标识符
	TenantID        string    `gorm:"column:tenant_id;type:uuid;not null" json:"tenant_id"`                                     // 工作空间ID
	ProviderName    string    `gorm:"column:provider_name;type:varchar(255);not null" json:"provider_name"`                     // 提供商名称
	ModelName       string    `gorm:"column:model_name;type:varchar(255);not null" json:"model_name"`                           // 模型名称
	ModelType       string    `gorm:"column:model_type;type:varchar(40);not null" json:"model_type"`                            // 模型类型
	EncryptedConfig string    `gorm:"column:encrypted_config;type:text" json:"encrypted_config"`                                // 加密配置
	IsValid         bool      `gorm:"column:is_valid;type:bool;not null;default:false" json:"is_valid"`                         // 是否有效
	CreatedAt       time.Time `gorm:"column:created_at;type:timestamp(6);not null;default:CURRENT_TIMESTAMP" json:"created_at"` // 创建时间
	UpdatedAt       time.Time `gorm:"column:updated_at;type:timestamp(6);not null;default:CURRENT_TIMESTAMP" json:"updated_at"` // 更新时间
}

func (ProviderModels) TableName() string {
	return "provider_models"
}

type ProviderModelSettings struct {
	ID                   uuid.UUID `gorm:"column:id;primary_key;type:uuid;default:uuid_generate_v4()" json:"id"`                   // 唯一标识符
	TenantID             uuid.UUID `gorm:"column:tenant_id;type:uuid;not null" json:"tenant_id"`                                   // 工作空间ID
	ProviderName         string    `gorm:"column:provider_name;type:varchar(255);not null" json:"provider_name"`                   // 提供商名称
	ModelName            string    `gorm:"column:model_name;type:varchar(255);not null" json:"model_name"`                         // 模型名称
	ModelType            string    `gorm:"column:model_type;type:varchar(40);not null" json:"model_type"`                          // 模型类型
	Enabled              bool      `gorm:"column:enabled;type:boolean;default:true" json:"enabled"`                                // 是否启用
	LoadBalancingEnabled bool      `gorm:"column:load_balancing_enabled;type:boolean;default:false" json:"load_balancing_enabled"` // 是否启用负载均衡
	CreatedAt            time.Time `gorm:"column:created_at;type:timestamp(6);default:CURRENT_TIMESTAMP(0)" json:"created_at"`     // 创建时间
	UpdatedAt            time.Time `gorm:"column:updated_at;type:timestamp(6);default:CURRENT_TIMESTAMP(0)" json:"updated_at"`     // 更新时间
}

func (ProviderModelSettings) TableName() string {
	return "provider_model_settings"
}
