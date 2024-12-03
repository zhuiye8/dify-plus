package gaia

import (
	"github.com/gin-gonic/gin"
)

type DashboardRouter struct{}

// InitDashboardRouter 初始化 dashboard表 路由信息
func (d *DashboardRouter) InitDashboardRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	dashboardRouterWithoutRecord := Router.Group("gaia/dashboard")
	{
		dashboardRouterWithoutRecord.GET("getAccountQuotaRankingData", dashboardApi.GetAccountQuotaRankingData)   // 分页获取账号额度排名
		dashboardRouterWithoutRecord.GET("getAppQuotaRankingData", dashboardApi.GetAppQuotaRankingData)           // 分页获取【应用】配额排名数据
		dashboardRouterWithoutRecord.GET("getAppTokenQuotaRankingData", dashboardApi.GetAppTokenQuotaRankingData) // 分页获取【应用密钥】配额排名数据列表
		dashboardRouterWithoutRecord.GET("getAppTokenDailyQuotaData", dashboardApi.GetAppTokenDailyQuotaData)     // 获取每天密钥花费数据列表
		dashboardRouterWithoutRecord.GET("getAiImageQuotaRankingData", dashboardApi.GetAiImageQuotaRankingData)   // 获取每天ai图片额度排名
	}
}
