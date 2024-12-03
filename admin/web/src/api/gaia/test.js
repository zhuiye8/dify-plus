import service from '@/utils/request'


// @Summary gaia应用请求测试结果列表
// @Produce  application/json
// @Param data body {username:"string",password:"string",newPassword:"string"}
// @Router /gaia/test/app/request [get]
export const gaiaAppRequestTesList = (data) => {
  return service({
    url: '/gaia/test/app/request/list',
    method: 'get',
    params: data
  })
}


// @Summary 发起gaia应用请求测试
// @Produce  application/json
// @Param data body {username:"string",password:"string",newPassword:"string"}
// @Router /gaia/test/app/request [post]
export const gaiaAppRequestTest = (data) => {
  return service({
    url: '/gaia/test/app/request',
    method: 'post',
    data: data
  })
}


// @Summary gaia应用请求测试批次列表
// @Produce  application/json
// @Param data body {username:"string",password:"string",newPassword:"string"}
// @Router /gaia/test/app/request/batch [get]
export const gaiaAppRequestBatch = (data) => {
  return service({
    url: '/gaia/test/app/request/batch',
    method: 'get',
    params: data
  })
}


// Extend Stop: Sync User
