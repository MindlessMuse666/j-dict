import axios from 'axios'
import router from '@/router'
import { useAuthStore } from '@/stores/auth'
import { useToast } from '@/composables/useToast'

const api = axios.create({
    baseURL: import.meta.env.VITE_API_URL || '/api',
    headers: {
        'Content-Type': 'application/json',
    },
})

// Интерцептор для добавления токена к запросам
api.interceptors.request.use(
    (config) => {
        const authStore = useAuthStore()
        const token = authStore.token

        if (token) {
            config.headers.Authorization = `Bearer ${token}`
        }

        return config
    },
    (error) => {
        return Promise.reject(error)
    }
)

// Интерцептор для обработки ошибок
api.interceptors.response.use(
    (response) => response,
    async (error) => {
        const authStore = useAuthStore()

        // Создаем экземпляр тостов
        let toast
        try {
            const toastModule = await import('@/composables/useToast')
            toast = toastModule.useToast()
        } catch (e) {
            console.error('Не удалось загрузить модуль тостов:', e)
            // Фолбэк на консоль
            const showError = (msg) => console.error('Toast Error:', msg)
            toast = { showError }
        }

        if (error.response?.status === 401) {
            authStore.clearAuth()
            if (!['Login', 'Register'].includes(router.currentRoute.value.name)) {
                router.push({ name: 'Login', query: { redirect: router.currentRoute.value.fullPath } })
            }
        } else if (error.response?.status === 429) {
            toast.showError?.('Слишком много запросов. Пожалуйста, подождите.')
        } else if (error.code === 'ERR_NETWORK') {
            toast.showError?.('Нет соединения с сервером. Проверьте подключение к интернету.')
        } else if (error.response?.status === 413) {
            toast.showError?.('Файл слишком большой. Максимальный размер: 10MB.')
        } else if (error.response?.status >= 500) {
            toast.showError?.('Внутренняя ошибка сервера. Пожалуйста, попробуйте позже.')
        }

        return Promise.reject(error)
    }
)

export { api }