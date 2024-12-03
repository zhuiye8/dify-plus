// 自动生成模板AccountMoneyExtend
package gaia

import (
	"github.com/gofrs/uuid/v5"
	"time"
)

// accountMoneyExtend表 结构体  AccountMoneyExtend
type AccountMoneyExtend struct {
	Id         *string    `json:"id" form:"id" gorm:"primarykey;column:id;comment:;"`                        //id字段
	AccountId  uuid.UUID  `json:"accountId" form:"accountId" gorm:"uniqueIndex;column:account_id;comment:;"` //accountId字段
	TotalQuota float64    `json:"totalQuota" form:"totalQuota" gorm:"column:total_quota;comment:;"`          //totalQuota字段
	UsedQuota  float64    `json:"usedQuota" form:"usedQuota" gorm:"column:used_quota;comment:;"`             //usedQuota字段
	CreatedAt  *time.Time `json:"createdAt" form:"createdAt" gorm:"column:created_at;comment:;size:6;"`      //createdAt字段
	UpdatedAt  *time.Time `json:"updatedAt" form:"updatedAt" gorm:"column:updated_at;comment:;size:6;"`      //updatedAt字段
}

// TableName accountMoneyExtend表 AccountMoneyExtend自定义表名 account_money_extend
func (AccountMoneyExtend) TableName() string {
	return "account_money_extend"
}
