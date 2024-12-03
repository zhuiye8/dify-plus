package gaia

import (
	"github.com/gin-gonic/gin"
)

type QuotaRouter struct{}

// InitQuotaRouter 初始化 quota表 路由信息
func (d *QuotaRouter) InitQuotaRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	dashboardRouterWithoutRecord := Router.Group("gaia/quota")
	{
		dashboardRouterWithoutRecord.POST("setUserQuota", quotaApi.SetUserQuota)            // 设置用户额度
		dashboardRouterWithoutRecord.GET("getManagementList", quotaApi.QuotaManagementList) // 额度管理列表
	}
}
