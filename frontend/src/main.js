import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { provideToast } from '@/composables/useToast'
import App from './App.vue'
import router from './router'
import './styles/tailwind.css'
import './styles/responsive.css'

const app = createApp(App)
const toast = provideToast(app)

app.use(createPinia())
app.use(router)

app.mount('#app')

if (import.meta.env.DEV) {
    window.$toast = toast
}