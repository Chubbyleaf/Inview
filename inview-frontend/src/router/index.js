import { createRouter, createWebHashHistory } from 'vue-router'
import Login from '../components/Login.vue'
import Dashboard from '../components/Dashboard'

const routes = [
  {
    path: '/', name: 'Login', component: Login
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: Dashboard,
    children: [
      {
        path: '/dashboard/overview',
        name: 'OverviewPanel',
        components: { OverviewPanel: () => import('../view/OverviewPanel.vue') }

      },
      {
        path: '/dashboard/task',
        name: 'TaskManagement',
        components: { TaskManagement: () => import('../view/TaskManagement.vue') }
      },
      {
        path: '/dashboard/log',
        name: 'SystemLog',
        components: { SystemLog: () => import('../view/SystemLog.vue') }
      },
      {
        path: '/dashboard/config',
        name: 'SystemConfig',
        components: { SystemConfig: () => import('../view/SystemConfig') }
      },
      {
        path: '/dashboard/account',
        name: 'AccountPanel',
        components: { AccountPanel: () => import('../view/AccountPanel.vue') }
      },
    ]

  },

]

const router = createRouter({
  // history: createWebHistory(process.env.BASE_URL),
  history: createWebHashHistory(process.env.BASE_URL),
  routes
})


router.beforeEach((to, from, next) => {
  // 判断是否已登录（检查 sessionStorage 中是否有 user_id 信息）
  let userId = localStorage.getItem('user_id');
  // 如果访问的不是登录页面，且没有用户 ID，则跳转到登录页面
  if (to.name !== 'Login' && !userId) {
    next({ name: 'Login' });
  } else {
    next();
  }
});


export default router