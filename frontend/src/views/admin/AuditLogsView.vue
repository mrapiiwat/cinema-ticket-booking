<template>
  <main class="p-8 max-w-7xl w-full mx-auto">
    <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-6">
      <div>
        <h2 class="text-lg font-semibold text-stone-800">System Logs</h2>
        <p class="text-sm text-stone-500 mt-1">Track system events, lock releases, and errors.</p>
      </div>
      <button
        @click="refreshLogs"
        class="flex items-center gap-2 px-4 py-2 bg-white border border-stone-200 text-stone-600 rounded-xl text-sm font-medium hover:bg-stone-50 transition-colors"
      >
        <svg
          :class="{ 'animate-spin': isRefreshing }"
          class="w-4 h-4"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="1.5"
            d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
          />
        </svg>
        Refresh
      </button>
    </div>

    <div
      class="bg-white border border-stone-100 rounded-2xl p-5 mb-6 shadow-sm flex flex-col sm:flex-row gap-4 items-center justify-between"
    >
      <div class="flex flex-col sm:flex-row gap-3 w-full sm:w-auto flex-1">
        <div class="relative flex-1 max-w-sm">
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
            v-model="filters.search"
            type="text"
            placeholder="Search message or booking ID..."
            class="w-full pl-9 pr-4 py-2 text-sm bg-stone-50 border border-stone-200 rounded-xl text-stone-700 placeholder-stone-400 focus:outline-none focus:border-stone-400 focus:bg-white transition-all"
          />
        </div>

        <select
          v-model="filters.type"
          class="px-3 py-2 text-sm bg-stone-50 border border-stone-200 rounded-xl text-stone-600 focus:outline-none focus:border-stone-400 focus:bg-white transition-all"
        >
          <option value="">All Events</option>
          <option value="BOOKING_SUCCESS">Booking Success</option>
          <option value="BOOKING_TIMEOUT">Booking Timeout</option>
          <option value="SEAT_RELEASED">Seat Released</option>
          <option value="SYSTEM_ERROR">System Error</option>
        </select>
      </div>
      <p class="text-xs text-stone-400 shrink-0">Showing {{ filteredLogs.length }} events</p>
    </div>

    <div class="bg-white border border-stone-100 rounded-2xl shadow-sm overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr
              class="border-b border-stone-100 bg-stone-50/50 text-stone-400 text-[11px] uppercase tracking-wider font-semibold"
            >
              <th class="py-4 px-6 w-48">Timestamp</th>
              <th class="py-4 px-6 w-40">Event Type</th>
              <th class="py-4 px-6 w-32">Booking ID</th>
              <th class="py-4 px-6">Message / Details</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-stone-50 text-sm text-stone-700">
            <tr
              v-for="log in filteredLogs"
              :key="log.id"
              class="hover:bg-stone-50/30 transition-colors font-mono text-[13px]"
            >
              <td class="py-4 px-6 text-stone-500">
                {{ log.timestamp }}
              </td>

              <td class="py-4 px-6">
                <span
                  :class="[
                    'px-2 py-0.5 rounded-md font-semibold border',
                    log.type === 'BOOKING_SUCCESS'
                      ? 'bg-green-50 text-green-600 border-green-100'
                      : log.type === 'BOOKING_TIMEOUT'
                        ? 'bg-amber-50 text-amber-600 border-amber-100'
                        : log.type === 'SEAT_RELEASED'
                          ? 'bg-blue-50 text-blue-600 border-blue-100'
                          : 'bg-red-50 text-red-600 border-red-100',
                  ]"
                >
                  {{ log.type }}
                </span>
              </td>

              <td class="py-4 px-6">
                <span v-if="log.bookingId" class="text-stone-800 font-medium"
                  >#{{ log.bookingId }}</span
                >
                <span v-else class="text-stone-300">-</span>
              </td>

              <td class="py-4 px-6 text-stone-600">
                {{ log.message }}
                <div v-if="log.meta" class="mt-1 text-[11px] text-stone-400">
                  {{ JSON.stringify(log.meta) }}
                </div>
              </td>
            </tr>

            <tr v-if="!filteredLogs.length">
              <td colspan="4" class="text-center py-12 text-stone-400 text-sm font-sans">
                No logs found matching your criteria.
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </main>
</template>

<script setup>
import { ref, computed } from 'vue'

const isRefreshing = ref(false)

const filters = ref({
  search: '',
  type: '',
})

// Mock Logs Data (ตรงตาม requirement ของโจทย์: Success, Timeout, Released, Error)
const logs = ref([
  {
    id: 1,
    timestamp: '2023-10-24 14:35:12',
    type: 'BOOKING_SUCCESS',
    bookingId: 'BK-0001',
    message: 'Payment received and booking confirmed.',
    meta: { user: 'john@example.com', seats: ['E5', 'E6'] },
  },
  {
    id: 2,
    timestamp: '2023-10-24 14:32:05',
    type: 'SYSTEM_ERROR',
    bookingId: null,
    message: 'Failed to acquire Redis lock for seats A1, A2.',
    meta: { movieId: 'MV-002', error: 'Redis timeout' },
  },
  {
    id: 3,
    timestamp: '2023-10-24 14:30:00',
    type: 'SEAT_RELEASED',
    bookingId: 'BK-0003',
    message: 'Lock expired. Seats released back to pool.',
    meta: { seats: ['F1'] },
  },
  {
    id: 4,
    timestamp: '2023-10-24 14:29:59',
    type: 'BOOKING_TIMEOUT',
    bookingId: 'BK-0003',
    message: 'User failed to complete payment within 5 minutes.',
    meta: { user: 'bob.j@outlook.com' },
  },
  {
    id: 5,
    timestamp: '2023-10-24 14:15:22',
    type: 'BOOKING_SUCCESS',
    bookingId: 'BK-0004',
    message: 'Payment received and booking confirmed.',
    meta: { user: 'sarah@sky.net', seats: ['B3', 'B4', 'B5'] },
  },
])

const filteredLogs = computed(() => {
  return logs.value.filter((log) => {
    const searchLower = filters.value.search.toLowerCase()
    const matchesSearch =
      (log.bookingId && log.bookingId.toLowerCase().includes(searchLower)) ||
      log.message.toLowerCase().includes(searchLower) ||
      (log.meta && JSON.stringify(log.meta).toLowerCase().includes(searchLower))

    const matchesType = !filters.value.type || log.type === filters.value.type

    return matchesSearch && matchesType
  })
})

const refreshLogs = () => {
  isRefreshing.value = true
  // Mock API call delay
  setTimeout(() => {
    isRefreshing.value = false
  }, 600)
}
</script>
