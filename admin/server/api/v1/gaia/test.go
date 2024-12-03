package gaia

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/gaia"
	"github.com/flipped-aurora/gin-vue-admin/server/model/gaia/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TestApi struct{}

// SyncDatabaseTableData
// @Tags Test
// @Summary 同步数据库表数据
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query gaiaReq.DashboardSearch true "分页获取账号额度排名列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /gaia/test/sync/database [post]
func (quotaApi *TestApi) SyncDatabaseTableData(c *gin.Context) {
	var data request.SyncDatabaseTableData
	if err := c.ShouldBindJSON(&data); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// 是否都不为空
	if len(data.LogTable) > 0 && len(data.NewTable) > 0 && len(data.KeyName) > 0 && len(data.OrderName) > 0 {
		go TestService.SyncDatabaseTableData(data.LogTable, data.NewTable, data.KeyName, data.OrderName, data.GroupName)
		response.OkWithDetailed("ok", "获取成功", c)
	} else {
		response.FailWithMessage("传参有误", c)
	}

}

// GaiaAppRequestTest
// @Tags Test
// @Summary 发起gaia应用请求测试
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query gaiaReq.DashboardSearch true "分页获取账号额度排名列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /gaia/test/app/request [post]
func (quotaApi *TestApi) GaiaAppRequestTest(c *gin.Context) {
	err := TestService.AppRequestTest()
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed("ok", "获取成功", c)
}

// GaiaAppRequestTestList
// @Tags Test
// @Summary gaia应用请求测试结果列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query gaiaReq.DashboardSearch true "分页获取账号额度排名列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /gaia/test/app/request/list [get]
func (quotaApi *TestApi) GaiaAppRequestTestList(c *gin.Context) {
	var pageInfo request.GetAppRequestTestRequest
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// list
	lock, list, total, err := TestService.AppRequestTestList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}
	var appIdList []string
	var appList []map[string]string
	var batchInfo []gaia.AppRequestTest
	if err = global.GVA_DB.Select("app_id").Where("batch_id = ?", pageInfo.BatchId).Group(
		"app_id").Find(&batchInfo).Error; err == nil {
		for _, v := range batchInfo {
			appIdList = append(appIdList, v.AppID)
		}
	}
	// 查询相关的app列表
	if len(appIdList) > 0 {
		var apps []gaia.Apps
		if err = global.GVA_DB.Select("id", "name").Where(
			"id IN (?)", appIdList).Find(&apps).Error; err == nil {
			for _, v := range apps {
				appList = append(appList, map[string]string{
					"value": v.ID.String(),
					"label": v.Name,
				})
			}
		}
	}
	response.OkWithDetailed(map[string]interface{}{
		"lock":      lock,
		"list":      list,
		"total":     total,
		"apps":      appList,
		"page":      pageInfo.Page,
		"page_size": pageInfo.PageSize,
	}, "获取成功", c)
}

// GaiaAppRequestTestBatch
// @Tags Test
// @Summary gaia应用请求测试批次列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query gaiaReq.DashboardSearch true "分页获取账号额度排名列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /gaia/test/app/request/batch [get]
func (quotaApi *TestApi) GaiaAppRequestTestBatch(c *gin.Context) {
	var pageInfo request.GetAppRequestTestRequest
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// list
	lock, list, total, err := TestService.AppRequestTestBatch(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(map[string]interface{}{
		"lock":      lock,
		"list":      list,
		"total":     total,
		"page":      pageInfo.Page,
		"page_size": pageInfo.PageSize,
	}, "获取成功", c)
}
