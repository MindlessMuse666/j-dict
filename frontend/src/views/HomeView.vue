<template>
  <div class="max-w-5xl mx-auto">
    <!-- Панель управления -->
    <div class="mb-8">
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-6">
        <div>
          <h2 class="text-3xl font-bold text-text-main font-jp">Мой словарь</h2>
          <!-- Плавное изменение счетчика слов -->
          <Transition name="fade-slide" mode="out-in">
            <p v-if="wordsStore.words.length" key="count" class="text-sm text-text-muted mt-1 font-medium">
              Всего слов: {{ wordsStore.words.length }}
            </p>
            <div v-else key="empty" class="h-5"></div>
          </Transition>
        </div>

        <div class="flex items-center gap-3">
          <!-- Переключатель компактного вида -->
          <div class="hidden sm:flex items-center bg-stone-100/80 p-1 rounded-xl shadow-inner mr-2">
            <button @click="isCompactView = false" :class="[
              'p-2 rounded-lg transition-all duration-200 flex items-center justify-center',
              !isCompactView ? 'bg-white text-primary shadow-sm' : 'text-stone-400 hover:text-stone-600'
            ]" title="Полный вид">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
              </svg>
            </button>
            <button @click="isCompactView = true" :class="[
              'p-2 rounded-lg transition-all duration-200 flex items-center justify-center',
              isCompactView ? 'bg-white text-primary shadow-sm' : 'text-stone-400 hover:text-stone-600'
            ]" title="Компактный вид">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 10h16M4 14h16M4 18h16" />
              </svg>
            </button>
          </div>

          <!-- Кнопка импорта -->
          <button @click="showImportModal = true"
            class="w-32 justify-center px-4 py-2 bg-stone-100 hover:bg-stone-200 text-text-main font-medium rounded-xl transition-all duration-200 hover:shadow-md flex items-center transform hover:-translate-y-0.5 border border-stone-200">
            <svg class="w-5 h-5 mr-2 text-secondary" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M9 19l3 3m0 0l3-3m-3 3V10" />
            </svg>
            Импорт
          </button>

          <!-- Кнопка добавления -->
          <button @click="showAddModal = true"
            class="w-32 justify-center px-4 py-2 bg-primary hover:bg-primary/90 text-white font-medium rounded-xl transition-all duration-200 hover:shadow-lg shadow-md flex items-center transform hover:-translate-y-0.5">
            <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
            </svg>
            Добавить
          </button>
        </div>
      </div>

      <!-- Поисковая панель -->
      <SearchBar @search="handleSearch" @advanced-search="handleAdvancedSearch" />
    </div>

    <!-- Сообщение об ошибке -->
    <Transition name="fade">
      <div v-if="wordsStore.error" class="mb-6 p-4 bg-red-50 border border-red-100 rounded-xl shadow-sm">
        <div class="flex items-center">
          <svg class="w-5 h-5 text-danger mr-3" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd"
              d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
              clip-rule="evenodd" />
          </svg>
          <span class="text-danger font-medium">{{ wordsStore.error }}</span>
        </div>
      </div>
    </Transition>

    <!-- Индикатор загрузки при первоначальной загрузке -->
    <div v-if="wordsStore.loading && wordsStore.words.length === 0" class="space-y-4">
      <WordSkeleton v-for="n in 5" :key="'skeleton-' + n" :is-compact="isCompactView" />
    </div>

    <!-- Сообщение о пустом списке -->
    <Transition name="fade">
      <div v-if="!wordsStore.loading && (!wordsStore.words || wordsStore.words.length === 0)">
        <EmptyState
          :title="isSearchActive ? 'Ничего не найдено' : 'Словарь пуст'"
          :description="emptyMessage"
        >
          <template #action v-if="!isSearchActive">
            <button @click="showAddModal = true"
              class="px-6 py-3 bg-primary hover:bg-primary/90 text-white font-medium rounded-lg transition-all duration-200 hover:shadow-lg shadow-md flex items-center mx-auto transform hover:-translate-y-0.5">
              <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
              </svg>
              Добавить первое слово
            </button>
          </template>
          <template #action v-else>
            <button @click="handleClearSearch"
              class="px-6 py-3 bg-white border border-gray-300 text-gray-700 font-medium rounded-lg hover:bg-gray-50 transition-all duration-200">
              Сбросить поиск
            </button>
          </template>
        </EmptyState>
      </div>
    </Transition>

    <!-- Список слов -->
    <TransitionGroup name="word-list" tag="div" class="words-list space-y-4">
      <WordCard v-for="word in wordsStore.words" :key="word.id" :word="word" :is-compact-view="isCompactView"
        @edit="handleEditWord(word)" @delete="handleWordDeleted(word.id)" />
    </TransitionGroup>

    <!-- Индикатор загрузки при подгрузке -->
    <Transition name="fade">
      <div v-if="wordsStore.loading && wordsStore.words.length > 0" class="mt-6">
        <WordSkeleton v-for="n in 3" :key="'skeleton-more-' + n" :is-compact="isCompactView" />
      </div>
    </Transition>

    <!-- Кнопка "Загрузить еще" -->
    <Transition name="fade">
      <div v-if="!wordsStore.loading && wordsStore.hasMore" class="mt-6 text-center">
        <button @click="loadMoreWords"
          class="px-6 py-3 bg-gray-100 hover:bg-gray-200 text-gray-800 font-medium rounded-lg transition-all duration-200 hover:shadow-sm">
          Загрузить еще
        </button>
      </div>
    </Transition>

    <!-- Сообщение о конце списка -->
    <Transition name="fade">
      <div v-if="!wordsStore.hasMore && wordsStore.words.length > 0" class="mt-6 pt-6 border-t text-center">
        <p class="text-gray-500">Вы просмотрели все слова в словаре</p>
      </div>
    </Transition>

    <!-- Модальные окна -->
    <WordFormModal v-if="showAddModal" :word="selectedWord" @close="closeModal" @saved="handleWordSaved"
      :is-open="showAddModal" />

    <WordFormModal v-if="showEditModal" :word="selectedWord" @close="closeModal" @saved="handleWordSaved"
      :is-open="showEditModal" />

    <ImportModal v-if="showImportModal" @close="showImportModal = false" @imported="handleImported"
      :is-open="showImportModal" />
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { useWordsStore } from '@/stores/words'
import { useToast } from '@/composables/useToast'
import { useHotkeys } from '@/composables/useHotkeys'
import WordCard from '@/components/WordCard.vue'
import WordFormModal from '@/components/WordFormModal.vue'
import SearchBar from '@/components/SearchBar.vue'
import WordSkeleton from '@/components/WordSkeleton.vue'
import ImportModal from '@/components/ImportModal.vue'
import EmptyState from '@/components/EmptyState.vue'

const wordsStore = useWordsStore()
const { showSuccess, showError } = useToast()

const searchBarRef = ref(null)
const isCompactView = ref(true)
const showAddModal = ref(false)
const showImportModal = ref(false)
const showEditModal = ref(false)
const selectedWord = ref(null)
const currentSearchQuery = ref('')
const isSearchActive = ref(false)
const importedWords = ref([])

// Компьютед свойство для отображения правильного сообщения
const emptyMessage = computed(() => {
  if (wordsStore.loading) return 'Загрузка...'
  if (isSearchActive.value) {
    return currentSearchQuery.value
      ? `По запросу "${currentSearchQuery.value}" ничего не найдено. Попробуйте изменить запрос или проверить написание.`
      : 'Ничего не найдено по заданным критериям. Попробуйте смягчить фильтры.'
  }
  return 'Ваш словарь пока пуст. Добавьте первые слова, чтобы начать изучение.'
})

// Бесконечная пагинация при скролле
const handleScroll = () => {
  if (isSearchActive.value || wordsStore.loading || !wordsStore.hasMore) {
    return
  }

  const scrollPosition = window.innerHeight + window.scrollY
  const pageHeight = document.documentElement.offsetHeight - 100

  if (scrollPosition >= pageHeight) {
    loadMoreWords()
  }
}

const handleImported = (importResult) => {
  showSuccess(`Успешно импортировано ${importResult.imported_count || 0} слов`)
  // Обновляем список слов
  wordsStore.fetchWords(20, true)
}

const loadMoreWords = () => {
  wordsStore.fetchWords()
}

const handleSearch = async (query) => {
  currentSearchQuery.value = query
  isSearchActive.value = query.trim().length > 0

  if (query.trim()) {
    await wordsStore.searchWords({ q: query })
  } else {
    isSearchActive.value = false
    wordsStore.fetchWords(20, true)
  }
}

const handleAdvancedSearch = async (params) => {
  const hasParams = Object.values(params).some(value =>
    (Array.isArray(value) && value.length > 0) ||
    (typeof value === 'string' && value.trim().length > 0)
  )

  isSearchActive.value = hasParams
  currentSearchQuery.value = ''

  if (hasParams) {
    await wordsStore.searchWords(params)
  } else {
    wordsStore.fetchWords(20, true)
  }
}

const handleClearSearch = () => {
  if (searchBarRef.value) {
    searchBarRef.value.clearSearch()
    searchBarRef.value.clearAdvancedFilters()
  }
  // Fallback if ref is not available for some reason
  currentSearchQuery.value = ''
  isSearchActive.value = false
  if (!searchBarRef.value) {
    wordsStore.fetchWords(20, true)
  }
}

const handleEditWord = (word) => {
  selectedWord.value = word
  showEditModal.value = true
}

const handleWordDeleted = (wordId) => {
  // Уведомление уже показывается в WordCard
  // Здесь можно добавить дополнительную логику, если нужно
}

const handleWordSaved = async (wordData) => {
  try {
    let result

    if (selectedWord.value) {
      result = await wordsStore.updateWord(selectedWord.value.id, wordData.data)
    } else {
      result = await wordsStore.createWord(wordData.data)
    }

    if (result.success) {
      showSuccess(
        selectedWord.value
          ? 'Слово успешно обновлено'
          : 'Слово успешно добавлено'
      )
      closeModal()
    } else {
      showError(result.error)
    }
  } catch (error) {
    showError('Ошибка при сохранении слова')
  }
}

const closeModal = () => {
  showAddModal.value = false
  showEditModal.value = false
  selectedWord.value = null
}

// Инициализация
onMounted(() => {
  useHotkeys({
    onNewWord: () => {
      if (!showAddModal.value && !showEditModal.value) {
        showAddModal.value = true
      }
    },
    onCloseModals: () => {
      if (showAddModal.value || showEditModal.value || showImportModal.value) {
        showAddModal.value = false
        showEditModal.value = false
        showImportModal.value = false
      }
    }
  })

  if (wordsStore.words.length === 0) {
    wordsStore.fetchWords()
  }

  window.addEventListener('scroll', handleScroll)
})

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
})
</script>

<style scoped>
/* Анимации для списка слов */
.word-list-enter-active,
.word-list-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.word-list-enter-from {
  opacity: 0;
  transform: translateY(10px);
}

.word-list-leave-to {
  opacity: 0;
  transform: translateX(-100%);
}

.word-list-leave-active {
  position: absolute;
  width: 100%;
}

.word-list-move {
  transition: transform 0.3s ease;
}

/* Плавные переходы */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.fade-slide-enter-active,
.fade-slide-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.fade-slide-enter-from {
  opacity: 0;
  transform: translateY(-10px);
}

.fade-slide-leave-to {
  opacity: 0;
  transform: translateY(10px);
}

/* Улучшение для мобильных устройств */
@media (max-width: 640px) {
  .mobile-stack {
    flex-direction: column;
    align-items: stretch;
  }

  .mobile-full-width {
    width: 100%;
  }

  .mobile-text-center {
    text-align: center;
  }

  .mobile-mt-2 {
    margin-top: 0.5rem;
  }
}
</style>