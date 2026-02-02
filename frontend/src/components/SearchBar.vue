<template>
    <div class="mb-8">
        <!-- Основной поиск -->
        <div class="relative mb-4 group">
            <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                <svg class="h-5 w-5 text-gray-400 group-focus-within:text-primary transition-colors duration-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                </svg>
            </div>
            <input v-model="searchQuery" type="text" placeholder="Поиск слова (русский или японский)..."
                class="w-full pl-11 pr-10 py-3.5 bg-surface border-none shadow-sm rounded-xl focus:ring-2 focus:ring-primary/50 focus:shadow-lg transition-all duration-300 placeholder-gray-400 text-text-main"
                @input="handleSearchInput" />

            <!-- Индикатор загрузки -->
            <div v-if="loading" class="absolute inset-y-0 right-0 pr-4 flex items-center">
                <div class="animate-spin rounded-full h-5 w-5 border-b-2 border-primary"></div>
            </div>

            <!-- Кнопка очистки -->
            <button v-if="searchQuery" @click="clearSearch"
                class="absolute inset-y-0 right-0 pr-4 flex items-center text-gray-400 hover:text-primary transition-colors"
                type="button">
                <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
            </button>
        </div>

        <!-- Кнопка расширенного поиска -->
        <div class="mb-4 flex justify-end">
            <button @click="showAdvanced = !showAdvanced"
                class="text-sm font-medium text-text-muted hover:text-primary flex items-center transition-colors px-3 py-1 rounded-md hover:bg-gray-100/50" type="button">
                <span>{{ showAdvanced ? 'Скрыть фильтры' : 'Расширенный поиск' }}</span>
                <svg class="w-4 h-4 ml-1 transition-transform duration-200" :class="{ 'rotate-180': showAdvanced }"
                    fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                </svg>
            </button>
        </div>

        <!-- Расширенный поиск -->
        <Transition enter-active-class="transition-all duration-300 ease-out"
            enter-from-class="transform opacity-0 -translate-y-4" enter-to-class="transform opacity-100 translate-y-0"
            leave-active-class="transition-all duration-200 ease-in"
            leave-from-class="transform opacity-100 translate-y-0" leave-to-class="transform opacity-0 -translate-y-4">
            <div v-if="showAdvanced" class="p-6 bg-surface rounded-xl shadow-sm mb-6 border border-gray-100">
                <!-- Поля для расширенного поиска -->
                <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-4">
                    <!-- Теги -->
                    <div>
                        <label class="block text-sm font-medium text-text-main mb-2">
                            Теги
                        </label>
                        <input v-model="advancedFilters.tags" type="text"
                            class="w-full px-4 py-2 bg-background border-none rounded-lg focus:ring-2 focus:ring-primary/50 transition-colors placeholder-gray-400"
                            placeholder="Например: n5, verbs" />
                    </div>

                    <!-- Онъёми -->
                    <div>
                        <label class="block text-sm font-medium text-text-main mb-2">
                            Онъёми (катакана)
                        </label>
                        <input v-model="advancedFilters.on" type="text"
                            class="w-full px-4 py-2 bg-background border-none rounded-lg focus:ring-2 focus:ring-primary/50 transition-colors placeholder-gray-400"
                            placeholder="Например: キ, シ" />
                    </div>

                    <!-- Кунъёми -->
                    <div>
                        <label class="block text-sm font-medium text-text-main mb-2">
                            Кунъёми (хирагана)
                        </label>
                        <input v-model="advancedFilters.kun" type="text"
                            class="w-full px-4 py-2 bg-background border-none rounded-lg focus:ring-2 focus:ring-primary/50 transition-colors placeholder-gray-400"
                            placeholder="Например: つくえ" />
                    </div>
                </div>

                <!-- Активные фильтры -->
                <div v-if="hasActiveFilters" class="pt-4 border-t border-gray-100 flex flex-wrap gap-2">
                    <span v-if="activeTags.length > 0"
                        class="inline-flex items-center px-3 py-1 text-xs font-medium bg-gray-100 text-gray-700 rounded-full">
                        Теги: {{ activeTags.join(', ') }}
                        <button @click="clearTagFilter('tags')"
                            class="ml-2 text-gray-500 hover:text-red-500 transition-colors text-sm">×</button>
                    </span>

                    <span v-if="activeOn.length > 0"
                        class="inline-flex items-center px-3 py-1 text-xs font-medium bg-onyomi text-amber-900 rounded-full">
                        On: {{ activeOn.join(', ') }}
                        <button @click="clearTagFilter('on')"
                            class="ml-2 text-amber-700 hover:text-red-600 transition-colors text-sm">×</button>
                    </span>

                    <span v-if="activeKun.length > 0"
                        class="inline-flex items-center px-3 py-1 text-xs bg-green-100 text-green-800 rounded-full">
                        Кунъёми: {{ activeKun.join(', ') }}
                        <button @click="clearTagFilter('kun')"
                            class="ml-2 text-green-600 hover:text-green-800 text-sm">×</button>
                    </span>
                </div>

                <!-- Кнопки расширенного поиска -->
                <div class="mt-4 flex justify-end space-x-3">
                    <button @click="clearAdvancedFilters" type="button"
                        class="px-4 py-2 text-sm text-gray-700 hover:text-gray-900 hover:bg-gray-100 rounded-md transition-colors">
                        Сбросить
                    </button>
                    <button @click="applyAdvancedSearch" type="button"
                        class="px-4 py-2 text-sm font-medium text-white bg-primary hover:bg-primary/90 rounded-md transition-colors">
                        Применить
                    </button>
                </div>
            </div>
        </Transition>
    </div>
</template>

<script setup>
import { ref, reactive, watch, computed } from 'vue'
import { debounce } from 'lodash-es'

const emit = defineEmits(['search', 'advanced-search'])

const showAdvanced = ref(false)
const searchQuery = ref('')
const loading = ref(false)

const advancedFilters = reactive({
    tags: '',
    on: '',
    kun: ''
})

// Дебаунс для основного поиска
const debouncedSearch = debounce(() => {
    loading.value = true
    emit('search', searchQuery.value.trim())
    setTimeout(() => {
        loading.value = false
    }, 300)
}, 800)

const handleSearchInput = () => {
    debouncedSearch()
}

const clearSearch = () => {
    searchQuery.value = ''
    emit('search', '')
}

const applyAdvancedSearch = () => {
    const params = {}

    if (advancedFilters.tags.trim()) {
        const tagsArray = advancedFilters.tags
            .split(',')
            .map(tag => tag.trim())
            .filter(tag => tag !== '')

        params.tags = tagsArray.join(',')
    }

    if (advancedFilters.on.trim()) {
        const onArray = advancedFilters.on
            .split(',')
            .map(item => item.trim())
            .filter(item => item !== '')

        params.on = onArray.join(',')
    }

    if (advancedFilters.kun.trim()) {
        const kunArray = advancedFilters.kun
            .split(',')
            .map(item => item.trim())
            .filter(item => item !== '')

        params.kun = kunArray.join(',')
    }

    console.log('Search params:', params)
    emit('advanced-search', params)
}

const clearAdvancedFilters = () => {
    advancedFilters.tags = ''
    advancedFilters.on = ''
    advancedFilters.kun = ''
    emit('advanced-search', {})
}

// При изменении основного поиска сбрасываем расширенные фильтры
watch(searchQuery, (newValue) => {
    if (newValue.trim() === '') {
        clearAdvancedFilters()
    }
})

// Свойство для активных фильтров
const activeTags = computed(() => {
    return advancedFilters.tags
        .split(',')
        .map(tag => tag.trim())
        .filter(tag => tag !== '')
})

const activeOn = computed(() => {
    return advancedFilters.on
        .split(',')
        .map(item => item.trim())
        .filter(item => item !== '')
})

const activeKun = computed(() => {
    return advancedFilters.kun
        .split(',')
        .map(item => item.trim())
        .filter(item => item !== '')
})

const hasActiveFilters = computed(() => {
    return activeTags.value.length > 0 ||
        activeOn.value.length > 0 ||
        activeKun.value.length > 0
})

// Функция очистки конкретного фильтра
const clearTagFilter = (type) => {
    if (type === 'tags') {
        advancedFilters.tags = ''
    } else if (type === 'on') {
        advancedFilters.on = ''
    } else if (type === 'kun') {
        advancedFilters.kun = ''
    }

    applyAdvancedSearch()
}

clearAdvancedFilters()
</script>

<style>
@import './SearchBar.css';
</style>