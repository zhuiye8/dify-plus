package gaia

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/gaia"
	gaiaReq "github.com/flipped-aurora/gin-vue-admin/server/model/gaia/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/gaia/response"
	"github.com/gofrs/uuid/v5"
	"go.uber.org/zap"
	"strings"
)

type QuotaService struct{}

// GetQuotaManagementData
// @Tags Quota
// @Summary 额度管理列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param info gaiaReq.GetAccountQuotaRankingDataReq
// @Return list []response.GetQuotaManagementDataResponse, total int64, err error
func (dashboardService *QuotaService) GetQuotaManagementData(info gaiaReq.GetAccountQuotaRankingDataReq) (
	list []response.GetQuotaManagementDataResponse, total int64, err error) {
	if info.PageSize == 0 {
		info.PageSize = 10
	}
	var uuidList []string
	limit := info.PageSize
	var accountList []gaia.Account
	s := strings.TrimSpace(info.Keyword)
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&gaia.AccountMoneyExtend{}).Order("used_quota desc")
	if len(s) > 0 {
		s = fmt.Sprintf("%%%s%%", s)
		if err = global.GVA_DB.Debug().Select("id").Where(
			"\"name\" LIKE ? OR \"email\" LIKE ?", s, s).Find(&accountList).Error; err == nil {
			for _, v := range accountList {
				uuidList = append(uuidList, v.ID.String())
			}
		}
		// len
		if len(uuidList) > 0 {
			db.Where("account_id IN (?)", uuidList)
		} else {
			db.Where("account_id is null")
		}
	}

	var accountMoneys []gaia.AccountMoneyExtend
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&accountMoneys).Error
	if err != nil {
		err = fmt.Errorf("查询账号额度信息失败：%s", err.Error())
		return
	}

	// 账号ID集合，方便后面一次性查出
	var accountIds []uuid.UUID
	for _, money := range accountMoneys {
		accountIds = append(accountIds, money.AccountId)
	}

	// 查询用户信息
	var accountInfos = make(map[uuid.UUID]gaia.Account)
	var accounts []gaia.Account
	err = global.GVA_DB.Model(&gaia.Account{}).Where("id in ?", accountIds).Find(&accounts).Error
	if err != nil {
		err = fmt.Errorf("查询账户信息失败：%s", err.Error())
		return
	}
	for _, account := range accounts {
		accountInfos[account.ID] = account
	}

	// 拼接结果
	for i, money := range accountMoneys {
		var accountInfo gaia.Account
		var isExist bool
		if accountInfo, isExist = accountInfos[money.AccountId]; !isExist {
			global.GVA_LOG.Error("账户信息不存在!", zap.String("account_id", money.AccountId.String()))
			continue
		}
		list = append(list, response.GetQuotaManagementDataResponse{
			Uid:        money.AccountId.String(),
			Ranking:    i + 1 + offset,
			Name:       accountInfo.Name,
			UsedQuota:  money.UsedQuota,
			TotalQuota: money.TotalQuota,
		})
	}

	return list, total, err
}

// SetUserQuota
// @Tags Quota
// @Summary 设置指定用户额度
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param info gaiaReq.SetUserQuotaRequest
// @Return err error
func (dashboardService *QuotaService) SetUserQuota(uid uuid.UUID, quota float64) error {

	return global.GVA_DB.Model(&gaia.AccountMoneyExtend{}).Where(
		"account_id = ?", uid).Updates(&map[string]interface{}{
		"total_quota": quota,
	}).Error
}
