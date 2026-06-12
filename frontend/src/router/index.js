import { createRouter, createWebHistory } from 'vue-router'
import { getSession } from '../api/client'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    // Public routes
    {
      path: '/',
      redirect: '/login',
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/auth/LoginView.vue'),
    },

    // User routes
    {
      path: '/user',
      component: () => import('../layouts/UserLayout.vue'),
      meta: { requiresAuth: true, role: 'USER' },
      children: [
        {
          path: '',
          name: 'movies',
          component: () => import('../views/user/MoviesView.vue'),
        },
        {
          path: 'booking/:movieId',
          name: 'booking',
          component: () => import('../views/user/BookingView.vue'),
        },
        {
          path: 'payment',
          name: 'payment',
          component: () => import('../views/user/PaymentView.vue'),
        },
        {
          path: 'history',
          name: 'history',
          component: () => import('../views/user/HistoryView.vue'),
        },
      ],
    },

    // Admin routes
    {
      path: '/admin',
      component: () => import('../layouts/AdminLayout.vue'),
      meta: { requiresAuth: true, role: 'ADMIN' },
      children: [
        {
          path: '',
          name: 'admin-dashboard',
          component: () => import('../views/admin/DashboardView.vue'),
        },
        {
          path: 'audit-logs',
          name: 'admin-logs',
          component: () => import('../views/admin/AuditLogsView.vue'),
        },
      ],
    },
  ],
})

router.beforeEach((to) => {
  const matched = to.matched.find((record) => record.meta.requiresAuth)
  if (!matched) return true

  const session = getSession()
  if (!session?.token) return { name: 'login' }
  if (matched.meta.role === 'ADMIN' && session.user?.role !== 'ADMIN') {
    return { name: 'movies' }
  }
  return true
})

export default router
