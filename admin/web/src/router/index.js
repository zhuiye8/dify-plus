import { createRouter, createWebHashHistory } from 'vue-router'

const routes = [{
  path: '/',
  redirect: '/login'
},
{
  path: '/init',
  name: 'Init',
  component: () => import('@/view/init/index.vue')
},
{
  path: '/login',
  name: 'Login',
  component: () => import('@/view/login/index.vue')
},
// 新增OA登录 Begin
{
  path: '/loginCallback',
  name: 'LoginCallback',
  component: () => import('@/view/login/callback.vue')
},
// 新增OA登录 End
{
  path: '/:catchAll(.*)',
  meta: {
    closeTab: true,
  },
  component: () => import('@/view/error/index.vue')
}
]

const router = createRouter({
  history: createWebHashHistory(),
  routes,
})

export default router
