import service from '@/utils/request'

// @Tags Tenants
// @Summary 用id查询tenants表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.Tenants true "用id查询tenants表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /tenants/findTenants [get]
export const findTenants = (params) => {
  return service({
    url: '/tenants/findTenants',
    method: 'get',
    params
  })
}

// @Tags Tenants
// @Summary 分页获取tenants表列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取tenants表列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /tenants/getTenantsList [get]
export const getTenantsList = (params) => {
  return service({
    url: '/tenants/getTenantsList',
    method: 'get',
    params
  })
}

// @Tags Tenants
// @Summary 获取所有工作区
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.Tenants true "获取所有工作区"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /tenants/getAllTenants [get]
export const getAllTenants = (params) => {
  return service({
    url: '/tenants/getAllTenants',
    method: 'get',
    params
  })
}
