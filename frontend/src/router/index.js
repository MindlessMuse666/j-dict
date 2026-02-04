import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

/**
 * Определение маршрутов приложения.
 * @type {import('vue-router').RouteRecordRaw[]}
 */
const routes = [
    {
        path: '/',
        name: 'Home',
        component: () => import('@/views/HomeView.vue'),
        meta: { requiresAuth: true }
    },
    {
        path: '/login',
        name: 'Login',
        component: () => import('@/views/LoginView.vue'),
        meta: { guestOnly: true }
    },
    {
        path: '/register',
        name: 'Register',
        component: () => import('@/views/RegisterView.vue'),
        meta: { guestOnly: true }
    },
    {
        path: '/profile',
        name: 'Profile',
        component: () => import('@/views/ProfileView.vue'),
        meta: { requiresAuth: true }
    },
    {
        path: '/:pathMatch(.*)*',
        name: 'NotFound',
        component: () => import('@/views/NotFoundView.vue')
    }
]

/**
 * Экземпляр Vue Router.
 */
const router = createRouter({
    history: createWebHistory(),
    routes
})

/**
 * Глобальный навигационный хук для обработки аутентификации.
 * Перенаправляет на Login, если требуется auth и пользователь не авторизован.
 * Перенаправляет на Home, если пользователь авторизован и пытается получить доступ к страницам только для гостей.
 */
router.beforeEach((to, from, next) => {
    const authStore = useAuthStore()
    const isAuthenticated = authStore.isAuthenticated

    if (to.meta.requiresAuth && !isAuthenticated) {
        // Перенаправляем на страницу входа, если требуется аутентификация
        next({ name: 'Login', query: { redirect: to.fullPath } })
    } else if (to.meta.guestOnly && isAuthenticated) {
        // Если пользователь уже авторизован, перенаправляем на главную
        next({ name: 'Home' })
    } else {
        next()
    }
})

export default router