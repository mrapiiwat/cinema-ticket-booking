<template>
  <div class="min-h-screen bg-stone-50">
    <main class="max-w-5xl mx-auto px-6 py-10">
      <div class="flex items-center justify-between mb-8 gap-4">
        <h1 class="text-xl font-semibold text-stone-800">Now Showing</h1>
        <div class="relative">
          <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-stone-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <circle cx="11" cy="11" r="8" stroke-width="1.5" />
            <path d="M21 21l-4.35-4.35" stroke-width="1.5" stroke-linecap="round" />
          </svg>
          <input
            v-model="search"
            type="text"
            placeholder="Search movies..."
            class="pl-9 pr-4 py-2 text-sm bg-white border border-stone-200 rounded-lg text-stone-700 placeholder-stone-300 focus:outline-none focus:border-stone-400 w-48"
          />
        </div>
      </div>

      <div v-if="loading" class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-5">
        <div v-for="item in 8" :key="item" class="aspect-2/3 rounded-xl bg-stone-100 animate-pulse"></div>
      </div>

      <div v-else-if="filtered.length" class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-5">
        <button
          v-for="movie in filtered"
          :key="movie.id"
          @click="$router.push({ name: 'booking', params: { movieId: movie.id } })"
          class="group cursor-pointer text-left"
        >
          <div class="aspect-2/3 rounded-xl overflow-hidden bg-stone-100 relative">
            <img
              v-if="movie.poster_url"
              :src="movie.poster_url"
              :alt="movie.title"
              class="w-full h-full object-cover transition-transform duration-300 group-hover:scale-105"
            />
            <div v-else class="w-full h-full flex items-center justify-center bg-gradient-to-br from-stone-100 via-white to-sky-50">
              <svg class="w-9 h-9 text-stone-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <rect x="2" y="7" width="20" height="13" rx="2" stroke-width="1.5" />
                <path d="M7 7V5a2 2 0 0 1 2-2h6a2 2 0 0 1 2 2v2" stroke-width="1.5" />
              </svg>
            </div>
            <div class="absolute inset-0 bg-stone-900/0 group-hover:bg-stone-900/10 transition-colors duration-200 rounded-xl"></div>
          </div>

          <div class="mt-2.5 px-0.5">
            <p class="text-sm font-medium text-stone-800 truncate">{{ movie.title }}</p>
            <div class="flex items-center gap-1.5 mt-1">
              <svg class="w-3 h-3 text-sky-400 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <circle cx="12" cy="12" r="10" stroke-width="1.5" />
                <path d="M12 6v6l4 2" stroke-width="1.5" stroke-linecap="round" />
              </svg>
              <span class="text-xs text-stone-400">{{ formatDuration(movie.duration_mins) }}</span>
            </div>
          </div>
        </button>
      </div>

      <div v-else class="flex flex-col items-center justify-center py-24 text-center">
        <svg class="w-10 h-10 text-stone-200 mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <circle cx="11" cy="11" r="8" stroke-width="1.5" />
          <path d="M21 21l-4.35-4.35" stroke-width="1.5" stroke-linecap="round" />
        </svg>
        <p class="text-sm text-stone-400">{{ error || `No movies found for "${search}"` }}</p>
      </div>
    </main>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { apiFetch } from '../../api/client'

const movies = ref([])
const search = ref('')
const loading = ref(true)
const error = ref('')

const filtered = computed(() =>
  movies.value.filter((movie) => movie.title.toLowerCase().includes(search.value.toLowerCase())),
)

const formatDuration = (minutes) => {
  const hours = Math.floor(minutes / 60)
  const mins = minutes % 60
  return hours ? `${hours}h ${mins}m` : `${mins}m`
}

onMounted(async () => {
  try {
    movies.value = await apiFetch('/movies')
  } catch (err) {
    error.value = err.message
  } finally {
    loading.value = false
  }
})
</script>
