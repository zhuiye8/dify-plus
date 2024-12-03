package gaia

import (
	"github.com/gin-gonic/gin"
)

type TenantsRouter struct{}

// InitTenantsRouter 初始化 tenants表 路由信息
func (s *TenantsRouter) InitTenantsRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	tenantsRouterWithoutRecord := Router.Group("tenants")
	{
		tenantsRouterWithoutRecord.GET("findTenants", tenantsApi.FindTenants)       // 根据ID获取tenants表
		tenantsRouterWithoutRecord.GET("getTenantsList", tenantsApi.GetTenantsList) // 获取tenants表列表
		tenantsRouterWithoutRecord.GET("getAllTenants", tenantsApi.GetAllTenants)   // 获取所有工作区
	}
}
