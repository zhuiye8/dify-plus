package gaia

import (
	"errors"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/gofrs/uuid/v5"
	"time"
)

const UserActive = "active"               // 用户状态: 活跃
const UserPending = "pending"             // 用户状态: 待办的
const UserUninitialized = "uninitialized" // 用户状态: 未初始化
const UserBanned = "banned"               // 用户状态: 禁止
const UserClosed = "closed"               // 用户状态: 关闭
const DefaultProviderType = "oauth2"      // 默认提供者类型: oauth2

// Account gaia 用户表
type Account struct {
	ID                uuid.UUID `json:"id" gorm:"primaryKey;comment:账户唯一标识符"`
	Name              string    `json:"name" gorm:"not null;comment:账户名称"`
	Email             string    `json:"email" gorm:"not null;index:account_email_idx;comment:账户邮箱"`
	Password          string    `json:"password" gorm:"comment:账户密码"`
	PasswordSalt      string    `json:"password_salt" gorm:"comment:加密密码的盐值"`
	Avatar            string    `json:"avatar" gorm:"comment:头像URL"`
	InterfaceLanguage string    `json:"interface_language" gorm:"comment:用户界面语言"`
	InterfaceTheme    string    `json:"interface_theme" gorm:"comment:用户界面主题"`
	Timezone          string    `json:"timezone" gorm:"comment:用户时区"`
	LastLoginAt       time.Time `json:"last_login_at" gorm:"comment:最后登录时间"`
	LastLoginIP       string    `json:"last_login_ip" gorm:"comment:最后登录的IP地址"`
	Status            string    `json:"status" gorm:"default:'active';not null;comment:账户状态"`
	InitializedAt     time.Time `json:"initialized_at" gorm:"comment:账户初始化时间"`
	CreatedAt         time.Time `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP;comment:账户创建时间"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP;comment:账户更新时间"`
	LastActiveAt      time.Time `json:"last_active_at" gorm:"not null;default:CURRENT_TIMESTAMP;comment:最后活跃时间"`
}

// AccountIntegrate gaia 用户提供者关联表
type AccountIntegrate struct {
	ID             uuid.UUID `json:"id" gorm:"index;comment:唯一标识"`
	AccountID      uuid.UUID `json:"account_id" gorm:"not null;comment:账户ID"`
	Provider       string    `json:"provider" gorm:"not null;comment:提供者类型"`
	OpenID         string    `json:"open_id" gorm:"not null;comment:开放ID"`
	EncryptedToken string    `json:"encrypted_token" gorm:"not null;comment:加密令牌"`
	CreatedAt      time.Time `json:"created_at" gorm:"not null;comment:创建时间"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"not null;comment:更新时间"`
}

// TenantAccountJoin gaia 用户和命名空间关联表
type TenantAccountJoin struct {
	ID        string    `json:"id" gorm:"primary_key;comment:租户账户连接ID"`
	TenantID  string    `json:"tenant_id" gorm:"not null;comment:租户ID"`
	AccountID string    `json:"account_id" gorm:"not null;comment:账户ID"`
	Role      string    `json:"role" gorm:"not null;default:normal;comment:角色"`
	InvitedBy *string   `json:"invited_by" gorm:"comment:邀请人ID"`
	CreatedAt time.Time `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP;comment:创建时间"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP;comment:更新时间"`
	Current   bool      `json:"current" gorm:"not null;default:false;comment:当前状态"`
}

// AccountDingTalkExtend gaia钉钉关联表
type AccountDingTalkExtend struct {
	ID       uuid.UUID `json:"id" gorm:"primaryKey;comment:账户唯一标识符"`
	DingTalk string    `json:"ding_talk" gorm:"index:account_ding_talk_idx;comment:关联钉钉id"`
}

func (Account) TableName() string               { return "accounts" }
func (AccountIntegrate) TableName() string      { return "account_integrates" }
func (TenantAccountJoin) TableName() string     { return "tenant_account_joins" }
func (AccountDingTalkExtend) TableName() string { return "account_ding_talk_extend" }

// GetAccount
// @description: Get user information through the user provider relationship table
// @return account, err error
func (i Account) GetAccount(username string) (integrate AccountIntegrate, err error) {
	// init
	if err = global.GVA_DB.Where("account_id IN (?) AND provider=?",
		[]string{i.ID.String(), username}, DefaultProviderType).First(&integrate).Error; err != nil {
		return integrate, errors.New("the query for the user-provider relationship could not be found")
	}
	// return
	return integrate, nil
}
