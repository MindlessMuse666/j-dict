<template>
    <Transition name="modal">
        <div v-if="isOpen" class="fixed inset-0 z-50 flex items-center justify-center">
            <!-- Фон с затемнением -->
            <div class="fixed inset-0 bg-black bg-opacity-50 transition-opacity" @click="$emit('close')"></div>

            <!-- Модальное окно -->
            <div class="relative z-50 bg-white rounded-xl shadow-xl w-full max-w-md mx-4 transform transition-all p-6">
                <div class="text-center">
                    <!-- Иконка предупреждения -->
                    <div class="mx-auto flex items-center justify-center h-12 w-12 rounded-full bg-red-100 mb-4">
                        <svg class="h-6 w-6 text-red-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                        </svg>
                    </div>

                    <!-- Заголовок и текст -->
                    <h3 class="text-lg leading-6 font-medium text-gray-900 mb-2">
                        {{ title }}
                    </h3>
                    <div class="mt-2">
                        <p class="text-sm text-gray-500">
                            {{ message }}
                        </p>
                    </div>

                    <!-- Кнопки -->
                    <div class="mt-6 flex justify-center space-x-3">
                        <button type="button" class="px-4 py-2 bg-white text-gray-700 border border-gray-300 rounded-lg hover:bg-gray-50 font-medium focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary transition-colors"
                            @click="$emit('close')">
                            Отмена
                        </button>
                        <button type="button" class="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 font-medium focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 transition-colors shadow-sm"
                            @click="$emit('confirm')">
                            Удалить
                        </button>
                    </div>
                </div>
            </div>
        </div>
    </Transition>
</template>

<script setup>
defineProps({
    isOpen: {
        type: Boolean,
        default: false
    },
    title: {
        type: String,
        default: 'Подтверждение действия'
    },
    message: {
        type: String,
        default: 'Вы уверены, что хотите выполнить это действие? Это действие необратимо.'
    }
})

defineEmits(['close', 'confirm'])
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