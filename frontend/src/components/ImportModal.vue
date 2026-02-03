<template>
    <Transition name="modal">
        <div v-if="isOpen" class="fixed inset-0 z-50 flex items-center justify-center">
            <!-- Фон с затемнением -->
            <div class="fixed inset-0 bg-black bg-opacity-50" @click="close"></div>

            <!-- Модальное окно -->
            <div class="relative z-50 bg-white rounded-xl shadow-xl w-full max-w-2xl mx-4 max-h-[90vh] overflow-hidden">
                <!-- Заголовок -->
                <div class="px-6 py-4 border-b border-stone-100 bg-stone-50">
                    <div class="flex items-center justify-between">
                        <h3 class="text-lg font-semibold text-text-main font-jp">Импорт CSV</h3>
                        <button @click="close" class="text-stone-400 hover:text-stone-600 transition-colors p-1">
                            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                    d="M6 18L18 6M6 6l12 12" />
                            </svg>
                        </button>
                    </div>
                </div>

                <!-- Контент -->
                <div class="overflow-y-auto max-h-[calc(90vh-8rem)] p-6">
                    <!-- Инструкция -->
                    <div class="mb-6 p-4 bg-stone-50 border border-stone-200 rounded-lg">
                        <h4 class="font-medium text-text-main mb-2 font-jp">Формат CSV файла:</h4>
                        <pre class="text-sm text-text-muted whitespace-pre-wrap font-mono bg-white p-2 rounded border border-stone-100">jp,ru,on,kun,ex_jp,ex_ru,tags
机,стол,キ,つくえ,机の上に本があります,На столе лежит книга,n5,мебель
本,книга,ホン,ほん,本を読みます,Читаю книгу,n5</pre>
                        <p class="text-sm text-text-muted mt-2">
                            • Разделитель: запятая<br>
                            • Первая строка - заголовки<br>
                            • Поля с массивами (теги) разделяйте точкой с запятой<br>
                            • Пустые значения оставляйте пустыми
                        </p>
                    </div>

                    <!-- Загрузка файла -->
                    <div class="mb-6">
                        <label class="block text-sm font-medium text-text-main mb-3">
                            Выберите CSV файл
                        </label>
                        <div @dragover="handleDragOver" @drop="handleDrop" @dragleave="isDragging = false" :class="[
                            'border-2 border-dashed rounded-lg p-8 text-center transition-all duration-200 cursor-pointer',
                            isDragging ? 'border-primary bg-primary/5 scale-[1.02]' :
                                selectedFile ? 'border-primary/50 bg-primary/5' :
                                    'border-stone-300 hover:border-primary/50 hover:bg-stone-50'
                        ]" @click="$refs.fileInput.click()">

                            <svg :class="[
                                'w-12 h-12 mx-auto mb-4 transition-colors',
                                selectedFile ? 'text-primary' :
                                    isDragging ? 'text-primary' : 'text-stone-400'
                            ]" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                    d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M9 19l3 3m0 0l3-3m-3 3V10" />
                            </svg>

                            <p class="text-text-main mb-2">
                                <span :class="[
                                    'font-medium',
                                    selectedFile ? 'text-primary' :
                                        isDragging ? 'text-primary' : 'text-primary'
                                ]">
                                    {{ selectedFile ? 'Файл выбран' : 'Нажмите для загрузки' }}
                                </span>
                                {{ selectedFile ? '' : 'или перетащите файл' }}
                            </p>

                            <p class="text-sm text-text-muted">
                                {{ selectedFile ?
                                    `${selectedFile.name} (${(selectedFile.size / 1024).toFixed(1)} KB)` :
                                'Поддерживаются файлы CSV до 10MB'
                                }}
                            </p>

                            <input ref="fileInput" type="file" accept=".csv,text/csv,application/vnd.ms-excel"
                                @change="handleFileSelect" class="hidden" />
                        </div>

                        <!-- Выбранный файл -->
                        <div v-if="selectedFile" class="mt-4 p-4 bg-stone-50 rounded-lg border border-stone-100">
                            <div class="flex items-center justify-between">
                                <div class="flex items-center">
                                    <svg class="w-5 h-5 text-stone-400 mr-3" fill="none" stroke="currentColor"
                                        viewBox="0 0 24 24">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                            d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                                    </svg>
                                    <div>
                                        <p class="font-medium text-text-main">{{ selectedFile.name }}</p>
                                        <p class="text-sm text-text-muted">{{ (selectedFile.size / 1024).toFixed(2) }} KB
                                        </p>
                                    </div>
                                </div>
                                <button @click="selectedFile = null" class="text-stone-400 hover:text-stone-600">
                                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                            d="M6 18L18 6M6 6l12 12" />
                                    </svg>
                                </button>
                            </div>
                        </div>
                    </div>

                    <!-- Статус импорта -->
                    <div v-if="importStore.importResult"
                        class="mb-6 p-4 bg-green-50 border border-green-200 rounded-lg">
                        <div class="flex items-start">
                            <svg class="w-5 h-5 text-green-500 mt-0.5 mr-3" fill="currentColor" viewBox="0 0 20 20">
                                <path fill-rule="evenodd"
                                    d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                                    clip-rule="evenodd" />
                            </svg>
                            <div>
                                <h4 class="font-medium text-green-800 mb-1">Импорт успешно завершен</h4>
                                <p class="text-sm text-green-700">
                                    Добавлено слов: {{ importStore.importResult.imported_count || 0 }}<br>
                                    <span v-if="importStore.importResult.failed_count > 0" class="text-red-600">
                                        Ошибок: {{ importStore.importResult.failed_count }}
                                    </span>
                                </p>
                                <div v-if="importStore.importResult.errors && importStore.importResult.errors.length > 0" 
                                     class="mt-2 text-xs text-red-600 max-h-32 overflow-y-auto bg-red-50 p-2 rounded border border-red-100">
                                    <ul class="list-disc list-inside">
                                        <li v-for="(err, idx) in importStore.importResult.errors" :key="idx">
                                            {{ err }}
                                        </li>
                                    </ul>
                                </div>
                            </div>
                        </div>
                    </div>

                    <!-- Ошибка импорта -->
                    <div v-if="importStore.error" class="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg">
                        <div class="flex items-start">
                            <svg class="w-5 h-5 text-red-500 mt-0.5 mr-3" fill="currentColor" viewBox="0 0 20 20">
                                <path fill-rule="evenodd"
                                    d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
                                    clip-rule="evenodd" />
                            </svg>
                            <div>
                                <h4 class="font-medium text-red-800 mb-1">Ошибка при импорте</h4>
                                <p class="text-sm text-red-700">{{ importStore.error }}</p>
                            </div>
                        </div>
                    </div>
                </div>

                <!-- Футер -->
                <div class="px-6 py-4 border-t border-stone-100 bg-stone-50">
                    <div class="flex justify-end space-x-3">
                        <button @click="close" type="button"
                            class="px-5 py-2.5 text-text-muted hover:text-text-main hover:bg-stone-100 rounded-lg transition-colors font-medium">
                            Отмена
                        </button>
                        <button @click="handleImport" :disabled="!selectedFile || importStore.isImporting" :class="[
                            'px-5 py-2.5 bg-primary hover:bg-primary/90 text-white font-medium rounded-lg transition-colors',
                            (!selectedFile || importStore.isImporting) ? 'opacity-50 cursor-not-allowed' : ''
                        ]">
                            <span v-if="importStore.isImporting">
                                <svg class="w-4 h-4 inline-block animate-spin mr-2" fill="none" viewBox="0 0 24 24">
                                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor"
                                        stroke-width="4"></circle>
                                    <path class="opacity-75" fill="currentColor"
                                        d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z">
                                    </path>
                                </svg>
                                Импорт...
                            </span>
                            <span v-else>Импортировать</span>
                        </button>
                    </div>
                </div>
            </div>
        </div>
    </Transition>
</template>

<script setup>
import { ref } from 'vue'
import { useImportStore } from '@/stores/import'
import { useToast } from '@/composables/useToast'

const props = defineProps({
    isOpen: {
        type: Boolean,
        default: false
    }
})

const emit = defineEmits(['close', 'imported'])

const importStore = useImportStore()
const { showSuccess, showError } = useToast()

const selectedFile = ref(null)
const isDragging = ref(false)

const handleFileSelect = (event) => {
    const file = event.target.files[0]
    if (file && isValidCSVFile(file)) {
        selectedFile.value = file
        // Сбрасываем значение input, чтобы можно было выбрать тот же файл снова
        event.target.value = ''
    } else {
        showError('Пожалуйста, выберите файл с расширением .csv')
        // Сбрасываем значение input
        event.target.value = ''
    }
}

const handleDragOver = (event) => {
    isDragging.value = true
    event.preventDefault()
}

const handleDrop = (event) => {
    isDragging.value = false
    event.preventDefault()

    const file = event.dataTransfer.files[0]
    if (file && isValidCSVFile(file)) {
        selectedFile.value = file
    } else {
        showError('Пожалуйста, перетащите файл с расширением .csv')
    }
}

const isValidCSVFile = (file) => {
    // Проверяем по расширению, а не по MIME-типу, так как он часто не определяется
    const fileName = file.name.toLowerCase()
    const isValidExtension = fileName.endsWith('.csv')

    // Дополнительная проверка размера (максимум 10MB)
    const isValidSize = file.size <= 10 * 1024 * 1024 // 10MB

    // Проверяем тип файла, если он доступен
    const isValidType = !file.type ||
        file.type === 'text/csv' ||
        file.type === 'application/vnd.ms-excel' ||
        file.type === 'text/plain'

    return isValidExtension && isValidSize && isValidType
}

const handleImport = async () => {
    if (!selectedFile.value) {
        showError('Пожалуйста, выберите файл')
        return
    }

    if (!isValidCSVFile(selectedFile.value)) {
        showError('Некорректный файл. Пожалуйста, выберите CSV файл до 10MB')
        return
    }

    const reader = new FileReader()

    reader.onload = async (e) => {
        try {
            const content = e.target.result

            // Проверяем, что файл не пустой
            if (!content.trim()) {
                showError('Файл пуст')
                return
            }

            // Проверяем структуру CSV (должна быть хотя бы одна строка)
            const lines = content.split('\n').filter(line => line.trim())
            if (lines.length === 0) {
                showError('Файл не содержит данных')
                return
            }

            // Проверяем заголовки (опционально)
            const headers = lines[0].split(',')
            const expectedHeaders = ['jp', 'ru', 'on', 'kun', 'ex_jp', 'ex_ru', 'tags']
            const hasValidHeaders = headers.some(header =>
                expectedHeaders.includes(header.trim().toLowerCase())
            )

            if (!hasValidHeaders && lines.length > 0) {
                // Если нет ожидаемых заголовков, это может быть файл без заголовков
                console.warn('CSV файл не содержит ожидаемых заголовков, но будет обработан')
            }

            const result = await importStore.importCSV(content)

            if (result.success) {
                if (result.data && result.data.imported_count > 0) {
                    showSuccess(`Успешно импортировано ${result.data.imported_count} слов`)
                    emit('imported', result.data)

                    // Закрываем модальное окно через 2 секунды
                    setTimeout(() => {
                        close()
                    }, 2000)
                } else {
                    showInfo('Нет новых слов для импорта')
                }
            } else {
                showError(result.error || 'Ошибка при импорте файла')
            }
        } catch (error) {
            console.error('Ошибка при чтении файла:', error)
            showError('Ошибка при чтении файла')
        }
    }

    reader.onerror = () => {
        showError('Ошибка при чтении файла')
    }

    reader.readAsText(selectedFile.value, 'UTF-8')
}

const close = () => {
    selectedFile.value = null
    importStore.clearImportResult()
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

.drag-enter {
    border-color: #3B82F6;
    background-color: rgba(59, 130, 246, 0.1);
}
</style>