package gaia

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	gaiaReq "github.com/flipped-aurora/gin-vue-admin/server/model/gaia/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ProvidersApi struct{}

// SyncProviders 同步模型
// @Tags Providers
// @Summary 同步模型
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body gaia.Providers true "同步模型"
// @Success 200 {object} response.Response{msg=string} "同步成功"
// @Router /providers/syncProviders [put]
func (providersApi *ProvidersApi) SyncProviders(c *gin.Context) {
	var req gaiaReq.SyncProviderReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	jwtToken := utils.GetToken(c)
	err = providersService.SyncProviders(jwtToken, req)
	if err != nil {
		global.GVA_LOG.Error("同步失败!", zap.Error(err))
		response.FailWithMessage("同步失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("同步成功", c)
}

// FindProviders 用id查询providers表
// @Tags Providers
// @Summary 用id查询providers表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query gaia.Providers true "用id查询providers表"
// @Success 200 {object} response.Response{data=gaia.Providers,msg=string} "查询成功"
// @Router /providers/findProviders [get]
func (providersApi *ProvidersApi) FindProviders(c *gin.Context) {
	id := c.Query("id")
	reproviders, err := providersService.GetProviders(id)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithData(reproviders, c)
}

// GetProvidersList 分页获取providers表列表
// @Tags Providers
// @Summary 分页获取providers表列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query gaiaReq.ProvidersSearch true "分页获取providers表列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /providers/getProvidersList [get]
func (providersApi *ProvidersApi) GetProvidersList(c *gin.Context) {
	var pageInfo gaiaReq.ProvidersSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := providersService.GetProvidersInfoList(pageInfo)
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
