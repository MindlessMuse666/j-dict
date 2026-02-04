<template>
    <Transition name="modal">
        <div v-if="isOpen" class="fixed inset-0 z-50 flex items-center justify-center">
            <!-- Фон с затемнением -->
            <div class="fixed inset-0 bg-black bg-opacity-50" @click="close"></div>

            <!-- Модальное окно - поверх фона -->
            <div class="relative z-50 bg-white rounded-xl shadow-xl w-full max-w-2xl max-h-[70vh] overflow-hidden mx-4">
                <!-- Заголовок -->
                <div class="px-5 py-3 border-b border-stone-200 sticky top-0 z-10"
                    :class="isEditing ? 'bg-amber-50' : 'bg-white'">
                    <div class="flex items-center justify-between">
                        <h3 class="text-lg font-semibold text-text-main font-jp">
                            <span v-if="isEditing">Редактирование</span>
                            <span v-else>Новое слово</span>
                        </h3>
                        <button @click="close" class="text-stone-400 hover:text-stone-600 transition-colors p-1">
                            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                    d="M6 18L18 6M6 6l12 12" />
                            </svg>
                        </button>
                    </div>
                </div>

                <!-- Контент с прокруткой -->
                <div class="overflow-y-auto max-h-[calc(70vh-8rem)] p-5 bg-surface">
                    <!-- Ошибки формы -->
                    <div v-if="Object.keys(errors).length > 0"
                        class="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg">
                        <div class="flex items-start">
                            <svg class="w-5 h-5 text-danger mt-0.5 mr-3" fill="currentColor" viewBox="0 0 20 20">
                                <path fill-rule="evenodd"
                                    d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
                                    clip-rule="evenodd" />
                            </svg>
                            <div>
                                <h4 class="font-medium text-red-800 mb-1">Исправьте следующие ошибки:</h4>
                                <ul class="text-sm text-red-700 list-disc list-inside">
                                    <li v-for="(error, field) in errors" :key="field">
                                        {{ error }}
                                    </li>
                                </ul>
                            </div>
                        </div>
                    </div>

                    <!-- Основная форма -->
                    <form @submit.prevent="handleSubmit" class="space-y-5">
                        <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
                            <!-- Японские написания -->
                            <div>
                                <label class="block text-sm font-medium text-text-main mb-1.5 font-jp">
                                    Японский
                                    <span class="text-primary">*</span>
                                </label>
                                <div class="space-y-2">
                                    <div v-for="(item, index) in form.jp" :key="'jp-' + index"
                                        class="flex items-center space-x-2">
                                        <input v-model="form.jp[index]" type="text" :class="[
                                            'flex-1 px-3 py-2 border rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent transition-colors font-jp text-sm',
                                            errors.jp ? 'border-red-300' : 'border-stone-200'
                                        ]" placeholder="Слово/фраза" />
                                        <button v-if="form.jp.length > 1" type="button"
                                            @click="removeArrayItem('jp', index)"
                                            class="p-1.5 text-stone-400 hover:text-danger transition-colors">
                                            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                                    d="M6 18L18 6M6 6l12 12" />
                                            </svg>
                                        </button>
                                    </div>
                                </div>
                                <button type="button" @click="addArrayItem('jp')"
                                    class="mt-1.5 text-xs text-primary hover:text-primary/80 flex items-center font-medium">
                                    <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                            d="M12 4v16m8-8H4" />
                                    </svg>
                                    Еще вариант
                                </button>
                            </div>

                            <!-- Русские переводы -->
                            <div>
                                <label class="block text-sm font-medium text-text-main mb-1.5">
                                    Русский
                                    <span class="text-primary">*</span>
                                </label>
                                <div class="space-y-2">
                                    <div v-for="(item, index) in form.ru" :key="'ru-' + index"
                                        class="flex items-center space-x-2">
                                        <input v-model="form.ru[index]" type="text" :class="[
                                            'flex-1 px-3 py-2 border rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent transition-colors text-sm',
                                            errors.ru ? 'border-red-300' : 'border-stone-200'
                                        ]" placeholder="Перевод" />
                                        <button v-if="form.ru.length > 1" type="button"
                                            @click="removeArrayItem('ru', index)"
                                            class="p-1.5 text-stone-400 hover:text-danger transition-colors">
                                            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                                    d="M6 18L18 6M6 6l12 12" />
                                            </svg>
                                        </button>
                                    </div>
                                </div>
                                <button type="button" @click="addArrayItem('ru')"
                                    class="mt-1.5 text-xs text-primary hover:text-primary/80 flex items-center font-medium">
                                    <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                            d="M12 4v16m8-8H4" />
                                    </svg>
                                    Еще перевод
                                </button>
                            </div>
                        </div>

                        <!-- Кнопка расширенного ввода -->
                        <div>
                            <button type="button" @click="showAdvanced = !showAdvanced"
                                class="text-sm text-primary hover:text-primary/80 flex items-center font-medium">
                                <span>{{ showAdvanced ? 'Скрыть' : 'Показать' }} дополнительные поля</span>
                                <svg class="w-4 h-4 ml-1 transition-transform" :class="{ 'rotate-180': showAdvanced }"
                                    fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                        d="M19 9l-7 7-7-7" />
                                </svg>
                            </button>
                        </div>

                        <!-- Расширенные поля -->
                        <Transition enter-active-class="transition-all duration-300 ease-out"
                            enter-from-class="transform opacity-0 -translate-y-4"
                            enter-to-class="transform opacity-100 translate-y-0"
                            leave-active-class="transition-all duration-200 ease-in"
                            leave-from-class="transform opacity-100 translate-y-0"
                            leave-to-class="transform opacity-0 -translate-y-4">
                            <div v-if="showAdvanced" class="space-y-5 pt-4 border-t border-stone-100">
                                <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
                                    <!-- Онъёми -->
                                    <div>
                                        <label class="block text-sm font-medium text-text-main mb-1.5 font-jp">
                                            Онъёми
                                            <span class="text-xs text-text-muted ml-1">(через запятую)</span>
                                        </label>
                                        <div class="space-y-2">
                                            <div v-for="(item, index) in form.on" :key="'on-' + index"
                                                class="flex items-center space-x-2">
                                                <input v-model="form.on[index]" type="text"
                                                    class="flex-1 px-3 py-2 border border-stone-200 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent transition-colors font-jp bg-onyomi/30 text-sm"
                                                    placeholder="Онъёми" />
                                                <button v-if="form.on.length > 1" type="button"
                                                    @click="removeArrayItem('on', index)"
                                                    class="p-1.5 text-stone-400 hover:text-danger transition-colors">
                                                    <svg class="w-4 h-4" fill="none" stroke="currentColor"
                                                        viewBox="0 0 24 24">
                                                        <path stroke-linecap="round" stroke-linejoin="round"
                                                            stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                                                    </svg>
                                                </button>
                                            </div>
                                        </div>
                                        <button type="button" @click="addArrayItem('on')"
                                            class="mt-1.5 text-xs text-primary hover:text-primary/80 flex items-center font-medium">
                                            <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor"
                                                viewBox="0 0 24 24">
                                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                                    d="M12 4v16m8-8H4" />
                                            </svg>
                                            Еще вариант
                                        </button>
                                    </div>

                                    <!-- Кунъёми -->
                                    <div>
                                        <label class="block text-sm font-medium text-text-main mb-1.5 font-jp">
                                            Кунъёми
                                            <span class="text-xs text-text-muted ml-1">(через запятую)</span>
                                        </label>
                                        <div class="space-y-2">
                                            <div v-for="(item, index) in form.kun" :key="'kun-' + index"
                                                class="flex items-center space-x-2">
                                                <input v-model="form.kun[index]" type="text"
                                                    class="flex-1 px-3 py-2 border border-stone-200 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent transition-colors font-jp bg-kunyomi/30 text-sm"
                                                    placeholder="Кунъёми" />
                                                <button v-if="form.kun.length > 1" type="button"
                                                    @click="removeArrayItem('kun', index)"
                                                    class="p-1.5 text-stone-400 hover:text-danger transition-colors">
                                                    <svg class="w-4 h-4" fill="none" stroke="currentColor"
                                                        viewBox="0 0 24 24">
                                                        <path stroke-linecap="round" stroke-linejoin="round"
                                                            stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                                                    </svg>
                                                </button>
                                            </div>
                                        </div>
                                        <button type="button" @click="addArrayItem('kun')"
                                            class="mt-1.5 text-xs text-primary hover:text-primary/80 flex items-center font-medium">
                                            <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor"
                                                viewBox="0 0 24 24">
                                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                                    d="M12 4v16m8-8H4" />
                                            </svg>
                                            Еще вариант
                                        </button>
                                    </div>
                                </div>

                                <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
                                    <!-- Примеры на японском -->
                                    <div>
                                        <label class="block text-sm font-medium text-text-main mb-1.5">
                                            Примеры (JP)
                                            <span class="text-xs text-text-muted ml-1">(макс. 3)</span>
                                        </label>
                                        <div class="space-y-2">
                                            <div v-for="(item, index) in form.ex_jp" :key="'ex_jp-' + index"
                                                class="flex items-center space-x-2">
                                                <textarea v-model="form.ex_jp[index]" rows="2"
                                                    class="flex-1 px-3 py-2 border border-stone-200 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent transition-colors resize-none text-sm"
                                                    placeholder="Пример на японском" />
                                                <button v-if="form.ex_jp.length > 1" type="button"
                                                    @click="removeArrayItem('ex_jp', index)"
                                                    class="p-1.5 text-stone-400 hover:text-danger transition-colors self-start">
                                                    <svg class="w-4 h-4" fill="none" stroke="currentColor"
                                                        viewBox="0 0 24 24">
                                                        <path stroke-linecap="round" stroke-linejoin="round"
                                                            stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                                                    </svg>
                                                </button>
                                            </div>
                                        </div>
                                        <button v-if="form.ex_jp.length < 3" type="button"
                                            @click="addArrayItem('ex_jp')"
                                            class="mt-1.5 text-xs text-primary hover:text-primary/80 flex items-center">
                                            <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor"
                                                viewBox="0 0 24 24">
                                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                                    d="M12 4v16m8-8H4" />
                                            </svg>
                                            Еще пример
                                        </button>
                                    </div>

                                    <!-- Переводы примеров -->
                                    <div>
                                        <label class="block text-sm font-medium text-text-main mb-1.5">
                                            Примеры (RU)
                                            <span class="text-xs text-text-muted ml-1">(макс. 3)</span>
                                        </label>
                                        <div class="space-y-2">
                                            <div v-for="(item, index) in form.ex_ru" :key="'ex_ru-' + index"
                                                class="flex items-center space-x-2">
                                                <textarea v-model="form.ex_ru[index]" rows="2"
                                                    class="flex-1 px-3 py-2 border border-stone-200 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent transition-colors resize-none text-sm"
                                                    placeholder="Перевод примера" />
                                                <button v-if="form.ex_ru.length > 1" type="button"
                                                    @click="removeArrayItem('ex_ru', index)"
                                                    class="p-1.5 text-stone-400 hover:text-danger transition-colors self-start">
                                                    <svg class="w-4 h-4" fill="none" stroke="currentColor"
                                                        viewBox="0 0 24 24">
                                                        <path stroke-linecap="round" stroke-linejoin="round"
                                                            stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                                                    </svg>
                                                </button>
                                            </div>
                                        </div>
                                        <button v-if="form.ex_ru.length < 3" type="button"
                                            @click="addArrayItem('ex_ru')"
                                            class="mt-1.5 text-xs text-primary hover:text-primary/80 flex items-center">
                                            <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor"
                                                viewBox="0 0 24 24">
                                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                                    d="M12 4v16m8-8H4" />
                                            </svg>
                                            Еще перевод
                                        </button>
                                    </div>
                                </div>

                                <!-- Теги -->
                                <div>
                                    <label class="block text-sm font-medium text-text-main mb-1.5">
                                        Теги
                                        <span class="text-xs text-text-muted ml-1">(макс. 5)</span>
                                    </label>
                                    <input v-model="tagsInput" type="text" @input="updateTagsFromInput"
                                        @blur="updateTagsFromInput"
                                        class="w-full px-3 py-2 border border-stone-200 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent transition-colors text-sm"
                                        placeholder="Теги через запятую (например: n5, база, дом)" />

                                    <!-- Список тегов -->
                                    <div v-if="form.tags.length > 0" class="mt-2 flex flex-wrap gap-1.5">
                                        <span v-for="(tag, index) in form.tags" :key="'tag-' + index"
                                            class="inline-flex items-center px-2 py-0.5 rounded-full text-xs bg-primary/10 text-primary">
                                            {{ tag }}
                                            <button type="button" @click="removeTag(index)"
                                                class="ml-1.5 text-primary/70 hover:text-primary">
                                                ×
                                            </button>
                                        </span>
                                    </div>
                                </div>
                            </div>
                        </Transition>
                    </form>
                </div>

                <!-- Футер с кнопками -->
                <div class="px-6 py-4 border-t border-stone-100 bg-stone-50 sticky bottom-0">
                    <div class="flex justify-end space-x-3">
                        <button @click="close" type="button"
                            class="px-5 py-2.5 text-text-main hover:text-stone-900 hover:bg-stone-100 rounded-lg transition-colors font-medium">
                            Отмена
                        </button>
                        <button @click="handleSubmit" :disabled="isSubmitting" :class="[
                            'px-5 py-2.5 bg-primary hover:bg-primary/90 text-white font-medium rounded-lg transition-colors',
                            isSubmitting ? 'opacity-50 cursor-not-allowed' : ''
                        ]">
                            <span v-if="isSubmitting">Сохранение...</span>
                            <span v-else>{{ isEditing ? 'Сохранить изменения' : 'Добавить слово' }}</span>
                        </button>
                    </div>
                </div>
            </div>
        </div>
    </Transition>
</template>

<script setup>
import { ref, reactive, watch } from 'vue'
import { useToast } from '@/composables/useToast'

/**
 * Свойства компонента WordFormModal.
 * @typedef {Object} Props
 * @property {Boolean} isOpen - Открыто ли модальное окно.
 * @property {Object|null} word - Объект слова для редактирования (null для создания нового).
 */
const props = defineProps({
    isOpen: {
        type: Boolean,
        default: false
    },
    word: {
        type: Object,
        default: null
    }
})

/**
 * События компонента WordFormModal.
 * @emits close - Событие закрытия модального окна.
 * @emits saved - Событие успешного сохранения слова.
 */
const emit = defineEmits(['close', 'saved'])

const { showSuccess, showError } = useToast()

const isEditing = ref(false)
const isSubmitting = ref(false)
const showAdvanced = ref(false)
const tagsInput = ref('')

const form = reactive({
    jp: [''],
    ru: [''],
    on: [''],
    kun: [''],
    ex_jp: [''],
    ex_ru: [''],
    tags: []
})

const errors = reactive({})

/**
 * Инициализирует состояние формы на основе свойств.
 */
const initForm = () => {
    resetForm()
    if (props.word) {
        isEditing.value = true
        loadWordData(props.word)
    } else {
        isEditing.value = false
    }
}

/**
 * Загружает данные существующего слова в форму.
 * @param {Object} word - Объект слова для загрузки.
 */
const loadWordData = (word) => {
    if (!word) return

    form.jp = Array.isArray(word.jp) && word.jp.length ? [...word.jp] : ['']
    form.ru = Array.isArray(word.ru) && word.ru.length ? [...word.ru] : ['']
    form.on = Array.isArray(word.on) && word.on.length ? [...word.on] : ['']
    form.kun = Array.isArray(word.kun) && word.kun.length ? [...word.kun] : ['']
    form.ex_jp = Array.isArray(word.ex_jp) && word.ex_jp.length ? [...word.ex_jp] : ['']
    form.ex_ru = Array.isArray(word.ex_ru) && word.ex_ru.length ? [...word.ex_ru] : ['']
    form.tags = Array.isArray(word.tags) && word.tags.length ? [...word.tags] : []

    // Обновляем поле ввода тегов
    tagsInput.value = form.tags.join(', ')

    // Автоматически раскрываем дополнительные поля, если в них есть данные
    if ((form.on.length > 1 || (form.on[0] && form.on[0].trim())) ||
        (form.kun.length > 1 || (form.kun[0] && form.kun[0].trim())) ||
        (form.ex_jp.length > 1 || (form.ex_jp[0] && form.ex_jp[0].trim())) ||
        (form.tags.length > 0)) {
        showAdvanced.value = true
    }
}

/**
 * Сбрасывает форму к начальному состоянию.
 */
const resetForm = () => {
    form.jp = ['']
    form.ru = ['']
    form.on = ['']
    form.kun = ['']
    form.ex_jp = ['']
    form.ex_ru = ['']
    form.tags = []
    tagsInput.value = ''
    showAdvanced.value = false
    Object.keys(errors).forEach(key => delete errors[key])
}

// Инициализация формы при открытии
watch(() => props.isOpen, (isOpen) => {
    if (isOpen) {
        initForm()
    }
}, { immediate: true })

// Отслеживание изменений слова при открытой форме
watch(() => props.word, (newWord) => {
    if (props.isOpen) {
        initForm()
    }
})

/**
 * Валидирует данные формы.
 * @returns {Boolean} - True, если данные валидны, иначе false.
 */
const validateForm = () => {
    let isValid = true
    Object.keys(errors).forEach(key => delete errors[key])

    // Проверка японских написаний
    const validJp = form.jp.filter(item => item.trim() !== '')
    if (validJp.length === 0) {
        errors.jp = 'Необходимо указать хотя бы одно японское написание'
        isValid = false
    } else {
        // Проверка на дубликаты
        const uniqueJp = [...new Set(validJp)]
        if (uniqueJp.length !== validJp.length) {
            errors.jp = 'Японские написания не должны повторяться'
            isValid = false
        }
    }

    // Проверка русских переводов
    const validRu = form.ru.filter(item => item.trim() !== '')
    if (validRu.length === 0) {
        errors.ru = 'Необходимо указать хотя бы один русский перевод'
        isValid = false
    }

    return isValid
}

/**
 * Добавляет новый элемент в поле-массив.
 * @param {String} field - Имя поля (например, 'jp', 'ru').
 */
const addArrayItem = (field) => {
    // Проверка лимитов для массивов
    if (field === 'ex_jp' && form.ex_jp.length >= 3) return
    if (field === 'ex_ru' && form.ex_ru.length >= 3) return
    if (field === 'tags' && form.tags.length >= 5) return

    if (field === 'jp' || field === 'ru' || field === 'on' || field === 'kun') {
        form[field].push('')
    } else if (field === 'ex_jp' || field === 'ex_ru') {
        form[field].push('')
    }
}

/**
 * Удаляет элемент из поля-массива.
 * @param {String} field - Имя поля.
 * @param {Number} index - Индекс удаляемого элемента.
 */
const removeArrayItem = (field, index) => {
    form[field].splice(index, 1)
    // Если массив пуст, добавляем пустой элемент для полей чтений
    if ((field === 'on' || field === 'kun') && form[field].length === 0) {
        form[field].push('')
    }
}

/**
 * Обновляет теги из поля ввода.
 * Обрабатывает как ввод (парсинг), так и потерю фокуса (форматирование).
 * @param {Event} [event] - DOM событие.
 */
const updateTagsFromInput = (event) => {
    const tags = tagsInput.value
        .split(',')
        .map(tag => tag.trim())
        .filter(tag => tag !== '')
        .slice(0, 5) // Ограничение на 5 тегов

    form.tags = tags

    // Форматируем отображаемое значение только при потере фокуса для удобства ввода
    if (event && event.type === 'blur') {
        tagsInput.value = tags.join(', ')
    }
}

/**
 * Удаляет конкретный тег.
 * @param {Number} index - Индекс удаляемого тега.
 */
const removeTag = (index) => {
    form.tags.splice(index, 1)
    tagsInput.value = form.tags.join(', ')
}

/**
 * Обрабатывает отправку формы.
 */
const handleSubmit = async () => {
    if (!validateForm()) return

    isSubmitting.value = true

    try {
        // Подготавливаем данные для отправки
        const data = {
            jp: form.jp.filter(item => item.trim() !== ''),
            ru: form.ru.filter(item => item.trim() !== ''),
            on: form.on[0] ? form.on.filter(item => item.trim() !== '') : undefined,
            kun: form.kun[0] ? form.kun.filter(item => item.trim() !== '') : undefined,
            ex_jp: form.ex_jp[0] ? form.ex_jp.filter(item => item.trim() !== '') : undefined,
            ex_ru: form.ex_ru[0] ? form.ex_ru.filter(item => item.trim() !== '') : undefined,
            tags: form.tags.length > 0 ? form.tags : undefined
        }

        // Отправляем данные через emit вместе с флагом редактирования
        emit('saved', { data, isEditing: isEditing.value, word: props.word })

    } catch (error) {
        showError('Ошибка при сохранении слова')
    } finally {
        isSubmitting.value = false
    }
}

/**
 * Закрывает модальное окно.
 */
const close = () => {
    resetForm()
    emit('close')
}
</script>

<style scoped>
.modal-enter-active,
.modal-leave-active {
    transition: opacity 0.3s ease;
}

.modal-enter-from,
.modal-leave-to {
    opacity: 0;
}
</style>