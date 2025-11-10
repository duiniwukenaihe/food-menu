import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { storeToRefs } from 'pinia'

import { useAuthStore } from '../stores/auth'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'home',
    component: () => import('../views/HomeView.vue'),
  },
  {
    path: '/dish/:id',
    name: 'dish-detail',
    component: () => import('../views/DishDetailView.vue'),
    props: true,
  },
  {
    path: '/admin',
    name: 'admin',
    component: () => import('../views/AdminDashboard.vue'),
    meta: { requiresAuth: true, requiresAdmin: true },
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'not-found',
    component: () => import('../views/NotFoundView.vue'),
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior: () => ({ top: 0 }),
})

router.beforeEach((to, _from, next) => {
  const auth = useAuthStore()
  const { isAuthenticated, role } = storeToRefs(auth)

  if (to.meta.requiresAuth && !isAuthenticated.value) {
    return next({ name: 'home', query: { redirect: to.fullPath } })
  }

  if (to.meta.requiresAdmin && role.value !== 'admin') {
    return next({ name: 'home' })
  }

  return next()
})

export default router
