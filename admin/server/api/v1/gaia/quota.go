package gaia

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	gaiaReq "github.com/flipped-aurora/gin-vue-admin/server/model/gaia/request"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
	"go.uber.org/zap"
)

type QuotaApi struct{}

// QuotaManagementList
// @Tags Quota
// @Summary 额度管理列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query gaiaReq.DashboardSearch true "分页获取账号额度排名列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /gaia/quota/getManagementList [get]
func (quotaApi *QuotaApi) QuotaManagementList(c *gin.Context) {
	var pageInfo gaiaReq.GetAccountQuotaRankingDataReq
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := QuotaService.GetQuotaManagementData(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// SetUserQuota
// @Tags Quota
// @Summary 设置用户额度
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query gaiaReq.DashboardSearch true "分页获取账号额度排名列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /gaia/quota/setUserQuota [post]
func (quotaApi *QuotaApi) SetUserQuota(c *gin.Context) {
	var err error
	var uid uuid.UUID
	var pageInfo gaiaReq.SetUserQuotaRequest
	if err = c.ShouldBindJSON(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	//
	if uid, err = uuid.FromString(pageInfo.Uid); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err = QuotaService.SetUserQuota(uid, pageInfo.Quota); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed("ok", "修改成功", c)
}
