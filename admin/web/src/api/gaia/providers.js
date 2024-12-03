import service from '@/utils/request'

// @Tags Providers
// @Summary 更新providers表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Providers true "更新providers表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /providers/syncProviders [put]
export const syncProviders = (data) => {
  return service({
    url: '/providers/syncProviders',
    method: 'put',
    data
  })
}

// @Tags Providers
// @Summary 用id查询providers表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.Providers true "用id查询providers表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /providers/findProviders [get]
export const findProviders = (params) => {
  return service({
    url: '/providers/findProviders',
    method: 'get',
    params
  })
}

// @Tags Providers
// @Summary 分页获取providers表列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取providers表列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /providers/getProvidersList [get]
export const getProvidersList = (params) => {
  return service({
    url: '/providers/getProvidersList',
    method: 'get',
    params
  })
}