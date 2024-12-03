package gaia

import "time"

const IndirectAccessUser = "end_user"
const UserUsingApiRequest = "service_api"
const UsernameUsingApiRequest = "gaia_test_api_user"

type EndUser struct {
	ID             string    `json:"id" gorm:"primary_key;comment:用户唯一标识"`
	TenantID       string    `json:"tenant_id" gorm:"comment:租户唯一标识"`
	AppID          string    `json:"app_id" gorm:"comment:应用唯一标识"` // 使用指针允许该字段为NULL
	Type           string    `json:"type" gorm:"comment:用户类型"`
	ExternalUserID string    `json:"external_user_id" gorm:"comment:外部用户唯一标识"` // 使用指针允许该字段为NULL
	Name           string    `json:"name" gorm:"comment:用户名称"`                 // 使用指针允许该字段为NULL
	IsAnonymous    bool      `json:"is_anonymous" gorm:"comment:是否匿名"`
	SessionID      string    `json:"session_id" gorm:"comment:会话唯一标识"`
	CreatedAt      time.Time `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"comment:更新时间"`
}

func (EndUser) TableName() string { return "end_users" }
