package gaia

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	gaiaReq "github.com/flipped-aurora/gin-vue-admin/server/model/gaia/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TenantsApi struct{}

// FindTenants 用id查询tenants表
// @Tags Tenants
// @Summary 用id查询tenants表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query gaia.Tenants true "用id查询tenants表"
// @Success 200 {object} response.Response{data=gaia.Tenants,msg=string} "查询成功"
// @Router /tenants/findTenants [get]
func (tenantsApi *TenantsApi) FindTenants(c *gin.Context) {
	id := c.Query("id")
	retenants, err := tenantsService.GetTenants(id)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithData(retenants, c)
}

// GetTenantsList 分页获取tenants表列表
// @Tags Tenants
// @Summary 分页获取tenants表列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query gaiaReq.TenantsSearch true "分页获取tenants表列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /tenants/getTenantsList [get]
func (tenantsApi *TenantsApi) GetTenantsList(c *gin.Context) {
	var pageInfo gaiaReq.TenantsSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := tenantsService.GetTenantsInfoList(pageInfo)
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

// GetAllTenants 获取所有工作区
// @Tags Tenants
// @Summary 获取所有工作区
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query gaia.Tenants true "获取所有工作区"
// @Success 200 {object} response.Response{data=gaia.Tenants,msg=string} "查询成功"
// @Router /tenants/getAllTenants [get]
func (tenantsApi *TenantsApi) GetAllTenants(c *gin.Context) {
	retenants, err := tenantsService.GetAllTenants()
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithData(retenants, c)
}
