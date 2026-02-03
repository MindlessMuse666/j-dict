<template>
    <div class="min-h-[60vh] flex flex-col items-center justify-center">
        <div class="text-center">
            <h1 class="text-9xl font-bold text-stone-200 font-jp">404</h1>
            <h2 class="text-2xl font-semibold text-text-main mt-4 font-jp">Страница не найдена</h2>
            <p class="text-text-muted mt-2 mb-6">Запрошенная страница не существует</p>
            <RouterLink :to="{ name: 'Home' }"
                class="px-6 py-3 bg-primary hover:bg-primary/90 text-white font-medium rounded-lg transition-all duration-200 shadow-md hover:shadow-lg transform active:scale-95">
                Вернуться на главную
            </RouterLink>
        </div>
    </div>
</template>

<script setup>
// Страница 404
</script>