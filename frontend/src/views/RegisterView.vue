<template>
    <div class="max-w-md mx-auto pt-8">
        <div class="bg-surface rounded-xl shadow-sm border border-stone-100 p-8">
            <div class="text-center mb-8">
                <h1 class="text-2xl font-bold text-text-main mb-2 font-jp">Регистрация</h1>
                <p class="text-text-muted">Создайте новый аккаунт</p>
            </div>

            <form @submit.prevent="handleRegister" class="space-y-6">
                <!-- Имя -->
                <div>
                    <label for="name" class="block text-sm font-medium text-text-main mb-2">
                        Имя
                    </label>
                    <input id="name" v-model="form.name" type="text" required
                        class="w-full px-4 py-3 border border-stone-200 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent transition-colors placeholder-stone-400"
                        placeholder="Иван Иванов" />
                </div>

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
                    <input id="password" v-model="form.password" type="password" required minlength="6"
                        class="w-full px-4 py-3 border border-stone-200 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent transition-colors placeholder-stone-400"
                        placeholder="Введите ваш пароль" />
                    <p class="mt-1 text-sm text-text-muted">Минимум 6 символов</p>
                </div>

                <!-- Подтверждение пароля -->
                <div>
                    <label for="confirmPassword" class="block text-sm font-medium text-text-main mb-2">
                        Подтверждение пароля
                    </label>
                    <input id="confirmPassword" v-model="form.confirmPassword" type="password" required
                        class="w-full px-4 py-3 border border-stone-200 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent transition-colors placeholder-stone-400"
                        placeholder="" />
                    <p v-if="passwordMismatch" class="mt-1 text-sm text-red-600">
                        Пароли не совпадают
                    </p>
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
                <button type="submit" :disabled="loading || passwordMismatch"
                    class="w-full py-3 px-4 bg-primary hover:bg-primary/90 text-white font-medium rounded-lg transition-all duration-200 shadow-md hover:shadow-lg disabled:opacity-50 disabled:cursor-not-allowed transform active:scale-95">
                    <span v-if="loading" class="flex items-center justify-center">
                        <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                        </svg>
                        Регистрация...
                    </span>
                    <span v-else>Зарегистрироваться</span>
                </button>

                <!-- Ссылка на вход -->
                <div class="text-center pt-4 border-t border-stone-100">
                    <p class="text-text-muted">
                        Уже есть аккаунт?
                        <RouterLink :to="{ name: 'Login' }" class="text-primary hover:text-primary/80 font-medium transition-colors">
                            Войти
                        </RouterLink>
                    </p>
                </div>
            </form>
        </div>
    </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useToast } from '@/composables/useToast'

const router = useRouter()
const authStore = useAuthStore()
const { showSuccess } = useToast()

const loading = ref(false)
const error = ref('')

const form = reactive({
    name: '',
    email: '',
    password: '',
    confirmPassword: ''
})

const passwordMismatch = computed(() => {
    return form.password !== form.confirmPassword && form.confirmPassword.length > 0
})

const handleRegister = async () => {
    if (passwordMismatch.value) return

    loading.value = true
    error.value = ''

    try {
        const result = await authStore.register(form.email, form.password, form.name)

        if (result.success) {
            showSuccess('Регистрация прошла успешно')
            router.push({ name: 'Home' })
        } else {
            error.value = result.error
        }
    } catch (err) {
        error.value = 'Произошла ошибка при регистрации'
    } finally {
        loading.value = false
    }
}
</script>