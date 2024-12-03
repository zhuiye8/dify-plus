import { CONVERSATION_ID_INFO } from '../base/chat/constants'
// import { fetchAccessToken } from '@/service/share' // 二开部分 - 额度限制，应用web端对登录用户计费
import { fetchAccessTokenExtend } from '@/service/web-extend'

export const checkOrSetAccessToken = async () => {
  const sharedToken = globalThis.location.pathname.split('/').slice(-1)[0]
  const accessToken = localStorage.getItem('token') || JSON.stringify({ [sharedToken]: '' })
  let accessTokenJson = { [sharedToken]: '' }
  try {
    accessTokenJson = JSON.parse(accessToken)
  }
  catch (e) {

  }
  if (!accessTokenJson[sharedToken]) {
    const loginToken = localStorage.getItem('console_token') || '' // 二开部分 - 额度限制，应用web端对登录用户计费
    const res = await fetchAccessTokenExtend(sharedToken, loginToken) // 二开部分 - 额度限制，应用web端对登录用户计费
    accessTokenJson[sharedToken] = res.access_token
    localStorage.setItem('token', JSON.stringify(accessTokenJson))
  }
}

export const setAccessToken = async (sharedToken: string, token: string) => {
  const accessToken = localStorage.getItem('token') || JSON.stringify({ [sharedToken]: '' })
  let accessTokenJson = { [sharedToken]: '' }
  try {
    accessTokenJson = JSON.parse(accessToken)
  }
  catch (e) {

  }

  localStorage.removeItem(CONVERSATION_ID_INFO)

  accessTokenJson[sharedToken] = token
  localStorage.setItem('token', JSON.stringify(accessTokenJson))
}

export const removeAccessToken = () => {
  const sharedToken = globalThis.location.pathname.split('/').slice(-1)[0]

  const accessToken = localStorage.getItem('token') || JSON.stringify({ [sharedToken]: '' })
  let accessTokenJson = { [sharedToken]: '' }
  try {
    accessTokenJson = JSON.parse(accessToken)
  }
  catch (e) {

  }

  localStorage.removeItem(CONVERSATION_ID_INFO)

  delete accessTokenJson[sharedToken]
  localStorage.setItem('token', JSON.stringify(accessTokenJson))
}
