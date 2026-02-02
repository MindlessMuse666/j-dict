<template>
    <div :class="[
        'bg-white rounded-lg border transition-all duration-300 hover:shadow-md interactive',
        isCompactView ? 'p-4' : 'p-6',
        wordState
    ]" :data-word-id="word.id">
        <!-- Компактный вид -->
        <div v-if="isCompactView" class="flex items-center justify-between">
            <!-- Японский текст -->
            <div class="flex-1 mr-4">
                <div class="flex flex-wrap items-center gap-2 mb-1">
                    <span v-for="(item, index) in word.jp" :key="'jp-' + index"
                        class="jp-text text-lg font-semibold text-gray-900">
                        {{ item }}
                    </span>
                </div>

                <!-- Чтения -->
                <div v-if="word.on || word.kun" class="flex flex-wrap gap-2 mt-1">
                    <span v-for="(reading, index) in word.on" :key="'on-' + index"
                        class="onyomi text-xs px-2 py-1 bg-onyomi/20 rounded">
                        {{ reading }}
                    </span>
                    <span v-for="(reading, index) in word.kun" :key="'kun-' + index"
                        class="kunyomi text-xs px-2 py-1 bg-kunyomi/20 rounded">
                        {{ reading }}
                    </span>
                </div>
            </div>

            <!-- Русский перевод -->
            <div class="flex-1 text-right">
                <div class="flex flex-wrap justify-end gap-2 mb-1">
                    <span v-for="(item, index) in word.ru" :key="'ru-' + index" class="text-gray-700">
                        {{ item }}
                    </span>
                </div>

                <!-- Теги -->
                <div v-if="word.tags && word.tags.length > 0" class="flex flex-wrap justify-end gap-1 mt-1">
                    <span v-for="(tag, index) in word.tags" :key="'tag-' + index"
                        class="text-xs px-2 py-0.5 bg-gray-100 text-gray-600 rounded">
                        {{ tag }}
                    </span>
                </div>
            </div>

            <!-- Кнопки действий -->
            <div class="ml-4 flex items-center space-x-2">
                <button @click="handleEdit"
                    class="p-2 text-gray-500 hover:text-primary hover:bg-gray-100 rounded-full transition-colors"
                    title="Редактировать">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                    </svg>
                </button>
                <button @click="handleDelete"
                    class="p-2 text-gray-500 hover:text-danger hover:bg-red-50 rounded-full transition-colors"
                    title="Удалить">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                    </svg>
                </button>
            </div>
        </div>

        <!-- Полный вид -->
        <div v-else>
            <!-- Заголовок с японскими написаниями -->
            <div class="mb-4">
                <div class="flex flex-wrap items-center gap-3 mb-3">
                    <h3 v-for="(item, index) in word.jp" :key="'jp-full-' + index"
                        class="jp-text text-2xl font-bold text-gray-900">
                        {{ item }}
                    </h3>
                </div>

                <!-- Чтения -->
                <div v-if="word.on || word.kun" class="flex flex-wrap gap-3 mb-3">
                    <div v-if="word.on && word.on.length > 0">
                        <span class="text-sm font-medium text-gray-600 mr-2">Онъёми:</span>
                        <span v-for="(reading, index) in word.on" :key="'on-full-' + index"
                            class="onyomi inline-block text-sm px-3 py-1 bg-onyomi/20 rounded-full mr-2 mb-1">
                            {{ reading }}
                        </span>
                    </div>
                    <div v-if="word.kun && word.kun.length > 0">
                        <span class="text-sm font-medium text-gray-600 mr-2">Кунъёми:</span>
                        <span v-for="(reading, index) in word.kun" :key="'kun-full-' + index"
                            class="kunyomi inline-block text-sm px-3 py-1 bg-kunyomi/20 rounded-full mr-2 mb-1">
                            {{ reading }}
                        </span>
                    </div>
                </div>
            </div>

            <!-- Русские переводы -->
            <div class="mb-6">
                <h4 class="text-sm font-medium text-gray-600 mb-2">Переводы:</h4>
                <ul class="space-y-1">
                    <li v-for="(item, index) in word.ru" :key="'ru-full-' + index" class="text-gray-800">
                        {{ item }}
                    </li>
                </ul>
            </div>

            <!-- Примеры -->
            <div v-if="word.ex_jp && word.ex_jp.length > 0" class="mb-6">
                <h4 class="text-sm font-medium text-gray-600 mb-2">Примеры на японском:</h4>
                <ul class="space-y-2">
                    <li v-for="(example, index) in word.ex_jp" :key="'ex-jp-' + index" class="jp-text text-gray-800">
                        {{ example }}
                    </li>
                </ul>
            </div>

            <div v-if="word.ex_ru && word.ex_ru.length > 0" class="mb-6">
                <h4 class="text-sm font-medium text-gray-600 mb-2">Переводы примеров:</h4>
                <ul class="space-y-2">
                    <li v-for="(example, index) in word.ex_ru" :key="'ex-ru-' + index" class="text-gray-800">
                        {{ example }}
                    </li>
                </ul>
            </div>

            <!-- Теги и метаданные -->
            <div class="flex items-center justify-between pt-4 border-t">
                <div>
                    <div v-if="word.tags && word.tags.length > 0" class="flex flex-wrap gap-2">
                        <span v-for="(tag, index) in word.tags" :key="'tag-full-' + index"
                            class="text-xs px-3 py-1 bg-gray-100 text-gray-600 rounded-full">
                            {{ tag }}
                        </span>
                    </div>
                    <div v-else class="text-xs text-gray-400">
                        Теги не указаны
                    </div>
                </div>

                <div class="text-xs text-gray-500">
                    {{ formatDate(word.updated_at) }}
                </div>
            </div>

            <!-- Кнопки действий -->
            <div class="mt-4 flex justify-end space-x-3">
                <button @click="handleEdit"
                    class="px-4 py-2 text-sm font-medium text-primary hover:bg-primary/10 rounded-md transition-colors">
                    Редактировать
                </button>
                <button @click="handleDelete"
                    class="px-4 py-2 text-sm font-medium text-danger hover:bg-danger/10 rounded-md transition-colors">
                    Удалить
                </button>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'
import { useToast } from '@/composables/useToast'
import { useWordsStore } from '@/stores/words'
import { format } from 'date-fns'

const props = defineProps({
    word: {
        type: Object,
        required: true
    },
    isCompactView: {
        type: Boolean,
        default: true
    }
})

const emit = defineEmits(['edit', 'delete'])

const { showSuccess, showError } = useToast()
const wordsStore = useWordsStore()
const wordState = ref('')

onMounted(() => {
    wordState.value = 'word-card-enter'
})

const handleEdit = () => {
    emit('edit', props.word)
}

const handleDelete = async () => {
    if (!confirm('Вы уверены, что хотите удалить это слово?')) return

    try {
        // Анимация удаления
        wordState.value = 'word-removing'

        // Ждем завершения анимации перед удалением
        setTimeout(async () => {
            const result = await wordsStore.deleteWord(props.word.id)

            if (result.success) {
                showSuccess('Слово успешно удалено')
                emit('delete', props.word.id)
            } else {
                showError(result.error)
                wordState.value = ''
            }
        }, 250)
    } catch (error) {
        showError('Ошибка при удалении слова')
        wordState.value = ''
    }
}

const formatDate = (dateString) => {
    if (!dateString) return ''
    try {
        return format(new Date(dateString), 'dd.MM.yyyy HH:mm')
    } catch {
        return dateString
    }
}

// Анимация обновления при изменении пропсов
watch(() => props.word, (newWord, oldWord) => {
    if (newWord.updated_at !== oldWord?.updated_at) {
        wordState.value = 'word-updated'
        setTimeout(() => {
            wordState.value = ''
        }, 1000)
    }
}, { deep: true })
</script>