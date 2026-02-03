<template>
    <div :class="[
        'bg-white rounded-lg border transition-all duration-200',
        isCompact ? 'p-4' : 'p-6'
    ]">
        <div v-if="isCompact" class="flex items-center justify-between">
            <!-- Японский текст (скелетон) -->
            <div class="flex-1 mr-4">
                <div class="flex items-center gap-2 mb-1">
                    <div class="h-6 w-16 bg-stone-200 rounded animate-pulse"></div>
                    <div class="h-6 w-12 bg-stone-200 rounded animate-pulse"></div>
                </div>
                <!-- Чтения (скелетон) -->
                <div class="flex gap-2 mt-1">
                    <div class="h-4 w-8 bg-stone-100 rounded animate-pulse"></div>
                    <div class="h-4 w-10 bg-stone-100 rounded animate-pulse"></div>
                </div>
            </div>

            <!-- Русский перевод (скелетон) -->
            <div class="flex-1 text-right">
                <div class="flex justify-end gap-2 mb-1">
                    <div class="h-6 w-20 bg-stone-200 rounded animate-pulse"></div>
                </div>
                <!-- Теги (скелетон) -->
                <div class="flex justify-end gap-1 mt-1">
                    <div class="h-4 w-10 bg-stone-100 rounded animate-pulse"></div>
                    <div class="h-4 w-8 bg-stone-100 rounded animate-pulse"></div>
                </div>
            </div>

            <!-- Кнопки (скелетон) -->
            <div class="ml-4 flex items-center space-x-2">
                <div class="h-8 w-8 bg-stone-100 rounded-full animate-pulse"></div>
                <div class="h-8 w-8 bg-stone-100 rounded-full animate-pulse"></div>
            </div>
        </div>

        <!-- Полный вид (скелетон) -->
        <div v-else>
            <div class="mb-4">
                <div class="h-8 w-24 bg-stone-200 rounded mb-3 animate-pulse"></div>
                <div class="flex gap-2 mb-3">
                    <div class="h-6 w-16 bg-stone-100 rounded animate-pulse"></div>
                    <div class="h-6 w-14 bg-stone-100 rounded animate-pulse"></div>
                </div>
            </div>

            <div class="mb-6">
                <div class="h-4 w-20 bg-stone-200 rounded mb-2 animate-pulse"></div>
                <div class="space-y-1">
                    <div class="h-4 w-full bg-stone-100 rounded animate-pulse"></div>
                    <div class="h-4 w-2/3 bg-stone-100 rounded animate-pulse"></div>
                </div>
            </div>

            <div class="flex items-center justify-between pt-4 border-t border-stone-100">
                <div class="flex gap-2">
                    <div class="h-6 w-12 bg-stone-100 rounded-full animate-pulse"></div>
                    <div class="h-6 w-10 bg-stone-100 rounded-full animate-pulse"></div>
                </div>
                <div class="h-4 w-24 bg-stone-100 rounded animate-pulse"></div>
            </div>

            <div class="mt-4 flex justify-end space-x-3">
                <div class="h-8 w-24 bg-stone-100 rounded animate-pulse"></div>
                <div class="h-8 w-20 bg-stone-100 rounded animate-pulse"></div>
            </div>
        </div>
    </div>
</template>

<script setup>
defineProps({
    isCompact: {
        type: Boolean,
        default: true
    }
})
</script>