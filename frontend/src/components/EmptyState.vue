<template>
  <div class="text-center py-16 px-4">
    <div class="mb-6 relative inline-block">
      <!-- Zen Circle (Enso) inspired SVG -->
      <svg class="w-32 h-32 mx-auto text-stone-200" viewBox="0 0 100 100" fill="none" stroke="currentColor">
        <path d="M50 10 C 20 10 10 30 10 50 C 10 80 30 90 50 90 C 80 90 90 70 90 50" 
              stroke-width="3" stroke-linecap="round" stroke-dasharray="10 5" opacity="0.5" />
        <circle cx="50" cy="50" r="15" class="text-primary/20" fill="currentColor" stroke="none" />
      </svg>
      <div class="absolute inset-0 flex items-center justify-center">
        <slot name="icon">
          <svg class="w-12 h-12 text-stone-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
              d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
          </svg>
        </slot>
      </div>
    </div>
    
    <h3 class="text-xl font-medium text-text-main mb-2 font-jp tracking-wide">
      {{ title }}
    </h3>
    
    <p class="text-text-muted mb-8 max-w-md mx-auto">
      {{ description }}
    </p>
    
    <slot name="action"></slot>
  </div>
</template>

<script setup>
defineProps({
  title: {
    type: String,
    default: 'Здесь пока пусто'
  },
  description: {
    type: String,
    default: 'Список пуст. Самое время добавить что-то новое.'
  }
})
</script>