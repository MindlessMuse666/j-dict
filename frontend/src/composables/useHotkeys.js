import { onMounted, onUnmounted } from 'vue'

export function useHotkeys(callbacks) {
    const handleKeyDown = (event) => {
        // Игнорируем горячие клавиши при фокусе в полях ввода
        const activeElement = document.activeElement
        const isInput = activeElement.tagName === 'INPUT' ||
            activeElement.tagName === 'TEXTAREA' ||
            activeElement.isContentEditable

        if (isInput) {
            return
        }

        // Shift + N для нового слова
        if (event.shiftKey && (event.key === 'N' || event.key === 'n' || event.key === 'ы' || event.key === 'Ы')) {
            event.preventDefault()
            event.stopPropagation()
            callbacks.onNewWord?.()
        }

        // Shift + Escape для закрытия модальных окон
        if (event.shiftKey && event.key === 'Escape') {
            event.preventDefault()
            event.stopPropagation()
            callbacks.onCloseModals?.()
        }

        // Escape для закрытия модальных окон (без Shift)
        if (event.key === 'Escape' && !event.shiftKey) {
            callbacks.onCloseModals?.()
        }

        // Ctrl/Cmd + F для фокуса на поиск
        if ((event.ctrlKey || event.metaKey) && event.key === 'f') {
            event.preventDefault()
            const searchInput = document.querySelector('input[placeholder*="японском"]')
            if (searchInput) {
                searchInput.focus()
                searchInput.select()
            }
        }

        // Ctrl/Cmd + S для сохранения
        if ((event.ctrlKey || event.metaKey) && event.key === 's') {
            event.preventDefault()
            callbacks.onSave?.()
        }
    }

    onMounted(() => {
        document.addEventListener('keydown', handleKeyDown)
    })

    onUnmounted(() => {
        document.removeEventListener('keydown', handleKeyDown)
    })
}

// Альтернативная версия для использования в setup
export function setupHotkeys(callbacks) {
    const handler = (event) => {
        // Игнорируем горячие клавиши при фокусе в полях ввода
        const activeElement = document.activeElement
        const isInput = activeElement.tagName === 'INPUT' ||
            activeElement.tagName === 'TEXTAREA' ||
            activeElement.isContentEditable

        if (isInput) {
            return
        }

        // Shift + N для нового слова
        if (event.shiftKey && (event.key === 'N' || event.key === 'n' || event.key === 'ы' || event.key === 'Ы')) {
            event.preventDefault()
            event.stopPropagation()
            callbacks.onNewWord?.()
        }

        // Shift + Escape для закрытия модальных окон
        if (event.shiftKey && event.key === 'Escape') {
            event.preventDefault()
            event.stopPropagation()
            callbacks.onCloseModals?.()
        }

        // Escape для закрытия модальных окон (без Shift)
        if (event.key === 'Escape' && !event.shiftKey) {
            callbacks.onCloseModals?.()
        }
    }

    onMounted(() => {
        document.addEventListener('keydown', handler)
    })

    onUnmounted(() => {
        document.removeEventListener('keydown', handler)
    })
}