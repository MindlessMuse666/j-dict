import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { api } from '@/api/axios'

/**
 * Стор для управления аутентификацией и состоянием пользователя.
 */
export const useAuthStore = defineStore('auth', () => {
    // Состояние
    const user = ref(null)
    const token = ref(localStorage.getItem('token'))

    // Вычисляемые свойства
    const isAuthenticated = computed(() => !!token.value)

    /**
     * Устанавливает данные аутентификации (пользователь и токен).
     * @param {Object} userData - Объект пользователя.
     * @param {string} authToken - JWT токен.
     */
    const setAuth = (userData, authToken) => {
        user.value = userData
        token.value = authToken
        localStorage.setItem('token', authToken)
    }

    /**
     * Очищает данные аутентификации.
     */
    const clearAuth = () => {
        user.value = null
        token.value = null
        localStorage.removeItem('token')
    }

    /**
     * Выполняет вход пользователя.
     * @param {string} email - Email пользователя.
     * @param {string} password - Пароль пользователя.
     * @returns {Promise<Object>} Результат с успехом.
     */
    const login = async (email, password) => {
        try {
            const response = await api.post('/auth/login', { email, password })
            const { data } = response.data

            setAuth(data.user, data.token)
            return { success: true }
        } catch (error) {
            return {
                success: false,
                error: error.response?.data?.error || 'Ошибка входа'
            }
        }
    }

    /**
     * Регистрирует нового пользователя.
     * @param {string} email - Email пользователя.
     * @param {string} password - Пароль пользователя.
     * @param {string} name - Имя пользователя.
     * @returns {Promise<Object>} Результат с успехом.
     */
    const register = async (email, password, name) => {
        try {
            const response = await api.post('/auth/register', { email, password, name })
            const { data } = response.data

            setAuth(data.user, data.token)
            return { success: true }
        } catch (error) {
            return {
                success: false,
                error: error.response?.data?.error || 'Ошибка регистрации'
            }
        }
    }

    /**
     * Выполняет выход пользователя.
     */
    const logout = async () => {
        try {
            await api.post('/auth/logout')
        } catch (error) {
            console.error('Ошибка при выходе:', error)
        } finally {
            clearAuth()
        }
    }

    /**
     * Получает данные текущего пользователя.
     */
    const fetchUser = async () => {
        if (!token.value) return

        try {
            const response = await api.get('/auth/me')
            user.value = response.data.data
        } catch (error) {
            console.error('Ошибка при получении данных пользователя:', error)
            if (error.response?.status === 401) {
                clearAuth()
            }
        }
    }

    /**
     * Загружает новый аватар для пользователя.
     * @param {File} file - Файл изображения.
     * @returns {Promise<Object>} Результат с успехом.
     */
    const uploadAvatar = async (file) => {
        try {
            const formData = new FormData()
            // Гарантируем, что имя файла имеет расширение .jpg для валидации на бэкенде
            formData.append('file', file, 'avatar.jpg')
            const response = await api.post('/users/avatar', formData, {
                headers: { 'Content-Type': 'multipart/form-data' },
                skipGlobalError: true
            })
            if (user.value) {
                user.value.avatar_url = response.data.url
            }
            return { success: true }
        } catch (error) {
            return { success: false, error: error.response?.data?.error || 'Ошибка загрузки аватара' }
        }
    }

    /**
     * Обновляет URL аватара напрямую.
     * @param {string} url - Новый URL аватара.
     * @returns {Promise<Object>} Результат с успехом.
     */
    const updateAvatarUrl = async (url) => {
        try {
            await api.put('/users/avatar', { avatar_url: url })
            if (user.value) {
                user.value.avatar_url = url
            }
            return { success: true }
        } catch (error) {
            return { success: false, error: error.response?.data?.error || 'Ошибка обновления аватара' }
        }
    }

    // Инициализируем данные пользователя, если есть токен
    if (token.value) {
        fetchUser()
    }

    return {
        user,
        token,
        isAuthenticated,
        login,
        register,
        logout,
        fetchUser,
        clearAuth,
        uploadAvatar,
        updateAvatarUrl
    }
})
