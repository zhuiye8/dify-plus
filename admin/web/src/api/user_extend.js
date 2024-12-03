import service from '@/utils/request'

// @Summary 用户OA登录
// @Produce  application/json
// @Param data body {authorize_code:"string"}
// @Router /base/login [post]
export const oaLogin = (data) => {
  return service({
    url: '/base/oaLogin',
    method: 'post',
    data: data
  })
}
