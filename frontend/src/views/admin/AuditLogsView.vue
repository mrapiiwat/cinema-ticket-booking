<template>
  <main class="p-8 max-w-7xl w-full mx-auto">
    <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-6">
      <div>
        <div class="flex items-center gap-2">
          <h2 class="text-lg font-semibold text-stone-800">System Logs</h2>
          <span class="flex h-2 w-2 relative">
            <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
            <span class="relative inline-flex rounded-full h-2 w-2 bg-emerald-500"></span>
          </span>
          <span class="text-[10px] text-emerald-600 font-medium uppercase tracking-wider">Live</span>
        </div>
        <p class="text-sm text-stone-500 mt-1">Track booking, timeout, release, and system events.</p>
      </div>
      <button
        @click="refreshLogs"
        class="flex items-center gap-2 px-4 py-2 bg-white border border-stone-200 text-stone-600 rounded-xl text-sm font-medium hover:bg-stone-50 transition-colors"
      >
        <svg :class="{ 'animate-spin': isRefreshing }" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
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

    <div class="bg-white border border-stone-100 rounded-2xl p-5 mb-6 shadow-sm flex flex-col sm:flex-row gap-4 items-center justify-between">
      <div class="flex flex-col sm:flex-row gap-3 w-full sm:w-auto flex-1">
        <div class="relative flex-1 max-w-sm">
          <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-stone-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <circle cx="11" cy="11" r="8" stroke-width="1.5" />
            <path d="M21 21l-4.35-4.35" stroke-width="1.5" stroke-linecap="round" />
          </svg>
          <input
            v-model="filters.search"
            type="text"
            placeholder="Search message or user..."
            class="w-full pl-9 pr-4 py-2 text-sm bg-stone-50 border border-stone-200 rounded-xl text-stone-700 placeholder-stone-400 focus:outline-none focus:border-stone-400 focus:bg-white transition-all"
          />
        </div>

        <select v-model="filters.type" class="px-3 py-2 text-sm bg-stone-50 border border-stone-200 rounded-xl text-stone-600 focus:outline-none focus:border-stone-400 focus:bg-white transition-all">
          <option value="">All Events</option>
          <option value="BOOKING_SUCCESS">Booking Success</option>
          <option value="BOOKING_TIMEOUT">Booking Timeout</option>
          <option value="SEAT_RELEASED">Seat Released</option>
          <option value="SEAT_LOCKED">Seat Locked</option>
          <option value="SYSTEM_ERROR">System Error</option>
        </select>
      </div>
      <p class="text-xs text-stone-400 shrink-0">Showing {{ filteredLogs.length }} events</p>
    </div>

    <div class="bg-white border border-stone-100 rounded-2xl shadow-sm overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="border-b border-stone-100 bg-stone-50/50 text-stone-400 text-[11px] uppercase tracking-wider font-semibold">
              <th class="py-4 px-6 w-48">Timestamp</th>
              <th class="py-4 px-6 w-44">Event Type</th>
              <th class="py-4 px-6 w-32">User</th>
              <th class="py-4 px-6">Details</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-stone-50 text-sm text-stone-700">
            <tr v-for="log in filteredLogs" :key="log.id" class="hover:bg-stone-50/30 transition-colors font-mono text-[13px]">
              <td class="py-4 px-6 text-stone-500">{{ formatDate(log.timestamp) }}</td>
              <td class="py-4 px-6">
                <span :class="['px-2 py-0.5 rounded-md font-semibold border', eventClass(log.event_type)]">
                  {{ log.event_type }}
                </span>
              </td>
              <td class="py-4 px-6 text-stone-400">{{ shortId(log.user_id) || '-' }}</td>
              <td class="py-4 px-6 text-stone-600">{{ log.details }}</td>
            </tr>

            <tr v-if="!filteredLogs.length">
              <td colspan="4" class="text-center py-12 text-stone-400 text-sm font-sans">
                {{ error || 'No logs found matching your criteria.' }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </main>
</template>

<script setup>
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { apiFetch, websocketUrl } from '../../api/client'

const isRefreshing = ref(false)
const logs = ref([])
const error = ref('')
const filters = ref({ search: '', type: '' })
let socket = null

const filteredLogs = computed(() => {
  return logs.value.filter((log) => {
    const searchLower = filters.value.search.toLowerCase()
    const matchesSearch =
      !searchLower ||
      log.details.toLowerCase().includes(searchLower) ||
      (log.user_id || '').toLowerCase().includes(searchLower)
    const matchesType = !filters.value.type || log.event_type === filters.value.type
    return matchesSearch && matchesType
  })
})

const eventClass = (type) => {
  if (type === 'BOOKING_SUCCESS') return 'bg-green-50 text-green-600 border-green-100'
  if (type === 'BOOKING_TIMEOUT') return 'bg-amber-50 text-amber-600 border-amber-100'
  if (type === 'SEAT_RELEASED' || type === 'SEAT_LOCKED') return 'bg-sky-50 text-sky-600 border-sky-100'
  return 'bg-red-50 text-red-600 border-red-100'
}

const shortId = (id) => String(id || '').slice(-6)
const formatDate = (value) =>
  new Intl.DateTimeFormat('en-US', { dateStyle: 'medium', timeStyle: 'medium' }).format(new Date(value))

const refreshLogs = async () => {
  isRefreshing.value = true
  error.value = ''
  try {
    logs.value = await apiFetch('/admin/audit-logs')
  } catch (err) {
    error.value = err.message
  } finally {
    isRefreshing.value = false
  }
}

const handleLiveLogs = (message) => {
  if (message.type === 'AUDIT_LOG' || message.payload?.log) {
    const newLog = message.payload.log
    logs.value.unshift(newLog) 
  } else if (message.type === 'SEAT_UPDATED') {
    logs.value.unshift({
      id: crypto.randomUUID(),
      timestamp: new Date().toISOString(),
      event_type: message.payload?.seats[0]?.status === 'LOCKED' ? 'SEAT_LOCKED' : 'SEAT_RELEASED',
      user_id: message.payload?.seats[0]?.locked_by || 'system',
      details: `Seats ${message.payload?.seats.map(s => s.seat_no).join(', ')} status changed.`
    })
  }
}

const connectWebSocket = () => {
  socket?.close()
  socket = new WebSocket(websocketUrl())
  socket.onmessage = (event) => {
    try {
      handleLiveLogs(JSON.parse(event.data))
    } catch (err) {
      console.error('WS Error:', err)
    }
  }
}

onMounted(async () => {
  await refreshLogs()
  connectWebSocket()
})

onUnmounted(() => {
  socket?.close()
})
</script>
