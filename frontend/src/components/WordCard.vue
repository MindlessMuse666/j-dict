<template>
    <div :class="[
        'bg-surface rounded-xl border-none shadow-sm hover:shadow-lg transition-all duration-300 transform hover:-translate-y-1',
        isCompactView ? 'p-4' : 'p-6'
    ]" :data-word-id="word.id">
        <!-- Компактный вид -->
        <div v-if="isCompactView" class="flex items-center justify-between">
            <!-- Японский текст -->
            <div class="flex-1 mr-4">
                <div class="flex flex-wrap items-center gap-3 mb-1">
                    <span v-for="(item, index) in word.jp" :key="'jp-' + index"
                        class="font-jp text-xl font-medium text-text-main tracking-wide">
                        {{ item }}
                    </span>
                </div>

                <!-- Чтения -->
                <div v-if="word.on || word.kun" class="flex flex-wrap gap-2 mt-2">
                    <span v-for="(reading, index) in word.on" :key="'on-' + index"
                        class="font-jp text-xs px-2 py-1 bg-onyomi text-amber-900 rounded-md shadow-sm">
                        {{ reading }}
                    </span>
                    <span v-for="(reading, index) in word.kun" :key="'kun-' + index"
                        class="font-jp text-xs px-2 py-1 bg-kunyomi text-sky-900 rounded-md shadow-sm">
                        {{ reading }}
                    </span>
                </div>
            </div>

            <!-- Русский перевод -->
            <div class="flex-1 text-right">
                <div class="flex flex-wrap justify-end gap-2 mb-2">
                    <span v-for="(item, index) in word.ru" :key="'ru-' + index" class="text-text-main font-medium">
                        {{ item }}
                    </span>
                </div>

                <!-- Теги -->
                <div v-if="word.tags && word.tags.length > 0" class="flex flex-wrap justify-end gap-1.5 mt-1">
                    <span v-for="(tag, index) in word.tags" :key="'tag-' + index"
                        class="text-xs px-2 py-0.5 bg-stone-100 text-stone-500 rounded-full border border-stone-200">
                        #{{ tag }}
                    </span>
                </div>
            </div>

            <!-- Кнопки действий -->
            <div class="ml-5 pl-4 border-l border-stone-100 flex items-center space-x-1">
                <button @click="handleEdit"
                    class="p-2 text-stone-400 hover:text-primary hover:bg-rose-50 rounded-full transition-colors"
                    title="Редактировать">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
                            d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                    </svg>
                </button>
                <button @click="handleDelete"
                    class="p-2 text-stone-400 hover:text-danger hover:bg-red-50 rounded-full transition-colors"
                    title="Удалить">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
                            d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                    </svg>
                </button>
            </div>
        </div>

        <!-- Полный вид -->
        <div v-else>
            <!-- Заголовок с японскими написаниями -->
            <div class="mb-5 pb-4 border-b border-stone-100 border-dashed">
                <div class="flex flex-wrap items-center gap-4 mb-3">
                    <h3 v-for="(item, index) in word.jp" :key="'jp-full-' + index"
                        class="font-jp text-3xl font-medium text-text-main tracking-wide">
                        {{ item }}
                    </h3>
                </div>

                <!-- Чтения -->
                <div v-if="word.on || word.kun" class="flex flex-wrap gap-3">
                    <div v-if="word.on && word.on.length > 0" class="flex items-center">
                        <span class="text-xs font-bold text-amber-700 uppercase tracking-wider mr-2">On</span>
                        <span v-for="(reading, index) in word.on" :key="'on-full-' + index"
                            class="font-jp inline-block text-sm px-3 py-1 bg-onyomi text-amber-900 rounded-md mr-2 shadow-sm">
                            {{ reading }}
                        </span>
                    </div>
                    <div v-if="word.kun && word.kun.length > 0" class="flex items-center ml-2">
                        <span class="text-xs font-bold text-sky-700 uppercase tracking-wider mr-2">Kun</span>
                        <span v-for="(reading, index) in word.kun" :key="'kun-full-' + index"
                            class="font-jp inline-block text-sm px-3 py-1 bg-kunyomi text-sky-900 rounded-md mr-2 shadow-sm">
                            {{ reading }}
                        </span>
                    </div>
                </div>
            </div>

            <!-- Русские переводы -->
            <div class="mb-6">
                <h4 class="text-sm font-medium text-text-muted mb-2">Переводы:</h4>
                <ul class="space-y-1">
                    <li v-for="(item, index) in word.ru" :key="'ru-full-' + index" class="text-text-main">
                        {{ item }}
                    </li>
                </ul>
            </div>

            <!-- Примеры -->
            <div v-if="word.ex_jp && word.ex_jp.length > 0" class="mb-6">
                <h4 class="text-sm font-medium text-text-muted mb-2">Примеры на японском:</h4>
                <ul class="space-y-2">
                    <li v-for="(example, index) in word.ex_jp" :key="'ex-jp-' + index" class="jp-text text-text-main">
                        {{ example }}
                    </li>
                </ul>
            </div>

            <div v-if="word.ex_ru && word.ex_ru.length > 0" class="mb-6">
                <h4 class="text-sm font-medium text-text-muted mb-2">Переводы примеров:</h4>
                <ul class="space-y-2">
                    <li v-for="(example, index) in word.ex_ru" :key="'ex-ru-' + index" class="text-text-main">
                        {{ example }}
                    </li>
                </ul>
            </div>

            <!-- Теги и метаданные -->
            <div class="flex items-center justify-between pt-4 border-t border-stone-100">
                <div>
                    <div v-if="word.tags && word.tags.length > 0" class="flex flex-wrap gap-2">
                        <span v-for="(tag, index) in word.tags" :key="'tag-full-' + index"
                            class="text-xs px-3 py-1 bg-stone-100 text-text-muted rounded-full">
                            {{ tag }}
                        </span>
                    </div>
                    <div v-else class="text-xs text-stone-400">
                        Теги не указаны
                    </div>
                </div>

                <div class="text-xs text-stone-400">
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
import { format } from 'date-fns'

/**
 * Свойства компонента WordCard.
 * @typedef {Object} Props
 * @property {Object} word - Объект слова для отображения.
 * @property {Boolean} [isCompactView=true] - Отображать ли карточку в компактном режиме.
 */
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

/**
 * События компонента WordCard.
 * @emits edit - Вызывается при нажатии кнопки редактирования. Payload: объект слова.
 * @emits delete - Вызывается при нажатии кнопки удаления. Payload: ID слова.
 */
const emit = defineEmits(['edit', 'delete'])

/**
 * Обрабатывает клик по кнопке редактирования.
 * Генерирует событие 'edit' с текущим словом.
 */
const handleEdit = () => {
    emit('edit', props.word)
}

/**
 * Обрабатывает клик по кнопке удаления.
 * Генерирует событие 'delete' с ID текущего слова.
 */
const handleDelete = () => {
    emit('delete', props.word.id)
}

/**
 * Форматирует строку даты в читаемый формат.
 * @param {string} dateString - Строка даты для форматирования.
 * @returns {string} Отформатированная дата (dd.MM.yyyy HH:mm) или исходная строка при ошибке.
 */
const formatDate = (dateString) => {
    if (!dateString) return ''
    try {
        return format(new Date(dateString), 'dd.MM.yyyy HH:mm')
    } catch {
        return dateString
    }
}
</script>
