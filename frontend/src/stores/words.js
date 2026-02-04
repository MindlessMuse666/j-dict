import { defineStore } from 'pinia'
import { ref } from 'vue'
import { api } from '@/api/axios'

/**
 * Стор для управления словами.
 * Обрабатывает получение, поиск, создание, обновление и удаление слов.
 */
export const useWordsStore = defineStore('words', () => {
    // Состояние
    const words = ref([])
    const totalCount = ref(0)
    const loading = ref(false)
    const error = ref(null)
    const hasMore = ref(true)
    const nextCursor = ref(0)

    /**
     * Получает список слов с пагинацией.
     * @param {number} [limit=20] - Количество слов для получения.
     * @param {boolean} [reset=false] - Сбрасывать ли список и курсор.
     */
    const fetchWords = async (limit = 20, reset = false) => {
        if (loading.value) return

        loading.value = true
        error.value = null

        if (reset) {
            words.value = []
            nextCursor.value = 0
            hasMore.value = true
        }

        try {
            const params = { limit, cursor: nextCursor.value }
            const response = await api.get('/words', { params })
            const { data } = response.data

            words.value = [...words.value, ...data.words]
            nextCursor.value = data.next_cursor
            hasMore.value = data.has_more
            if (data.total_count !== undefined) {
                totalCount.value = data.total_count
            }
        } catch (err) {
            error.value = err.response?.data?.error || 'Ошибка при загрузке слов'
            console.error('Ошибка при загрузке слов:', err)
        } finally {
            loading.value = false
        }
    }

    /**
     * Ищет слова по критериям.
     * @param {Object} [params={}] - Параметры поиска (query, tags, on, kun и т.д.).
     */
    const searchWords = async (params = {}) => {
        loading.value = true
        error.value = null

        try {
            // Сбрасываем пагинацию для поиска
            words.value = []
            nextCursor.value = 0
            hasMore.value = true

            const queryParams = { ...params }

            if (Array.isArray(queryParams.tags)) {
                queryParams.tags = queryParams.tags.join(',')
            }
            if (Array.isArray(queryParams.on)) {
                queryParams.on = queryParams.on.join(',')
            }
            if (Array.isArray(queryParams.kun)) {
                queryParams.kun = queryParams.kun.join(',')
            }

            const response = await api.get('/words/search', {
                params: queryParams,
                paramsSerializer: {
                    indexes: null
                }
            })

            const { data } = response.data

            words.value = data.words || []
            nextCursor.value = data.next_cursor || 0
            hasMore.value = data.has_more || false
            if (data.total_count !== undefined) {
                totalCount.value = data.total_count
            }
        } catch (err) {
            error.value = err.response?.data?.error || 'Ошибка при поиске слов'
            console.error('Ошибка при поиске слов:', err)
        } finally {
            loading.value = false
        }
    }

    /**
     * Создает новое слово.
     * @param {Object} wordData - Данные нового слова.
     * @returns {Promise<Object>} Результат с успехом и данными/ошибкой.
     */
    const createWord = async (wordData) => {
        try {
            const response = await api.post('/words', wordData)
            const newWord = response.data.data

            // Добавляем новое слово в начало списка
            words.value = [newWord, ...words.value]
            totalCount.value++
            return { success: true, word: newWord }
        } catch (err) {
            return {
                success: false,
                error: err.response?.data?.error || 'Ошибка при создании слова'
            }
        }
    }

    /**
     * Обновляет существующее слово.
     * @param {number|string} id - ID слова для обновления.
     * @param {Object} wordData - Обновленные данные.
     * @returns {Promise<Object>} Результат с успехом и данными/ошибкой.
     */
    const updateWord = async (id, wordData) => {
        try {
            const response = await api.patch(`/words/${id}`, wordData)
            const updatedWord = response.data.data

            // Обновляем слово в списке
            const index = words.value.findIndex(word => word.id === id)
            if (index !== -1) {
                words.value[index] = updatedWord
            }

            return { success: true, word: updatedWord }
        } catch (err) {
            return {
                success: false,
                error: err.response?.data?.error || 'Ошибка при обновлении слова'
            }
        }
    }

    /**
     * Удаляет слово.
     * @param {number|string} id - ID слова для удаления.
     * @returns {Promise<Object>} Результат с успехом.
     */
    const deleteWord = async (id) => {
        try {
            await api.delete(`/words/${id}`)

            // Удаляем слово из списка
            words.value = words.value.filter(word => word.id !== id)
            totalCount.value--
            return { success: true }
        } catch (err) {
            return {
                success: false,
                error: err.response?.data?.error || 'Ошибка при удалении слова'
            }
        }
    }

    /**
     * Очищает список слов и сбрасывает состояние.
     */
    const clearWords = () => {
        words.value = []
        nextCursor.value = 0
        hasMore.value = true
        error.value = null
    }

    return {
        words,
        totalCount,
        loading,
        error,
        hasMore,
        nextCursor,
        fetchWords,
        searchWords,
        createWord,
        updateWord,
        deleteWord,
        clearWords
    }
})
