<template>
    <header
        class="bg-surface shadow-sm border-b border-stone-200/50 backdrop-blur-sm bg-surface/90 sticky top-0 z-50 transition-colors duration-300">
        <div class="container mx-auto px-4 py-3">
            <div class="flex items-center justify-between">
                <!-- Логотип и название -->
                <RouterLink :to="{ name: 'Home' }" class="flex items-center space-x-3 group">
                    <img src="/logo.svg" alt="Logo" class="w-10 h-10 object-contain drop-shadow-sm" />
                    <div>
                        <h1 class="text-xl font-bold text-primary tracking-tight font-jp">
                            j~dict!^w^
                        </h1>
                        <RouterLink v-if="authStore.user" :to="{ name: 'Profile' }"
                            class="text-xs text-text-muted font-jp hover:text-primary transition-colors cursor-pointer block">
                            {{ authStore.user.name }}
                        </RouterLink>
                    </div>
                </RouterLink>

                <!-- Навигация -->
                <div class="flex items-center space-x-4">
                    <template v-if="authStore.isAuthenticated">
                        <div class="flex items-center gap-3">
                            <RouterLink :to="{ name: isProfilePage ? 'Home' : 'Profile' }"
                                class="px-5 py-2 text-sm font-medium text-white bg-primary hover:bg-primary/90 rounded-xl transition-all duration-200 shadow-md hover:shadow-lg active:translate-y-0 flex items-center gap-2">
                                <svg v-if="isProfilePage" class="w-4 h-4" fill="none" stroke="currentColor"
                                    viewBox="0 0 24 24">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                        d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253">
                                    </path>
                                </svg>
                                <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                        d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path>
                                </svg>
                                {{ isProfilePage ? 'К словарю' : 'Профиль' }}
                            </RouterLink>
                            <div
                                class="w-10 h-10 rounded-full overflow-hidden border border-stone-200 shadow-sm bg-background">
                                <img :src="avatarUrl" @error="handleImageError" class="w-full h-full object-cover">
                            </div>
                        </div>
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

/**
 * Компонент шапки приложения.
 * Отображает логотип, навигацию и информацию о пользователе.
 */
import { useAuthStore } from '@/stores/auth'
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const authStore = useAuthStore()
const route = useRoute()

const isProfilePage = computed(() => route.name === 'Profile')

const avatarUrl = computed(() => {
    if (authStore.user?.avatar_url) {
        return authStore.user.avatar_url
    }
    return '/assets/default_avatars/~catishe~cat~.jpg'
})

const handleImageError = (e) => {
    console.error('Header avatar load error:', e.target.src)
    e.target.src = '/assets/default_avatars/~catishe~cat~.jpg'
}
</script>