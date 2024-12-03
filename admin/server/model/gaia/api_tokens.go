package gaia

import (
	"github.com/gofrs/uuid/v5"
	"time"
)

type ApiTokens struct {
	ID         uuid.UUID  `gorm:"column:id;type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`                     // 主键ID
	AppID      uuid.UUID  `gorm:"column:app_id;type:uuid" json:"app_id"`                                                    // 应用ID
	Type       string     `gorm:"column:type;type:varchar(16);not null" json:"type"`                                        // 令牌类型
	Token      string     `gorm:"column:token;type:varchar(255);not null" json:"token"`                                     // 令牌值
	LastUsedAt *time.Time `gorm:"column:last_used_at;type:timestamp(6)" json:"last_used_at"`                                // 最后使用时间
	CreatedAt  time.Time  `gorm:"column:created_at;type:timestamp(6);not null;default:CURRENT_TIMESTAMP" json:"created_at"` // 创建时间
	TenantID   uuid.UUID  `gorm:"column:tenant_id;type:uuid" json:"tenant_id"`                                              // 工作空间ID
}

func (*ApiTokens) TableName() string {
	return "api_tokens"
}
func (a *ApiTokens) GenerateToken() string {
	return a.Token[:3] + "..." + a.Token[len(a.Token)-23:]
}

type ApiTokenMoneyDailyStatExtend struct {
	ID               uuid.UUID `gorm:"column:id;type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`                     // 主键ID
	AppTokenID       uuid.UUID `gorm:"column:app_token_id;type:uuid;not null" json:"app_token_id"`                               // 应用令牌ID
	AccumulatedQuota float64   `gorm:"column:accumulated_quota;type:numeric(16,7)" json:"accumulated_quota"`                     // 累计配额
	DayUsedQuota     float64   `gorm:"column:day_used_quota;type:numeric(16,7)" json:"day_used_quota"`                           // 日使用配额
	DayLimitQuota    float64   `gorm:"column:day_limit_quota;type:numeric(16,7)" json:"day_limit_quota"`                         // 日限额配额
	StatAt           time.Time `gorm:"column:stat_at;type:timestamp(6);not null" json:"stat_at"`                                 // 统计时间
	CreatedAt        time.Time `gorm:"column:created_at;type:timestamp(6);not null;default:CURRENT_TIMESTAMP" json:"created_at"` // 创建时间
	UpdatedAt        time.Time `gorm:"column:updated_at;type:timestamp(6);not null;default:CURRENT_TIMESTAMP" json:"updated_at"` // 更新时间
}

func (*ApiTokenMoneyDailyStatExtend) TableName() string {
	return "api_token_money_daily_stat_extend"
}

type ApiTokenMoneyMonthlyStatExtend struct {
	ID               uuid.UUID `gorm:"column:id;type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`                     // 主键ID
	AppTokenID       uuid.UUID `gorm:"column:app_token_id;type:uuid;not null" json:"app_token_id"`                               // 应用令牌ID
	AccumulatedQuota float64   `gorm:"column:accumulated_quota;type:numeric(16,7)" json:"accumulated_quota"`                     // 累计配额
	MonthUsedQuota   float64   `gorm:"column:month_used_quota;type:numeric(16,7)" json:"month_used_quota"`                       // 月使用配额
	MonthLimitQuota  float64   `gorm:"column:month_limit_quota;type:numeric(16,7)" json:"month_limit_quota"`                     // 月限额配额
	StatAt           time.Time `gorm:"column:stat_at;type:timestamp(6);not null" json:"stat_at"`                                 // 统计时间
	CreatedAt        time.Time `gorm:"column:created_at;type:timestamp(6);not null;default:CURRENT_TIMESTAMP" json:"created_at"` // 创建时间
	UpdatedAt        time.Time `gorm:"column:updated_at;type:timestamp(6);not null;default:CURRENT_TIMESTAMP" json:"updated_at"` // 更新时间
}

func (*ApiTokenMoneyMonthlyStatExtend) TableName() string {
	return "api_token_money_monthly_stat_extend"
}

type ApiTokenMoneyExtend struct {
	ID               uuid.UUID `gorm:"column:id;type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`                     // 主键ID
	AppTokenID       uuid.UUID `gorm:"column:app_token_id;type:uuid" json:"app_token_id"`                                        // 应用令牌ID
	AccumulatedQuota float64   `gorm:"column:accumulated_quota;type:numeric(16,7)" json:"accumulated_quota"`                     // 累计配额
	DayUsedQuota     float64   `gorm:"column:day_used_quota;type:numeric(16,7)" json:"day_used_quota"`                           // 日使用配额
	MonthUsedQuota   float64   `gorm:"column:month_used_quota;type:numeric(16,7)" json:"month_used_quota"`                       // 月使用配额
	DayLimitQuota    float64   `gorm:"column:day_limit_quota;type:numeric(16,7)" json:"day_limit_quota"`                         // 日限额配额
	MonthLimitQuota  float64   `gorm:"column:month_limit_quota;type:numeric(16,7)" json:"month_limit_quota"`                     // 月限额配额
	IsDeleted        bool      `gorm:"column:is_deleted;type:bool;not null;default:false" json:"is_deleted"`                     // 是否删除
	UpdatedAt        time.Time `gorm:"column:updated_at;type:timestamp(6);not null;default:CURRENT_TIMESTAMP" json:"updated_at"` // 更新时间
	CreatedAt        time.Time `gorm:"column:created_at;type:timestamp(6);not null;default:CURRENT_TIMESTAMP" json:"created_at"` // 创建时间
	Description      string    `gorm:"column:description;type:varchar(50)" json:"description"`                                   // 描述
}

func (*ApiTokenMoneyExtend) TableName() string {
	return "api_token_money_extend"
}
