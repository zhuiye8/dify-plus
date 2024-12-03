package gaia

import (
	"github.com/gofrs/uuid/v5"
	"time"
)

type Apps struct {
	ID                  uuid.UUID `gorm:"column:id;type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`                        // 应用ID
	TenantID            uuid.UUID `gorm:"column:tenant_id;type:uuid;not null" json:"tenant_id"`                                        // 工作空间ID
	Name                string    `gorm:"column:name;type:varchar(255);not null" json:"name"`                                          // 应用名称
	Mode                string    `gorm:"column:mode;type:varchar(255);not null" json:"mode"`                                          // 应用模式
	Icon                string    `gorm:"column:icon;type:varchar(255)" json:"icon"`                                                   // 应用图标
	IconBackground      string    `gorm:"column:icon_background;type:varchar(255)" json:"icon_background"`                             // 应用图标背景
	AppModelConfigID    uuid.UUID `gorm:"column:app_model_config_id;type:uuid" json:"app_model_config_id"`                             // 应用模型配置ID
	Status              string    `gorm:"column:status;type:varchar(255);not null;default:normal" json:"status"`                       // 应用状态
	EnableSite          bool      `gorm:"column:enable_site;not null" json:"enable_site"`                                              // 是否启用站点
	EnableAPI           bool      `gorm:"column:enable_api;not null" json:"enable_api"`                                                // 是否启用API
	APIRPM              int       `gorm:"column:api_rpm;type:int4;not null;default:0" json:"api_rpm"`                                  // API每分钟请求限制
	APIRPH              int       `gorm:"column:api_rph;type:int4;not null;default:0" json:"api_rph"`                                  // API每小时请求限制
	IsDemo              bool      `gorm:"column:is_demo;not null;default:false" json:"is_demo"`                                        // 是否为演示应用
	IsPublic            bool      `gorm:"column:is_public;not null;default:false" json:"is_public"`                                    // 是否为公开应用
	CreatedAt           time.Time `gorm:"column:created_at;type:timestamp(6);not null;default:CURRENT_TIMESTAMP(0)" json:"created_at"` // 创建时间
	UpdatedAt           time.Time `gorm:"column:updated_at;type:timestamp(6);not null;default:CURRENT_TIMESTAMP(0)" json:"updated_at"` // 更新时间
	IsUniversal         bool      `gorm:"column:is_universal;not null;default:false" json:"is_universal"`                              // 是否为通用应用
	WorkflowID          uuid.UUID `gorm:"column:workflow_id;type:uuid" json:"workflow_id"`                                             // 工作流ID
	Description         string    `gorm:"column:description;type:text;not null;default:''" json:"description"`                         // 应用描述
	Tracing             string    `gorm:"column:tracing;type:text" json:"tracing"`                                                     // 追踪信息
	MaxActiveRequests   int       `gorm:"column:max_active_requests;type:int4" json:"max_active_requests"`                             // 最大活跃请求数
	IconType            string    `gorm:"column:icon_type;type:varchar(255)" json:"icon_type"`                                         // 图标类型
	CreatedBy           uuid.UUID `gorm:"column:created_by;type:uuid" json:"created_by"`                                               // 创建者ID
	UpdatedBy           uuid.UUID `gorm:"column:updated_by;type:uuid" json:"updated_by"`                                               // 更新者ID
	UseIconAsAnswerIcon bool      `gorm:"column:use_icon_as_answer_icon;not null;default:false" json:"use_icon_as_answer_icon"`        // 是否使用图标作为回答图标
}

func (Apps) TableName() string {
	return "apps"
}

type AppStatisticsExtend struct {
	ID     uuid.UUID `gorm:"column:id;type:uuid;primary_key;default:uuid_generate_v4()" json:"id"` // 唯一标识符
	AppID  uuid.UUID `gorm:"column:app_id;type:uuid;not null" json:"app_id"`                       // 应用ID
	Number int       `gorm:"column:number;type:integer;not null" json:"number"`                    // 统计数量
}

func (AppStatisticsExtend) TableName() string {
	return "app_statistics_extend"
}
