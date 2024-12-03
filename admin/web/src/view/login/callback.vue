<template>
  <div v-show="false">该页面用于对接oa-oauth2.0回调登录</div>
</template>

<script setup>
import { ElMessage } from 'element-plus'
import { useRoute } from 'vue-router'
import { useUserStore } from '@/pinia/modules/user'

defineOptions({
  name: 'LoginCallback',
})
const route = useRoute()
const userStore = useUserStore()
const oaLogin = async() => {
  return await userStore.OaLoginIn(route.query.code)
}
const callback = async() => {
  if (route.query.code === undefined || route.query.code === '') {
    ElMessage({
      type: 'error',
      message: '登录失败，授权码缺失，3秒后跳转到登录页',
      showClose: true,
    })
    // 3秒后跳转登录页
    setTimeout(() => {
      window.location.href = '/'
    }, 3000)
    return false
  }
  const flag = await oaLogin()
  if (!flag) {
    ElMessage({
      type: 'error',
      message: '登录失败，3秒后跳转到登录页',
      showClose: true,
    })
    // 3秒后跳转登录页
    setTimeout(() => {
      window.location.href = '/'
    }, 3000)
  }
  return
}
callback()
</script>
