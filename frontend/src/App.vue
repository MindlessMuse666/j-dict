<template>
    <div class="min-h-screen bg-background bg-washi-pattern font-sans text-text-main">
        <AppHeader />
        <main class="container mx-auto px-4 py-8 relative z-10">
            <RouterView />
        </main>
        <ToastContainer />
    </div>
</template>

<script setup>
/**
 * Корневой компонент App.
 * Отвечает за структуру макета и глобальные провайдеры.
 */
import AppHeader from '@/components/AppHeader.vue'
import ToastContainer from '@/components/ToastContainer.vue'
import { provideToast } from '@/composables/useToast'

// Предоставляем тосты на уровне корневого компонента
provideToast()
</script>