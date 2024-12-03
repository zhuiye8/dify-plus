import { get } from './base'

// 二开部分Begin - 额度限制，应用web端对登录用户计费
export const fetchAccessTokenExtend = async (appCode: string, loginToken: string) => {
  const headers = new Headers()
  headers.append('X-App-Code', appCode)
  headers.append('Authorization-extend', `Bearer ${loginToken}`)
  return get('/passport-extend', { headers }) as Promise<{ access_token: string }>
}
// 二开部分End - 额度限制，应用web端对登录用户计费
