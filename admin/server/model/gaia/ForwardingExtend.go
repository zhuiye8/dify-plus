package gaia

import (
	"github.com/goccy/go-json"
	"github.com/gofrs/uuid/v5"
	"github.com/richardlehane/msoleps/types"
	"time"
)

type ForwardingExtend struct {
	ID          uuid.UUID `gorm:"column:id;type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`                   // 唯一标识符
	Path        string    `gorm:"column:path;type:varchar(255);not null" json:"path"`                                     // 转发路径
	Address     string    `gorm:"column:address;type:varchar(255);not null" json:"address"`                               // 转发地址
	Header      string    `gorm:"column:header;type:text;not null;default:'[]'::text" json:"header"`                      // 请求头信息
	Description string    `gorm:"column:description;type:text;not null;default:''::character varying" json:"description"` // 描述信息
}

func (ForwardingExtend) TableName() string {
	return "forwarding_extend"
}

type ForwardingAddressExtend struct {
	ID           uuid.UUID `gorm:"column:id;primary_key;type:uuid;default:uuid_generate_v4()" json:"id"`    // 唯一标识符
	ForwardingID uuid.UUID `gorm:"column:forwarding_id;type:uuid;not null" json:"forwarding_id"`            // 转发ID
	Path         string    `gorm:"column:path;type:varchar(255);not null" json:"path"`                      // 路径
	Models       string    `gorm:"column:models;type:varchar(255);not null" json:"models"`                  // 模型
	Description  string    `gorm:"column:description;type:text;default:''" json:"description"`              // 描述
	ContentType  int       `gorm:"column:content_type;type:integer;not null;default:0" json:"content_type"` // 内容类型
	Billing      string    `gorm:"column:billing;type:text;default:'[]'" json:"billing"`                    // 计费信息
	Status       bool      `gorm:"column:status;type:boolean;default:true" json:"status"`                   // 状态
}

func (ForwardingAddressExtend) TableName() string {
	return "public.forwarding_address_extend"
}

type AccountLayoverRecordExtend struct {
	ID           uuid.UUID       `gorm:"column:id;type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`                     // 记录ID
	AccountID    uuid.UUID       `gorm:"column:account_id;type:uuid;not null" json:"account_id"`                                   // 账户ID
	ForwardingID uuid.UUID       `gorm:"column:forwarding_id;type:uuid;not null" json:"forwarding_id"`                             // 转发ID
	Money        types.Decimal   `gorm:"column:money;type:numeric(16,7)" json:"money"`                                             // 金额
	Info         json.RawMessage `gorm:"column:info;type:json" json:"info"`                                                        // 附加信息
	CreatedAt    time.Time       `gorm:"column:created_at;type:timestamp(6);not null;default:CURRENT_TIMESTAMP" json:"created_at"` // 创建时间
}

func (AccountLayoverRecordExtend) TableName() string {
	return "account_layover_record_extend"
}
