package gaia

import (
	"github.com/gin-gonic/gin"
)

type TestRouter struct{}

// InitTestRouter 初始化 测试表 路由信息
func (d *TestRouter) InitTestRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	dashboardRouterWithoutRecord := Router.Group("gaia/test")
	{
		dashboardRouterWithoutRecord.POST("app/request", testApi.GaiaAppRequestTest)           // 发起gaia应用请求测试
		dashboardRouterWithoutRecord.GET("app/request/list", testApi.GaiaAppRequestTestList)   // gaia应用请求测试结果列表
		dashboardRouterWithoutRecord.GET("app/request/batch", testApi.GaiaAppRequestTestBatch) // gaia应用请求测试批次列表
		dashboardRouterWithoutRecord.POST("sync/database", testApi.SyncDatabaseTableData)      // 同步数据库表数据
	}
}
