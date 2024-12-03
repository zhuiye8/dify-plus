package gaia

import (
	"github.com/gofrs/uuid/v5"
	"time"
)

// ModelSyncConfigExtend 模型同步配置表
type ModelSyncConfigExtend struct {
	ID        uuid.UUID `gorm:"column:id;type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`                        // 唯一标识
	ModelID   string    `gorm:"column:model_id;type:uuid;uniqueIndex:unique_model_id" json:"model_id"`                       // 关联的模型ID
	IsAll     bool      `gorm:"column:is_all;default:true" json:"is_all"`                                                    // 是否同步所有
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp(6);not null;default:CURRENT_TIMESTAMP(0)" json:"created_at"` // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp(6);not null;default:CURRENT_TIMESTAMP(0)" json:"updated_at"` // 更新时间
}

func (ModelSyncConfigExtend) TableName() string {
	return "model_sync_config_extend"
}

// TenantModelSyncExtend 工作空间已同步模型关联表
type TenantModelSyncExtend struct {
	ID            uuid.UUID `gorm:"column:id;primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`                // 唯一标识符
	TenantID      string    `gorm:"column:tenant_id;type:uuid;not null" json:"tenant_id"`                               // 工作空间ID
	ModelID       uuid.UUID `gorm:"column:model_id;type:uuid;not null;index" json:"model_id"`                           // 模型ID
	OriginModelID string    `gorm:"column:origin_model_id;type:varchar(255);not null" json:"origin_model_id"`           // 原始模型ID
	IsAll         bool      `gorm:"column:is_all;type:boolean;default:false" json:"is_all"`                             // 是否全部
	CreatedAt     time.Time `gorm:"column:created_at;type:timestamp(6);default:CURRENT_TIMESTAMP(0)" json:"created_at"` // 创建时间
	UpdatedAt     time.Time `gorm:"column:updated_at;type:timestamp(6);default:CURRENT_TIMESTAMP(0)" json:"updated_at"` // 更新时间
}

func (TenantModelSyncExtend) TableName() string {
	return "tenant_model_sync_extend"
}
