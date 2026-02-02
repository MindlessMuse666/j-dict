<template>
    <div class="fixed bottom-4 right-4 z-50 flex flex-col space-y-2 max-w-md">
        <TransitionGroup name="toast">
            <div v-for="toast in toasts" :key="toast.id" :class="[
                'px-4 py-3 rounded-lg shadow-lg',
                toastClasses[toast.type]
            ]">
                <div class="flex items-center justify-between">
                    <div class="flex items-center">
                        <!-- Иконки для разных типов тостов -->
                        <svg v-if="toast.type === 'success'" class="w-5 h-5 mr-2 text-green-600" fill="currentColor"
                            viewBox="0 0 20 20">
                            <path fill-rule="evenodd"
                                d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                                clip-rule="evenodd" />
                        </svg>
                        <svg v-if="toast.type === 'error'" class="w-5 h-5 mr-2 text-red-600" fill="currentColor"
                            viewBox="0 0 20 20">
                            <path fill-rule="evenodd"
                                d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
                                clip-rule="evenodd" />
                        </svg>
                        <svg v-if="toast.type === 'warning'" class="w-5 h-5 mr-2 text-yellow-600" fill="currentColor"
                            viewBox="0 0 20 20">
                            <path fill-rule="evenodd"
                                d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z"
                                clip-rule="evenodd" />
                        </svg>
                        <svg v-if="toast.type === 'info'" class="w-5 h-5 mr-2 text-blue-600" fill="currentColor"
                            viewBox="0 0 20 20">
                            <path fill-rule="evenodd"
                                d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z"
                                clip-rule="evenodd" />
                        </svg>
                        <span class="text-sm">{{ toast.message }}</span>
                    </div>
                    <button @click="removeToast(toast.id)"
                        class="ml-4 text-lg hover:opacity-70 transition-opacity text-gray-500">
                        &times;
                    </button>
                </div>
            </div>
        </TransitionGroup>
    </div>
</template>

<script setup>
import { useToast } from '@/composables/useToast'

const { toasts, removeToast } = useToast()

const toastClasses = {
    success: 'bg-green-50 text-green-800 border border-green-200',
    error: 'bg-red-50 text-red-800 border border-red-200',
    warning: 'bg-yellow-50 text-yellow-800 border border-yellow-200',
    info: 'bg-blue-50 text-blue-800 border border-blue-200'
}
</script>

<style scoped>
.toast-enter-active,
.toast-leave-active {
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.toast-enter-from,
.toast-leave-to {
    opacity: 0;
    transform: translateX(100%);
}

.toast-leave-active {
    position: absolute;
}

.toast-move {
    transition: transform 0.3s ease;
}
</style>