import service from '@/utils/request'

// @Tags dashboard
// @Summary 分页获取账户配额排名数据
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取账户配额排名数据"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /gaia/dashboard/getAccountQuotaRankingData [get]
export const getAccountQuotaRankingData = (params) => {
  return service({
    url: '/gaia/dashboard/getAccountQuotaRankingData',
    method: 'get',
    params
  })
}

// @Tags dashboard
// @Summary 分页获取【应用】配额排名数据
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取【应用】配额排名数据"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /gaia/dashboard/getAppQuotaRankingData [get]
export const getAppQuotaRankingData = (params) => {
  return service({
    url: '/gaia/dashboard/getAppQuotaRankingData',
    method: 'get',
    params
  })
}

// @Tags dashboard
// @Summary 分页获取【应用密钥】配额排名数据列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取【应用密钥】配额排名数据列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /gaia/dashboard/getAppTokenQuotaRankingData [get]
export const getAppTokenQuotaRankingData = (params) => {
  return service({
    url: '/gaia/dashboard/getAppTokenQuotaRankingData',
    method: 'get',
    params
  })
}


// @Tags dashboard
// @Summary 获取每天密钥花费数据列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "获取每天密钥花费数据列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /gaia/dashboard/getAppTokenDailyQuotaData [get]
export const getAppTokenDailyQuotaData = (params) => {
  return service({
    url: '/gaia/dashboard/getAppTokenDailyQuotaData',
    method: 'get',
    params
  })
}


// @Tags dashboard
// @Summary 获取AI作图花费排行列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "获取AI作图花费排行列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /gaia/dashboard/getAiImageQuotaRankingData [get]
export const getAiImageQuotaRankingData = (params) => {
  return service({
    url: '/gaia/dashboard/getAiImageQuotaRankingData',
    method: 'get',
    params
  })
}

