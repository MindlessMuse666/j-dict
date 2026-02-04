<template>
  <div class="min-h-screen bg-background text-text-primary p-6">
    <div class="max-w-4xl mx-auto">
      <h1 class="text-3xl font-bold mb-8 text-primary">Личный кабинет</h1>

      <div class="bg-surface rounded-2xl shadow-lg p-8 flex flex-col md:flex-row gap-8 items-start">
        <!-- Avatar Section -->
        <div class="relative group flex-shrink-0 mx-auto md:mx-0">
          <div class="w-40 h-40 rounded-full overflow-hidden border-4 border-primary/20 shadow-md bg-background relative z-0">
            <img :src="avatarUrl" @error="handleImageError" alt="Avatar" class="w-full h-full object-cover" />
          </div>
          
          <!-- Hover Overlay -->
          <div class="absolute inset-0 bg-black/50 rounded-full flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity duration-200 cursor-pointer z-10"
               @click.stop="showAvatarMenu = !showAvatarMenu">
            <svg class="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
            </svg>
          </div>

          <!-- Avatar Menu Dropdown -->
          <div v-if="showAvatarMenu" class="fixed inset-0 z-20" @click="showAvatarMenu = false"></div>
          
          <div v-if="showAvatarMenu" 
               class="absolute top-full left-1/2 -translate-x-1/2 mt-2 w-48 bg-surface rounded-xl shadow-xl border border-border z-30 overflow-hidden">
            <button @click="openUploadModal" class="w-full text-left px-4 py-2 hover:bg-background/50 transition-colors flex items-center text-text-primary">
              <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12" /></svg>
              Загрузить фото
            </button>
            <button @click="openPresetModal" class="w-full text-left px-4 py-2 hover:bg-background/50 transition-colors flex items-center text-text-primary">
              <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" /></svg>
              Выбрать готовую
            </button>
            <button @click="deleteAvatar" class="w-full text-left px-4 py-2 hover:bg-red-500/10 text-red-500 transition-colors flex items-center">
              <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" /></svg>
              Удалить
            </button>
          </div>
        </div>

        <!-- User Info -->
        <div class="flex-1 w-full">
          <div class="flex flex-col sm:flex-row justify-between items-start gap-4">
            <div class="space-y-1">
              <h2 class="text-3xl font-bold text-text-primary tracking-tight">{{ user?.name || 'Пользователь' }}</h2>
              <div class="flex items-center text-text-secondary">
                <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" /></svg>
                <span>{{ user?.email }}</span>
              </div>
              <div class="flex items-center text-text-secondary mt-1">
                <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" /></svg>
                <span>{{ userRole }}</span>
              </div>
            </div>
            
            <button @click="confirmLogout" class="w-full sm:w-auto px-5 py-2.5 text-sm font-medium text-white bg-primary hover:bg-primary/90 rounded-xl transition-all duration-200 shadow-md hover:shadow-lg active:translate-y-0 flex items-center justify-center gap-2">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" /></svg>
              Выйти
            </button>
          </div>
          
          <div class="h-px bg-border/50 my-4"></div>
          
          <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
             <div class="bg-background px-4 py-3 rounded-xl border border-border">
                <div class="flex items-center gap-2">
                  <svg class="w-4 h-4 text-primary" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" /></svg>
                  <span class="text-text-secondary text-sm font-medium">Изучено слов</span>
                  <span class="text-xl font-bold text-primary ml-1">{{ wordsCount }}</span>
                </div>
             </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Upload/Crop Modal -->
    <div v-if="showUploadModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm">
      <div class="bg-surface rounded-2xl shadow-xl max-w-2xl w-full p-6 relative">
        <button v-if="!selectedFileResult" @click="closeUploadModal" class="absolute top-4 right-4 text-text-secondary hover:text-text-primary transition-colors">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" /></svg>
        </button>

        <h3 class="text-xl font-bold mb-4 text-text-primary">Загрузка аватара</h3>
        
        <div v-if="!selectedFile" 
             class="border-2 border-dashed rounded-xl p-8 text-center cursor-pointer transition-all duration-200"
             :class="[
               isDragging ? 'border-primary bg-primary/5 scale-[1.02]' : 'border-border hover:border-primary bg-background'
             ]"
             @dragover.prevent 
             @dragenter.prevent="isDragging = true"
             @dragleave.prevent="isDragging = false"
             @drop.prevent="handleDrop" 
             @click="$refs.fileInput.click()">
          <input type="file" ref="fileInput" class="hidden" accept="image/jpeg,image/png" @change="handleFileSelect" />
          <svg class="w-12 h-12 mx-auto mb-2 transition-colors duration-200" 
               :class="isDragging ? 'text-primary' : 'text-text-secondary'"
               fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
          </svg>
          <p class="transition-colors duration-200" :class="isDragging ? 'text-primary font-medium' : 'text-text-secondary'">
            {{ isDragging ? 'Отпустите файл для загрузки' : 'Нажмите или перетащите изображение сюда' }}
          </p>
        </div>

        <div v-else class="flex flex-col items-center">
          <div 
            ref="cropperContainer"
            class="relative w-full h-80 bg-stone-900 rounded-xl overflow-hidden cursor-move touch-none select-none"
            @mousedown="startDrag"
            @touchstart.prevent="startDrag"
            @mousemove="onDrag"
            @touchmove.prevent="onDrag"
            @mouseup="stopDrag"
            @touchend="stopDrag"
            @mouseleave="stopDrag"
            @wheel.prevent="onWheel"
          >
            <img 
              ref="cropperImage"
              :src="selectedFileResult" 
              class="absolute max-w-none origin-center pointer-events-none will-change-transform"
              :style="{
                transform: `translate(${imageState.x}px, ${imageState.y}px) scale(${imageState.scale})`,
                width: `${imageState.width}px`,
                height: `${imageState.height}px`,
                left: '50%',
                top: '50%',
                marginLeft: `-${imageState.width / 2}px`,
                marginTop: `-${imageState.height / 2}px`
              }"
              @load="onImageLoad"
            />
            
            <!-- Overlay with circular cutout -->
            <div class="absolute inset-0 pointer-events-none z-10">
               <div class="absolute inset-0 bg-black/60"></div>
               <!-- The clear circle/square area -->
               <div 
                 class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 border-2 border-white/80 rounded-full shadow-[0_0_0_9999px_rgba(0,0,0,0.6)] box-content"
                 :style="{ width: `${CROP_SIZE}px`, height: `${CROP_SIZE}px` }"
               ></div>
            </div>
          </div>

          <!-- Controls -->
          <div class="w-full mt-6 px-2 flex items-center gap-4">
            <svg class="w-5 h-5 text-text-secondary flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 12H4"/></svg>
            <input 
              type="range" 
              v-model.number="imageState.scale" 
              :min="minScale" 
              :max="maxScale" 
              step="0.01"
              class="flex-1 h-1.5 bg-stone-200 rounded-lg appearance-none cursor-pointer accent-primary hover:accent-primary/80"
            >
            <svg class="w-5 h-5 text-text-secondary flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/></svg>
          </div>
        </div>
        
        <div class="flex justify-end gap-3 mt-6" v-if="selectedFileResult">
          <button @click="closeUploadModal" class="px-6 py-2 bg-transparent hover:bg-stone-100 text-text-secondary font-medium rounded-xl transition-all">Отмена</button>
          <button @click="saveAvatar" class="px-6 py-2 bg-primary hover:bg-primary/90 text-white font-medium rounded-xl transition-all shadow-md hover:shadow-lg active:translate-y-0">Сохранить</button>
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

    <!-- Logout Confirmation Modal -->
    <div v-if="showLogoutModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm">
      <div class="bg-surface rounded-2xl shadow-xl max-w-md w-full p-6">
        <h3 class="text-xl font-bold mb-4 text-text-primary">Вы точно хотите выйти?</h3>
        <div class="flex justify-end gap-3 mt-6">
          <button @click="showLogoutModal = false" class="px-6 py-2 bg-transparent hover:bg-stone-100 text-text-secondary font-medium rounded-xl transition-all">Нет</button>
          <button @click="handleLogout" class="px-6 py-2 bg-primary hover:bg-primary/90 text-white font-medium rounded-xl transition-all shadow-md hover:shadow-lg active:translate-y-0">Да</button>
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
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useWordsStore } from '@/stores/words'
import { useRouter } from 'vue-router'

const authStore = useAuthStore()
const wordsStore = useWordsStore()
const router = useRouter()

const user = computed(() => authStore.user)
const wordsCount = computed(() => wordsStore.totalCount)
const userRole = computed(() => {
  const role = user.value?.role || 'user'
  return role === 'admin' ? 'Администратор' : 'Пользователь'
})

const confirmLogout = () => {
  showLogoutModal.value = true
}

const handleLogout = async () => {
  showLogoutModal.value = false
  await authStore.logout()
  router.push({ name: 'Login' })
  showToast('Вы успешно вышли из системы')
}

const avatarUrl = computed(() => {
  if (user.value?.avatar_url) {
    // If url starts with http/https, use it. 
    // If it starts with /uploads, it's backend served.
    // If it starts with /assets, it's static asset.
    // Add timestamp to bust cache
    return `${user.value.avatar_url}?t=${Date.now()}`
  }
  return '/assets/default_avatars/~catishe~cat~.jpg'
})

const handleImageError = (e) => {
  e.target.src = '/assets/default_avatars/~catishe~cat~.jpg'
}

const showAvatarMenu = ref(false)
const showUploadModal = ref(false)
const showPresetModal = ref(false)
const showLogoutModal = ref(false)
const selectedFile = ref(null)
const selectedFileResult = ref(null)
const cropperContainer = ref(null)
const cropperImage = ref(null)
const fileInput = ref(null)
const isDragging = ref(false)

// Custom Cropper State
const CROP_SIZE = 280
const imageState = ref({
  x: 0,
  y: 0,
  scale: 1,
  width: 0,
  height: 0,
  naturalWidth: 0,
  naturalHeight: 0
})
const minScale = ref(0.1)
const maxScale = ref(3)
const dragStart = ref({ x: 0, y: 0, imageX: 0, imageY: 0 })
const isImageDragging = ref(false)

const onImageLoad = (e) => {
  const img = e.target
  imageState.value.naturalWidth = img.naturalWidth
  imageState.value.naturalHeight = img.naturalHeight
  
  // Calculate initial dimensions to fit container
  // We want the image to be at least large enough to cover the crop area
  const containerW = 320 // approx width of container in modal
  const containerH = 320 
  
  // Reset
  imageState.value.x = 0
  imageState.value.y = 0
  
  // Set display size (keep aspect ratio)
  // We set base width to match natural width initially, but scale it down
  imageState.value.width = img.naturalWidth
  imageState.value.height = img.naturalHeight
  
  // Calculate min scale to cover CROP_SIZE
  const scaleW = CROP_SIZE / img.naturalWidth
  const scaleH = CROP_SIZE / img.naturalHeight
  minScale.value = Math.max(scaleW, scaleH)
  maxScale.value = minScale.value * 5
  
  // Set initial scale to fit nicely (slightly larger than min)
  imageState.value.scale = minScale.value
}

const startDrag = (e) => {
  isImageDragging.value = true
  const clientX = e.type.includes('mouse') ? e.clientX : e.touches[0].clientX
  const clientY = e.type.includes('mouse') ? e.clientY : e.touches[0].clientY
  
  dragStart.value = {
    x: clientX,
    y: clientY,
    imageX: imageState.value.x,
    imageY: imageState.value.y
  }
}

const onDrag = (e) => {
  if (!isImageDragging.value) return
  
  const clientX = e.type.includes('mouse') ? e.clientX : e.touches[0].clientX
  const clientY = e.type.includes('mouse') ? e.clientY : e.touches[0].clientY
  
  const deltaX = clientX - dragStart.value.x
  const deltaY = clientY - dragStart.value.y
  
  imageState.value.x = dragStart.value.imageX + deltaX
  imageState.value.y = dragStart.value.imageY + deltaY
}

const stopDrag = () => {
  isImageDragging.value = false
}

const onWheel = (e) => {
  const delta = e.deltaY > 0 ? -0.1 : 0.1
  const newScale = Math.max(minScale.value, Math.min(maxScale.value, imageState.value.scale + delta))
  imageState.value.scale = newScale
}

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
  if (selectedFileResult.value) {
    URL.revokeObjectURL(selectedFileResult.value)
  }
  selectedFile.value = null
  selectedFileResult.value = null
  isDragging.value = false
}

const handleFileSelect = (e) => {
  const file = e.target.files[0]
  if (file) processFile(file)
  // Reset input so same file can be selected again
  e.target.value = ''
}

const handleDrop = (e) => {
  isDragging.value = false
  const file = e.dataTransfer.files[0]
  if (file) processFile(file)
}

const processFile = (file) => {
  if (!file.type.startsWith('image/')) {
    showToast('Пожалуйста, загрузите изображение', 'error')
    return
  }
  
  // Cleanup previous object URL if exists
  if (selectedFileResult.value) {
    URL.revokeObjectURL(selectedFileResult.value)
  }

  selectedFile.value = file
  // Use createObjectURL instead of FileReader to avoid large strings and potential stack overflows
  selectedFileResult.value = URL.createObjectURL(file)
}

const saveAvatar = () => {
  if (!selectedFile.value) return

  const img = new Image()
  img.src = selectedFileResult.value
  
  img.onload = () => {
    const canvas = document.createElement('canvas')
    const ctx = canvas.getContext('2d')
    
    // Set canvas to crop size (high res)
    const OUTPUT_SIZE = 500
    canvas.width = OUTPUT_SIZE
    canvas.height = OUTPUT_SIZE
    
    // Calculate mapping
    // We need to map the "visible crop area" relative to the image, to the canvas
    // 1. Center of container is (containerW/2, containerH/2)
    // 2. Image is at (x, y) relative to center + offset (since we centered it with CSS left:50% top:50% margin:-w/2)
    // Let's simplify:
    // The image center is at (ContainerCenter + x, ContainerCenter + y)
    // The Crop Area center is at ContainerCenter.
    // So Crop Area Center relative to Image Center is (-x, -y) (scaled pixels)
    // In unscaled pixels: (-x / scale, -y / scale)
    // Crop Area TopLeft relative to Image Center (unscaled): (-x/scale - CROP_SIZE/2/scale, -y/scale - CROP_SIZE/2/scale)
    // Image Center is at (naturalWidth/2, naturalHeight/2)
    
    const scale = imageState.value.scale
    const centerOffsetX = -imageState.value.x / scale
    const centerOffsetY = -imageState.value.y / scale
    
    const cropSizeUnscaled = CROP_SIZE / scale
    
    const sourceX = (imageState.value.naturalWidth / 2) + centerOffsetX - (cropSizeUnscaled / 2)
    const sourceY = (imageState.value.naturalHeight / 2) + centerOffsetY - (cropSizeUnscaled / 2)
    
    ctx.drawImage(
      img, 
      sourceX, sourceY, cropSizeUnscaled, cropSizeUnscaled, 
      0, 0, OUTPUT_SIZE, OUTPUT_SIZE
    )
    
    canvas.toBlob(async (blob) => {
      const res = await authStore.uploadAvatar(blob)
      if (res.success) {
        showToast('Аватар успешно обновлен')
        closeUploadModal()
      } else {
        const errorMap = {
          'Only JPG and PNG are allowed': 'Разрешены только форматы JPG и PNG',
          'File too large': 'Файл слишком большой',
          'Invalid file type': 'Неверный тип файла'
        }
        showToast(errorMap[res.error] || res.error, 'error')
      }
    }, 'image/jpeg', 0.9)
  }
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

onUnmounted(() => {
  if (selectedFileResult.value) {
    URL.revokeObjectURL(selectedFileResult.value)
  }
})
</script>
