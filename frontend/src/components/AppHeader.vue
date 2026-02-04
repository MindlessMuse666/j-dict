<template>
    <header class="bg-surface shadow-sm border-b border-stone-200/50 backdrop-blur-sm bg-surface/90 sticky top-0 z-50 transition-colors duration-300">
        <div class="container mx-auto px-4 py-3">
            <div class="flex items-center justify-between">
                <!-- Логотип и название -->
                <RouterLink :to="{ name: 'Home' }" class="flex items-center space-x-3 group">
                    <img src="/logo.svg" alt="Logo" class="w-10 h-10 object-contain drop-shadow-sm" />
                    <div>
                        <h1 class="text-xl font-bold text-primary tracking-tight font-jp">
                            ~j~ru~dict!^^
                        </h1>
                        <p v-if="authStore.user" class="text-xs text-text-muted font-jp">
                            {{ authStore.user.name }}
                        </p>
                    </div>
                </RouterLink>

                <!-- Навигация -->
                <div class="flex items-center space-x-4">
                    <template v-if="authStore.isAuthenticated">
                        <button @click="handleLogout"
                            class="px-5 py-2 text-sm font-medium text-white bg-primary hover:bg-primary/90 rounded-xl transition-all duration-200 shadow-md hover:shadow-lg active:translate-y-0">
                            Выйти
                        </button>
                    </template>
                    <template v-else>
                        <RouterLink :to="{ name: 'Login' }"
                            class="px-4 py-2 text-sm font-medium text-text-muted hover:text-primary hover:bg-stone-100 rounded-lg transition-all duration-200">
                            Вход
                        </RouterLink>
                        <RouterLink :to="{ name: 'Register' }"
                            class="px-4 py-2 text-sm font-medium text-white bg-primary hover:bg-primary/90 shadow-md hover:shadow-lg rounded-lg transition-all duration-200 transform hover:-translate-y-0.5 active:scale-95">
                            Регистрация
                        </RouterLink>
                    </template>
                </div>
            </div>
        </div>
    </header>
</template>

<script setup>
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'
import { useToast } from '@/composables/useToast'

const authStore = useAuthStore()
const router = useRouter()
const { showSuccess } = useToast()

const handleLogout = async () => {
    await authStore.logout()
    showSuccess('Вы успешно вышли из системы')
    router.push({ name: 'Login' })
}
</script>