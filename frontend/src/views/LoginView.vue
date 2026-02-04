<template>
    <div class="max-w-md mx-auto pt-8">
        <div class="bg-surface rounded-xl shadow-sm border border-stone-100 p-8">
            <div class="text-center mb-8">
                <h1 class="text-2xl font-bold text-text-main mb-2 font-jp">Вход в систему</h1>
                <p class="text-text-muted">Введите свои учетные данные</p>
            </div>

            <form @submit.prevent="handleLogin" class="space-y-6">
                <!-- Email -->
                <div>
                    <label for="email" class="block text-sm font-medium text-text-main mb-2">
                        Email
                    </label>
                    <input id="email" v-model="form.email" type="email" required
                        class="w-full px-4 py-3 border border-stone-200 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent transition-colors placeholder-stone-400"
                        placeholder="Введите ваш email" />
                </div>

                <!-- Пароль -->
                <div>
                    <label for="password" class="block text-sm font-medium text-text-main mb-2">
                        Пароль
                    </label>
                    <input id="password" v-model="form.password" type="password" required
                        class="w-full px-4 py-3 border border-stone-200 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent transition-colors placeholder-stone-400"
                        placeholder="Введите ваш пароль" />
                </div>

                <!-- Ошибка -->
                <div v-if="error" class="p-4 bg-red-50 border border-red-100 rounded-lg">
                    <div class="flex items-center">
                        <svg class="w-5 h-5 text-red-500 mr-3" fill="currentColor" viewBox="0 0 20 20">
                            <path fill-rule="evenodd"
                                d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
                                clip-rule="evenodd" />
                        </svg>
                        <span class="text-red-800">{{ error }}</span>
                    </div>
                </div>

                <!-- Кнопка отправки -->
                <button type="submit" :disabled="loading"
                    class="w-full py-3 px-4 bg-primary hover:bg-primary/90 text-white font-medium rounded-lg transition-all duration-200 shadow-md hover:shadow-lg disabled:opacity-50 disabled:cursor-not-allowed transform active:scale-95">
                    <span v-if="loading" class="flex items-center justify-center">
                        <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                        </svg>
                        Вход...
                    </span>
                    <span v-else>Войти</span>
                </button>

                <!-- Ссылка на регистрацию -->
                <div class="text-center pt-4 border-t border-stone-100">
                    <p class="text-text-muted">
                        Нет аккаунта?
                        <RouterLink :to="{ name: 'Register' }" class="text-primary hover:text-primary/80 font-medium transition-colors">
                            Зарегистрироваться
                        </RouterLink>
                    </p>
                </div>
            </form>
        </div>
    </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useToast } from '@/composables/useToast'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const { showSuccess } = useToast()

const loading = ref(false)
const error = ref('')

const form = reactive({
    email: '',
    password: ''
})

const handleLogin = async () => {
    loading.value = true
    error.value = ''

    try {
        const result = await authStore.login(form.email, form.password)

        if (result.success) {
            showSuccess('Вы успешно вошли в систему')

            // Перенаправляем на страницу, с которой пришел пользователь, или на главную
            const redirect = route.query.redirect || '/'
            router.push(redirect)
        } else {
            error.value = result.error
        }
    } catch (err) {
        error.value = 'Произошла ошибка при входе'
    } finally {
        loading.value = false
    }
}
</script>