package gaia

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type ProvidersRouter struct{}

// InitProvidersRouter 初始化 providers表 路由信息
func (s *ProvidersRouter) InitProvidersRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	providersRouter := Router.Group("providers").Use(middleware.OperationRecord())
	providersRouterWithoutRecord := Router.Group("providers")
	{
		providersRouter.PUT("syncProviders", providersApi.SyncProviders) // 同步模型
	}
	{
		providersRouterWithoutRecord.GET("findProviders", providersApi.FindProviders)       // 根据ID获取providers表
		providersRouterWithoutRecord.GET("getProvidersList", providersApi.GetProvidersList) // 获取providers表列表
	}
}
