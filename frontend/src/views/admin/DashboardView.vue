<template>
  <div class="min-h-screen bg-stone-50 flex">
    <div class="flex-1 flex flex-col min-w-0">
      <main class="flex-1 p-8 overflow-y-auto max-w-7xl w-full mx-auto">
        <div class="flex justify-between items-center mb-6">
          <div class="flex items-center gap-2">
            <h1 class="text-xl font-semibold text-stone-800">Admin Dashboard</h1>
            <span class="flex h-2 w-2 relative">
              <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
              <span class="relative inline-flex rounded-full h-2 w-2 bg-emerald-500"></span>
            </span>
            <span class="text-[10px] text-emerald-600 font-medium uppercase tracking-wider">Live Stats</span>
          </div>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-3 gap-5 mb-8">
          <div class="bg-white border border-stone-100 rounded-2xl p-5 shadow-sm">
            <p class="text-xs text-stone-400 tracking-wider mb-1">TOTAL REVENUE</p>
            <p class="text-2xl font-semibold text-stone-800">฿{{ stats.revenue.toLocaleString() }}</p>
          </div>
          <div class="bg-white border border-stone-100 rounded-2xl p-5 shadow-sm">
            <p class="text-xs text-stone-400 tracking-wider mb-1">CONFIRMED BOOKINGS</p>
            <p class="text-2xl font-semibold text-emerald-700">{{ stats.confirmedCount }} <span class="text-xs font-normal text-stone-400">orders</span></p>
          </div>
          <div class="bg-white border border-stone-100 rounded-2xl p-5 shadow-sm">
            <p class="text-xs text-amber-600 tracking-wider mb-1">ACTIVE REDIS LOCKS</p>
            <p class="text-2xl font-semibold text-amber-700">{{ stats.activeLocks }} <span class="text-xs font-normal text-amber-500">holding</span></p>
          </div>
        </div>

        <div class="bg-white border border-stone-100 rounded-2xl p-5 mb-6 shadow-sm flex flex-col sm:flex-row gap-4 items-center justify-between">
          <div class="flex flex-col sm:flex-row gap-3 w-full sm:w-auto flex-1">
            <div class="relative flex-1 max-w-xs">
              <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-stone-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <circle cx="11" cy="11" r="8" stroke-width="1.5" />
                <path d="M21 21l-4.35-4.35" stroke-width="1.5" stroke-linecap="round" />
              </svg>
              <input
                v-model="filters.search"
                type="text"
                placeholder="Search booking ID / user..."
                class="w-full pl-9 pr-4 py-2 text-sm bg-stone-50 border border-stone-200 rounded-xl text-stone-700 placeholder-stone-400 focus:outline-none focus:border-stone-400 focus:bg-white transition-all"
              />
            </div>

            <select v-model="filters.status" class="px-3 py-2 text-sm bg-stone-50 border border-stone-200 rounded-xl text-stone-600 focus:outline-none focus:border-stone-400 focus:bg-white transition-all">
              <option value="">All Status</option>
              <option value="SUCCESS">Success</option>
              <option value="PENDING">Pending</option>
              <option value="TIMEOUT">Timeout</option>
            </select>
          </div>

          <button @click="loadBookings" class="px-4 py-2 bg-stone-900 text-white rounded-xl text-sm font-medium hover:bg-stone-700 transition-colors">
            Force Refresh
          </button>
        </div>

        <div class="bg-white border border-stone-100 rounded-2xl shadow-sm overflow-hidden">
          <div class="overflow-x-auto">
            <table class="w-full text-left border-collapse">
              <thead>
                <tr class="border-b border-stone-100 bg-stone-50/50 text-stone-400 text-[11px] uppercase tracking-wider font-semibold">
                  <th class="py-4 px-6">Booking ID</th>
                  <th class="py-4 px-6">Showtime</th>
                  <th class="py-4 px-6">User</th>
                  <th class="py-4 px-6 text-center">Seats</th>
                  <th class="py-4 px-6 text-right">Total Price</th>
                  <th class="py-4 px-6 text-center">Status</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-stone-50 text-sm text-stone-700">
                <tr v-for="booking in filteredBookings" :key="booking.id" class="hover:bg-stone-50/30 transition-colors">
                  <td class="py-4 px-6 font-mono text-xs text-stone-400">#{{ shortId(booking.id) }}</td>
                  <td class="py-4 px-6">
                    <p class="font-medium text-stone-800">Showtime {{ shortId(booking.showtime_id) }}</p>
                    <p class="text-xs text-stone-400 mt-0.5">{{ formatDate(booking.created_at) }}</p>
                  </td>
                  <td class="py-4 px-6 font-mono text-[11px] text-stone-500">{{ shortId(booking.user_id) }}</td>
                  <td class="py-4 px-6 text-center">
                    <div class="flex gap-1 justify-center flex-wrap">
                      <span v-for="seat in booking.seats" :key="seat" class="px-1.5 py-0.5 bg-stone-100 text-stone-600 text-[10px] rounded font-medium">{{ seat }}</span>
                    </div>
                  </td>
                  <td class="py-4 px-6 text-right font-medium text-stone-800">฿{{ booking.total_price.toLocaleString() }}</td>
                  <td class="py-4 px-6 text-center">
                    <span
                      :class="[
                        'text-[10px] uppercase tracking-wider px-2 py-0.5 rounded-full font-semibold inline-block',
                        booking.status === 'SUCCESS'
                          ? 'bg-green-50 text-green-600'
                          : booking.status === 'PENDING'
                            ? 'bg-amber-50 text-amber-600 animate-pulse'
                            : 'bg-stone-100 text-stone-400',
                      ]"
                    >
                      {{ booking.status }}
                    </span>
                  </td>
                </tr>

                <tr v-if="!filteredBookings.length">
                  <td colspan="6" class="text-center py-12 text-stone-400 text-sm">{{ error || 'No booking records found.' }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </main>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { apiFetch, websocketUrl } from '../../api/client'

const filters = ref({ search: '', status: '' })
const bookings = ref([])
const error = ref('')
let socket = null

const stats = computed(() => {
  return bookings.value.reduce(
    (acc, booking) => {
      if (booking.status === 'SUCCESS') {
        acc.revenue += booking.total_price
        acc.confirmedCount += 1
      }
      if (booking.status === 'PENDING') acc.activeLocks += 1
      return acc
    },
    { revenue: 0, confirmedCount: 0, activeLocks: 0 },
  )
})

const filteredBookings = computed(() => {
  const search = filters.value.search.toLowerCase()
  return bookings.value.filter((booking) => {
    const matchesSearch =
      !search ||
      booking.id.toLowerCase().includes(search) ||
      booking.user_id.toLowerCase().includes(search) ||
      booking.showtime_id.toLowerCase().includes(search)
    const matchesStatus = !filters.value.status || booking.status === filters.value.status
    return matchesSearch && matchesStatus
  })
})

const shortId = (id) => String(id || '').slice(-6)
const formatDate = (value) =>
  new Intl.DateTimeFormat('en-US', { dateStyle: 'medium', timeStyle: 'short' }).format(new Date(value))

const loadBookings = async () => {
  error.value = ''
  try {
    bookings.value = await apiFetch('/admin/bookings')
  } catch (err) {
    error.value = err.message
  }
}

const applyLiveBookingUpdates = (message) => {
  if (message.type === 'BOOKING_UPDATED' || message.type === 'BOOKING_CREATED') {
    const updatedBooking = message.payload?.booking
    if (!updatedBooking) return

    const index = bookings.value.findIndex((b) => b.id === updatedBooking.id)
    if (index !== -1) {
      bookings.value[index] = updatedBooking
    } else {
      bookings.value.unshift(updatedBooking)
    }
  }
}

const connectWebSocket = () => {
  socket?.close()
  socket = new WebSocket(websocketUrl())
  socket.onmessage = (event) => {
    try {
      applyLiveBookingUpdates(JSON.parse(event.data))
    } catch (err) {
      console.error(err)
    }
  }
}

onMounted(async () => {
  await loadBookings()
  connectWebSocket()
})

onUnmounted(() => {
  socket?.close()
})
</script>
