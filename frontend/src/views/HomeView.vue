<template>
  <div class="max-w-6xl mx-auto">
    <!-- Панель управления -->
    <div class="mb-6">
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-4">
        <div>
          <h2 class="text-2xl font-bold text-gray-800">Мой словарь</h2>
          <!-- Плавное изменение счетчика слов -->
          <Transition name="fade-slide" mode="out-in">
            <p v-if="wordsStore.words.length" key="count" class="text-sm text-gray-600 mt-1">
              Всего слов: {{ wordsStore.words.length }}
            </p>
            <div v-else key="empty" class="h-5"></div>
          </Transition>
        </div>

        <div class="flex items-center space-x-4">
          <!-- Переключатель компактного вида -->
          <div class="hidden sm:flex items-center">
            <span class="text-sm text-gray-700 mr-2">Компактный вид:</span>
            <button @click="isCompactView = !isCompactView" :class="[
              'relative inline-flex h-6 w-11 items-center rounded-full transition-colors duration-200',
              isCompactView ? 'bg-primary' : 'bg-gray-200'
            ]">
              <span :class="[
                'inline-block h-4 w-4 transform rounded-full bg-white transition-transform duration-200',
                isCompactView ? 'translate-x-6' : 'translate-x-1'
              ]" />
            </button>
          </div>

          <!-- Кнопка импорта -->
          <button @click="showImportModal = true"
            class="px-4 py-2 bg-secondary hover:bg-secondary/90 text-white font-medium rounded-lg transition-all duration-200 hover:shadow-md flex items-center">
            <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M9 19l3 3m0 0l3-3m-3 3V10" />
            </svg>
            Импорт CSV
          </button>

          <!-- Кнопка добавления -->
          <button @click="showAddModal = true"
            class="px-4 py-2 bg-primary hover:bg-primary/90 text-white font-medium rounded-lg transition-all duration-200 hover:shadow-md flex items-center">
            <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
            </svg>
            Добавить слово
          </button>
        </div>
      </div>

      <!-- Поисковая панель -->
      <SearchBar @search="handleSearch" @advanced-search="handleAdvancedSearch" />
    </div>

    <!-- Сообщение об ошибке -->
    <Transition name="fade">
      <div v-if="wordsStore.error" class="mb-4 p-4 bg-red-50 border border-red-200 rounded-lg">
        <div class="flex items-center">
          <svg class="w-5 h-5 text-red-500 mr-3" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd"
              d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
              clip-rule="evenodd" />
          </svg>
          <span class="text-red-800">{{ wordsStore.error }}</span>
        </div>
      </div>
    </Transition>

    <!-- Индикатор загрузки при первоначальной загрузке -->
    <div v-if="wordsStore.loading && wordsStore.words.length === 0" class="space-y-4">
      <WordSkeleton v-for="n in 5" :key="'skeleton-' + n" />
    </div>

    <!-- Сообщение о пустом списке -->
    <Transition name="fade">
      <div v-if="!wordsStore.loading && (!wordsStore.words || wordsStore.words.length === 0)"
        class="text-center py-12">
        <div class="text-gray-400 mb-4">
          <svg class="w-16 h-16 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1"
              d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.746 0 3.332.477 4.5 1.253v13C19.832 18.477 18.246 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
          </svg>
        </div>
        <h3 class="text-lg font-medium text-gray-700 mb-2">
          {{ emptyMessage }}
        </h3>
        <p class="text-gray-500 mb-4" v-if="!isSearchActive">
          Добавьте первое слово, нажав кнопку "Добавить слово"
        </p>
        <p class="text-gray-500 mb-4" v-else>
          Попробуйте изменить параметры поиска
        </p>
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

const wordsStore = useWordsStore()
const { showSuccess, showError } = useToast()

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
      ? `По запросу "${currentSearchQuery.value}" ничего не найдено`
      : 'Ничего не найдено'
  }
  return 'Словарь пуст'
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
  showSuccess(`Успешно импортировано ${importResult.added_count || 0} слов`)
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

  await wordsStore.searchWords(params)
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