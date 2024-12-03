package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type DashboardSearch struct {
	request.PageInfo
}

type GetAccountQuotaRankingDataReq struct {
	request.PageInfo
}

// GetAppQuotaRankingDataReq 获取应用配额排名数据
type GetAppQuotaRankingDataReq struct {
	request.PageInfo
}

// GetAppTokenQuotaRankingDataReq 获取应用配额排名数据
type GetAppTokenQuotaRankingDataReq struct {
	request.PageInfo
}

type GetAppTokenDailyQuotaDataReq struct {
	request.PageInfo
	AppId  string    `json:"app_id" form:"app_id"`   // 应用ID
	StatAt time.Time `json:"stat_at" form:"stat_at"` // 统计时间
}

// GetAiImageQuotaRankingDataReq 获取AI图片使用量排名数据
type GetAiImageQuotaRankingDataReq struct {
	request.PageInfo
	StatAt time.Time `json:"stat_at" form:"stat_at"` // 统计时间
}
