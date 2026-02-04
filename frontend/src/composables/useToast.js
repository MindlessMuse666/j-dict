import { ref, provide, inject } from 'vue'

// Создаем глобальный ключ для инъекции
const ToastSymbol = Symbol.for('toast')

// Создаем глобальное хранилище тостов
const globalToasts = ref([])

export function createToast() {
    const toasts = ref([])

    const showToast = (message, type = 'info', duration = 3000) => {
        const id = Date.now() + Math.random()
        const toast = { id, message, type }

        globalToasts.value.push(toast)

        // Автоматическое скрытие
        setTimeout(() => {
            removeToast(id)
        }, duration)

        return id
    }

    const showSuccess = (message, duration = 3000) => {
        return showToast(message, 'success', duration)
    }

    const showError = (message, duration = 5000) => {
        return showToast(message, 'error', duration)
    }

    const showWarning = (message, duration = 4000) => {
        return showToast(message, 'warning', duration)
    }

    const showInfo = (message, duration = 3000) => {
        return showToast(message, 'info', duration)
    }

    const removeToast = (id) => {
        const index = globalToasts.value.findIndex(toast => toast.id === id)
        if (index !== -1) {
            globalToasts.value.splice(index, 1)
        }
    }

    return {
        toasts: globalToasts,
        showToast,
        showSuccess,
        showError,
        showWarning,
        showInfo,
        removeToast
    }
}

// Composition function для использования в компонентах
export function useToast() {
    const toast = inject(ToastSymbol)

    if (!toast) {
        // Fallback: create a new toast instance if not provided (useful for async/outside setup)
        // But better is to just return the methods that operate on globalToasts directly
        // Since globalToasts is defined in module scope, we can access it.
        // This avoids the error when useToast is called outside component setup or before provide
        return createToast()
    }

    return toast
}

// Функция для предоставления тостов в корневом компоненте
export function provideToast(app) {
    const toast = createToast()
    if (app) {
        app.provide(ToastSymbol, toast)
    } else {
        provide(ToastSymbol, toast)
    }
    return toast
}