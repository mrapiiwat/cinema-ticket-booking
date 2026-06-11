<template>
  <div class="min-h-screen bg-stone-50">
    <main class="max-w-2xl mx-auto px-6 py-10">
      <div class="mb-8">
        <h1 class="text-xl font-semibold text-stone-800">My Tickets</h1>
        <p class="text-sm text-stone-400 mt-1">Manage and view your movie bookings</p>
      </div>

      <div class="flex p-1 bg-stone-200/50 rounded-xl mb-8">
        <button
          v-for="tab in ['upcoming', 'past']"
          :key="tab"
          @click="activeTab = tab"
          :class="[
            'flex-1 py-2 text-sm font-medium rounded-lg transition-all duration-200 capitalize',
            activeTab === tab
              ? 'bg-white text-stone-800 shadow-sm'
              : 'text-stone-400 hover:text-stone-600',
          ]"
        >
          {{ tab }}
        </button>
      </div>

      <div v-if="filteredBookings.length" class="space-y-6">
        <div v-for="ticket in filteredBookings" :key="ticket.id" class="group relative">
          <div
            class="bg-white border border-stone-100 rounded-2xl overflow-hidden shadow-sm transition-hover duration-300 hover:shadow-md"
          >
            <div class="flex">
              <div class="w-24 bg-stone-100 shrink-0 relative overflow-hidden">
                <div class="absolute inset-0 flex items-center justify-center text-stone-300">
                  <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path
                      d="M7 4v16M17 4v16M3 8h4m10 0h4M3 12h18M3 16h4m10 0h4M4 20h16a2 2 0 0 0 2-2V6a2 2 0 0 0-2-2H4a2 2 0 0 0-2 2v12a2 2 0 0 0 2 2z"
                      stroke-width="1.5"
                    />
                  </svg>
                </div>
                <img v-if="ticket.poster" :src="ticket.poster" class="w-full h-full object-cover" />
              </div>

              <div class="flex-1 p-5">
                <div class="flex justify-between items-start mb-2">
                  <h3 class="font-semibold text-stone-800 leading-tight">{{ ticket.movie }}</h3>
                  <span
                    :class="[
                      'text-[10px] uppercase tracking-wider px-2 py-0.5 rounded-full font-bold',
                      ticket.status === 'confirmed'
                        ? 'bg-green-50 text-green-600'
                        : 'bg-stone-100 text-stone-400',
                    ]"
                  >
                    {{ ticket.status }}
                  </span>
                </div>

                <div class="space-y-1.5 mb-4">
                  <div class="flex items-center gap-2 text-xs text-stone-500">
                    <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path
                        d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
                        stroke-width="1.5"
                      />
                    </svg>
                    {{ ticket.date }} · {{ ticket.time }}
                  </div>
                  <div class="flex items-center gap-2 text-xs text-stone-500">
                    <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path
                        d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"
                        stroke-width="1.5"
                      />
                      <path d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" stroke-width="1.5" />
                    </svg>
                    {{ ticket.cinema }}
                  </div>
                </div>

                <div
                  class="border-t border-dashed border-stone-100 pt-4 flex items-center justify-between"
                >
                  <div class="flex gap-1.5">
                    <span
                      v-for="seat in ticket.seats"
                      :key="seat"
                      class="text-[11px] font-medium text-stone-600 px-2 py-0.5 bg-stone-50 border border-stone-100 rounded"
                    >
                      {{ seat }}
                    </span>
                  </div>
                  <p class="text-[10px] text-stone-300 font-mono">#{{ ticket.bookingId }}</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div v-else class="flex flex-col items-center justify-center py-20 text-center">
        <div class="w-16 h-16 bg-stone-100 rounded-full flex items-center justify-center mb-4">
          <svg class="w-8 h-8 text-stone-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path
              d="M15 5v2m0 4v2m0 4v2M5 5a2 2 0 00-2 2v3a2 2 0 110 4v3a2 2 0 002 2h14a2 2 0 002-2v-3a2 2 0 110-4V7a2 2 0 00-2-2H5z"
              stroke-width="1.5"
            />
          </svg>
        </div>
        <p class="text-stone-400 text-sm">No {{ activeTab }} tickets found.</p>
        <button
          @click="$router.push({ name: 'movies' })"
          class="mt-4 text-sm font-medium text-stone-800 hover:text-stone-600"
        >
          Book a movie now →
        </button>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const activeTab = ref('upcoming')

// Mock data for tickets history
const bookings = ref([
  {
    id: 1,
    bookingId: 'BK-9942',
    movie: 'Dune: Part Two',
    date: 'Today, 24 Oct',
    time: '14:30',
    cinema: 'CinePass Central · Screen 4',
    seats: ['E5', 'E6'],
    status: 'confirmed',
    type: 'upcoming',
  },
  {
    id: 2,
    bookingId: 'BK-8821',
    movie: 'Oppenheimer',
    date: '15 Oct 2023',
    time: '18:00',
    cinema: 'CinePass Central · Screen 1',
    seats: ['B12'],
    status: 'completed',
    type: 'past',
  },
  {
    id: 3,
    bookingId: 'BK-7750',
    movie: 'Past Lives',
    date: '10 Oct 2023',
    time: '13:00',
    cinema: 'CinePass EmQuartier · Screen 2',
    seats: ['F1', 'F2'],
    status: 'completed',
    type: 'past',
  },
])

const filteredBookings = computed(() => {
  return bookings.value.filter((b) => b.type === activeTab.value)
})
</script>

<style scoped>
.border-dashed {
  background-image: linear-gradient(to right, #e7e5e4 50%, rgba(255, 255, 255, 0) 0%);
  background-position: bottom;
  background-size: 8px 1px;
  background-repeat: repeat-x;
  border-top: none;
}
</style>
