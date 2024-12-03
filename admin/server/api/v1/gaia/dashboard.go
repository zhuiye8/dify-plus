package gaia

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	gaiaReq "github.com/flipped-aurora/gin-vue-admin/server/model/gaia/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DashboardApi struct{}

// GetAccountQuotaRankingData 分页获取账号额度排名
// @Tags Dashboard
// @Summary 分页获取dashboard表列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query gaiaReq.DashboardSearch true "分页获取账号额度排名列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /dashboard/getAccountQuotaRankingData [get]
func (dashboardApi *DashboardApi) GetAccountQuotaRankingData(c *gin.Context) {
	var pageInfo gaiaReq.GetAccountQuotaRankingDataReq
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := dashboardService.GetAccountQuotaRankingData(pageInfo)
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

// GetAppQuotaRankingData 分页获取【应用】配额排名数据
// @Tags Dashboard
// @Summary 分页获取dashboard表列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query gaiaReq.DashboardSearch true "分页获取【应用】配额排名数据"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /dashboard/getAppQuotaRankingData [get]
func (dashboardApi *DashboardApi) GetAppQuotaRankingData(c *gin.Context) {
	var pageInfo gaiaReq.GetAppQuotaRankingDataReq
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := dashboardService.GetAppQuotaRankingData(pageInfo)
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

// GetAppTokenQuotaRankingData 分页获取【应用密钥】配额排名数据列表
// @Tags Dashboard
// @Summary 分页获取dashboard表列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query gaiaReq.DashboardSearch true "分页获取【应用密钥】配额排名数据列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /dashboard/getAppTokenQuotaRankingData [get]
func (dashboardApi *DashboardApi) GetAppTokenQuotaRankingData(c *gin.Context) {
	var pageInfo gaiaReq.GetAppTokenQuotaRankingDataReq
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := dashboardService.GetAppTokenQuotaRankingData(pageInfo)
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

// GetAppTokenDailyQuotaData 获取每天密钥花费数据列表
// @Tags Dashboard
// @Summary 分页获取dashboard表列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query gaiaReq.DashboardSearch true "获取每天密钥花费数据列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /dashboard/getAppTokenDailyQuotaData [get]
func (dashboardApi *DashboardApi) GetAppTokenDailyQuotaData(c *gin.Context) {
	var pageInfo gaiaReq.GetAppTokenDailyQuotaDataReq
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, err := dashboardService.GetAppTokenDailyQuotaData(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// GetAiImageQuotaRankingData 获取每天ai图片花费数据列表
// @Tags Dashboard
// @Summary 获取每天ai图片花费数据列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query gaiaReq.DashboardSearch true "获取每天ai图片花费数据列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /dashboard/getAppTokenDailyQuotaData [get]
func (dashboardApi *DashboardApi) GetAiImageQuotaRankingData(c *gin.Context) {
	var pageInfo gaiaReq.GetAiImageQuotaRankingDataReq
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, err := dashboardService.GetAiImageQuotaRankingData(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}
