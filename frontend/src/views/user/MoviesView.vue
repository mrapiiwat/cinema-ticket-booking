<template>
  <div class="min-h-screen bg-stone-50">
    <main class="max-w-5xl mx-auto px-6 py-10">
      <div class="flex items-center justify-between mb-8 gap-4">
        <h1 class="text-xl font-semibold text-stone-800">Now Showing</h1>
        <div class="relative">
          <svg
            class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-stone-300"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
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

      <div v-if="filtered.length" class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-5">
        <div
          v-for="movie in filtered"
          :key="movie.id"
          @click="$router.push({ name: 'booking', params: { movieId: movie.id } })"
          class="group cursor-pointer"
        >
          <div class="aspect-2/3 rounded-xl overflow-hidden bg-stone-100 relative">
            <img
              v-if="movie.poster"
              :src="movie.poster"
              :alt="movie.title"
              class="w-full h-full object-cover transition-transform duration-300 group-hover:scale-105"
            />
            <div v-else class="w-full h-full flex items-center justify-center">
              <svg
                class="w-8 h-8 text-stone-300"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <rect x="2" y="7" width="20" height="13" rx="2" stroke-width="1.5" />
                <path d="M7 7V5a2 2 0 0 1 2-2h6a2 2 0 0 1 2 2v2" stroke-width="1.5" />
              </svg>
            </div>
            <div
              class="absolute inset-0 bg-stone-900/0 group-hover:bg-stone-900/10 transition-colors duration-200 rounded-xl"
            ></div>
          </div>

          <div class="mt-2.5 px-0.5">
            <p class="text-sm font-medium text-stone-800 truncate">{{ movie.title }}</p>
            <div class="flex items-center gap-1.5 mt-1">
              <svg
                class="w-3 h-3 text-stone-300 shrink-0"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <circle cx="12" cy="12" r="10" stroke-width="1.5" />
                <path d="M12 6v6l4 2" stroke-width="1.5" stroke-linecap="round" />
              </svg>
              <span class="text-xs text-stone-400">{{ movie.time }}</span>
            </div>
          </div>
        </div>
      </div>

      <div v-else class="flex flex-col items-center justify-center py-24 text-center">
        <svg
          class="w-10 h-10 text-stone-200 mb-3"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <circle cx="11" cy="11" r="8" stroke-width="1.5" />
          <path d="M21 21l-4.35-4.35" stroke-width="1.5" stroke-linecap="round" />
        </svg>
        <p class="text-sm text-stone-400">
          No movies found for "<span class="text-stone-600">{{ search }}</span
          >"
        </p>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

// Mock movies — replace with API call
const movies = ref([
  { id: 1, title: 'Dune: Part Two', time: '14:30 · 2h 46m', poster: '' },
  { id: 2, title: 'Oppenheimer', time: '15:00 · 3h 0m', poster: '' },
  { id: 3, title: 'Poor Things', time: '16:00 · 2h 21m', poster: '' },
  { id: 4, title: 'The Zone of Interest', time: '17:30 · 1h 45m', poster: '' },
  { id: 5, title: 'Killers of the Flower Moon', time: '18:00 · 3h 26m', poster: '' },
  { id: 6, title: 'Past Lives', time: '19:00 · 1h 46m', poster: '' },
  { id: 7, title: 'Saltburn', time: '20:15 · 2h 11m', poster: '' },
  { id: 8, title: 'All of Us Strangers', time: '21:00 · 1h 58m', poster: '' },
])

const search = ref('')

const filtered = computed(() =>
  movies.value.filter((m) => m.title.toLowerCase().includes(search.value.toLowerCase())),
)
</script>
