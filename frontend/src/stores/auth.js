import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { api } from '@/api/axios'

export const useAuthStore = defineStore('auth', () => {
    const user = ref(null)
    const token = ref(localStorage.getItem('token'))

    const isAuthenticated = computed(() => !!token.value)

    const setAuth = (userData, authToken) => {
        user.value = userData
        token.value = authToken
        localStorage.setItem('token', authToken)
    }

    const clearAuth = () => {
        user.value = null
        token.value = null
        localStorage.removeItem('token')
    }

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

    const logout = async () => {
        try {
            await api.post('/auth/logout')
        } catch (error) {
            console.error('Ошибка при выходе:', error)
        } finally {
            clearAuth()
        }
    }

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

    const uploadAvatar = async (file) => {
        try {
            const formData = new FormData()
            formData.append('file', file)
            const response = await api.post('/users/avatar', formData, {
                headers: { 'Content-Type': 'multipart/form-data' }
            })
            if (user.value) {
                user.value.avatar_url = response.data.url
            }
            return { success: true }
        } catch (error) {
            return { success: false, error: error.response?.data?.error || 'Ошибка загрузки аватара' }
        }
    }

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

    // При инициализации загружаем данные пользователя, если есть токен
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