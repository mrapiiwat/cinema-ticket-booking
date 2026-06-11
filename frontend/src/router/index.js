import { createRouter, createWebHistory } from 'vue-router'

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
          beforeEnter: (_to, from, next) => {
            if (from.name === 'booking') {
              next()
            } else {
              next({ name: 'movies' })
            }
          },
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

export default router
