<template>
  <div class="min-h-screen bg-background text-text-primary p-6">
    <div class="max-w-4xl mx-auto">
      <h1 class="text-3xl font-bold mb-8 text-primary">Личный кабинет</h1>

      <div class="bg-surface rounded-2xl shadow-lg p-8 flex flex-col md:flex-row items-start gap-8">
        <!-- Avatar Section -->
        <div class="relative group">
          <div class="w-40 h-40 rounded-full overflow-hidden border-4 border-primary/20 shadow-md bg-background">
            <img :src="avatarUrl" alt="Avatar" class="w-full h-full object-cover" />
          </div>
          
          <!-- Hover Overlay -->
          <div class="absolute inset-0 bg-black/50 rounded-full flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity duration-200 cursor-pointer"
               @click.stop="showAvatarMenu = !showAvatarMenu">
            <svg class="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
            </svg>
          </div>

          <!-- Avatar Menu Dropdown -->
          <div v-if="showAvatarMenu" 
               class="absolute top-full left-0 mt-2 w-48 bg-surface rounded-xl shadow-xl border border-border z-10 overflow-hidden"
               v-click-outside="() => showAvatarMenu = false">
            <button @click="openUploadModal" class="w-full text-left px-4 py-2 hover:bg-background/50 transition-colors flex items-center text-text-primary">
              <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12" /></svg>
              Загрузить фото
            </button>
            <button @click="openPresetModal" class="w-full text-left px-4 py-2 hover:bg-background/50 transition-colors flex items-center text-text-primary">
              <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" /></svg>
              Выбрать из списка
            </button>
            <button @click="deleteAvatar" class="w-full text-left px-4 py-2 hover:bg-red-500/10 text-red-500 transition-colors flex items-center">
              <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" /></svg>
              Удалить
            </button>
          </div>
        </div>

        <!-- User Info -->
        <div class="flex-1">
          <div class="flex justify-between items-start">
            <div>
              <h2 class="text-2xl font-bold mb-2 text-text-primary">{{ user?.name || 'Пользователь' }}</h2>
              <p class="text-text-secondary mb-4">{{ user?.email }}</p>
            </div>
            <button @click="handleLogout" class="px-5 py-2 text-sm font-medium text-white bg-primary hover:bg-primary/90 rounded-xl transition-all duration-200 shadow-md hover:shadow-lg active:translate-y-0">
              Выйти
            </button>
          </div>
          
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 mt-8">
             <div class="bg-background p-4 rounded-xl border border-border">
                <div class="text-text-secondary text-sm mb-1">Всего слов</div>
                <div class="text-2xl font-bold text-primary">{{ wordsCount }}</div>
             </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Upload/Crop Modal -->
    <div v-if="showUploadModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm">
      <div class="bg-surface rounded-2xl shadow-xl max-w-2xl w-full p-6">
        <h3 class="text-xl font-bold mb-4 text-text-primary">Загрузка аватара</h3>
        
        <div v-if="!selectedFile" class="border-2 border-dashed border-border rounded-xl p-8 text-center cursor-pointer hover:border-primary transition-colors bg-background"
             @dragover.prevent @drop.prevent="handleDrop" @click="$refs.fileInput.click()">
          <input type="file" ref="fileInput" class="hidden" accept="image/jpeg,image/png" @change="handleFileSelect" />
          <svg class="w-12 h-12 text-text-secondary mx-auto mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
          </svg>
          <p class="text-text-secondary">Нажмите или перетащите изображение сюда</p>
        </div>

        <div v-else class="h-96 bg-black rounded-xl overflow-hidden mb-4 relative">
          <cropper
            ref="cropper"
            class="h-full"
            :src="selectedFileResult"
            :stencil-props="{ aspectRatio: 1 }"
          />
        </div>

        <div class="flex justify-end gap-3 mt-4">
          <button @click="closeUploadModal" class="px-4 py-2 text-text-secondary hover:bg-background rounded-lg transition-colors">Отмена</button>
          <button v-if="selectedFile" @click="saveAvatar" class="px-4 py-2 bg-primary hover:bg-primary/90 text-white rounded-lg transition-colors shadow-md">Сохранить</button>
        </div>
      </div>
    </div>

    <!-- Preset Modal -->
    <div v-if="showPresetModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm">
      <div class="bg-surface rounded-2xl shadow-xl max-w-2xl w-full p-6">
        <h3 class="text-xl font-bold mb-4 text-text-primary">Выберите аватара</h3>
        <div class="grid grid-cols-2 sm:grid-cols-4 gap-4">
           <button v-for="preset in presets" :key="preset" 
                   @click="selectPreset(preset)"
                   class="rounded-full overflow-hidden border-2 border-transparent hover:border-primary transition-all focus:outline-none focus:ring-2 focus:ring-primary aspect-square group relative">
              <img :src="`/assets/default_avatars/${preset}`" class="w-full h-full object-cover transform transition-transform group-hover:scale-110" />
           </button>
        </div>
        <div class="flex justify-end gap-3 mt-4">
          <button @click="showPresetModal = false" class="px-4 py-2 text-text-secondary hover:bg-background rounded-lg transition-colors">Отмена</button>
        </div>
      </div>
    </div>

    <!-- Toast -->
    <Transition
      enter-active-class="transform ease-out duration-300 transition"
      enter-from-class="translate-y-2 opacity-0 sm:translate-y-0 sm:translate-x-2"
      enter-to-class="translate-y-0 opacity-100 sm:translate-x-0"
      leave-active-class="transition ease-in duration-100"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div v-if="toast.show" class="fixed bottom-4 right-4 z-50 bg-surface border border-border rounded-xl shadow-lg p-4 flex items-center gap-3">
        <div class="w-2 h-2 rounded-full" :class="toast.type === 'error' ? 'bg-red-500' : 'bg-green-500'"></div>
        <p class="text-sm font-medium text-text-primary">{{ toast.message }}</p>
      </div>
    </Transition>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useWordsStore } from '@/stores/words'
import { useRouter } from 'vue-router'
import { Cropper } from 'vue-advanced-cropper'
import 'vue-advanced-cropper/dist/style.css'

const authStore = useAuthStore()
const wordsStore = useWordsStore()
const router = useRouter()

const user = computed(() => authStore.user)
const wordsCount = computed(() => wordsStore.totalCount)

const handleLogout = async () => {
  await authStore.logout()
  router.push({ name: 'Login' })
}

const avatarUrl = computed(() => {
  if (user.value?.avatar_url) {
    // If url starts with http/https, use it. 
    // If it starts with /uploads, it's backend served.
    // If it starts with /assets, it's static asset.
    return user.value.avatar_url
  }
  return '/assets/default_avatars/~catishe~cat~.jpg'
})

const showAvatarMenu = ref(false)
const showUploadModal = ref(false)
const showPresetModal = ref(false)
const selectedFile = ref(null)
const selectedFileResult = ref(null)
const cropper = ref(null)
const fileInput = ref(null)

const toast = ref({ show: false, message: '', type: 'success' })
const showToast = (message, type = 'success') => {
  toast.value = { show: true, message, type }
  setTimeout(() => toast.value.show = false, 3000)
}

const presets = [
  '~atom~bomb~cat~.jpg',
  '~catishe~cat~.jpg',
  '~cranchy~view~cat~.jpg',
  '~kayfishe~cat~.jpg',
  '~lunishe~cat~.jpg',
  '~mr~cat~.jpg',
  '~void~cat~.jpg',
  '~winter~cat~.jpg'
]

const vClickOutside = {
  mounted(el, binding) {
    el.clickOutsideEvent = function(event) {
      if (!(el === event.target || el.contains(event.target))) {
        binding.value(event)
      }
    }
    document.body.addEventListener('click', el.clickOutsideEvent)
  },
  unmounted(el) {
    document.body.removeEventListener('click', el.clickOutsideEvent)
  }
}

const openUploadModal = () => {
  showAvatarMenu.value = false
  showUploadModal.value = true
  selectedFile.value = null
  selectedFileResult.value = null
}

const openPresetModal = () => {
  showAvatarMenu.value = false
  showPresetModal.value = true
}

const closeUploadModal = () => {
  showUploadModal.value = false
  selectedFile.value = null
  selectedFileResult.value = null
}

const handleFileSelect = (e) => {
  const file = e.target.files[0]
  if (file) processFile(file)
  // Reset input so same file can be selected again
  e.target.value = ''
}

const handleDrop = (e) => {
  const file = e.dataTransfer.files[0]
  if (file) processFile(file)
}

const processFile = (file) => {
  if (!file.type.startsWith('image/')) {
    showToast('Пожалуйста, загрузите изображение', 'error')
    return
  }
  selectedFile.value = file
  const reader = new FileReader()
  reader.onload = (e) => {
    selectedFileResult.value = e.target.result
  }
  reader.readAsDataURL(file)
}

const saveAvatar = () => {
  const { canvas } = cropper.value.getResult()
  canvas.toBlob(async (blob) => {
    const res = await authStore.uploadAvatar(blob)
    if (res.success) {
      showToast('Аватар успешно обновлен')
      closeUploadModal()
    } else {
      showToast(res.error, 'error')
    }
  }, 'image/jpeg')
}

const selectPreset = async (filename) => {
  const url = `/assets/default_avatars/${filename}`
  const res = await authStore.updateAvatarUrl(url)
  if (res.success) {
    showToast('Аватар обновлен')
    showPresetModal.value = false
  } else {
    showToast(res.error, 'error')
  }
}

const deleteAvatar = async () => {
  const url = '/assets/default_avatars/~catishe~cat~.jpg'
  const res = await authStore.updateAvatarUrl(url)
  if (res.success) {
    showToast('Аватар сброшен')
    showAvatarMenu.value = false
  } else {
    showToast(res.error, 'error')
  }
}

onMounted(() => {
    authStore.fetchUser()
    wordsStore.fetchWords()
})
</script>
