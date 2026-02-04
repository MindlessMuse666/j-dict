import { defineStore } from 'pinia'
import { ref } from 'vue'
import { api } from '@/api/axios'

export const useWordsStore = defineStore('words', () => {
    const words = ref([])
    const totalCount = ref(0)
    const loading = ref(false)
    const error = ref(null)
    const hasMore = ref(true)
    const nextCursor = ref(0)

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

    const searchWords = async (params = {}) => {
        loading.value = true
        error.value = null

        try {
            // Сбрасываем пагинацию при поиске
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