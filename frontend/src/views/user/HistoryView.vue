<template>
  <div class="min-h-screen bg-stone-50">
    <main class="max-w-2xl mx-auto px-6 py-10">
      <div class="mb-8">
        <h1 class="text-xl font-semibold text-stone-800">My Tickets</h1>
        <p class="text-sm text-stone-400 mt-1">Manage and view your movie bookings</p>
      </div>

      <div class="flex p-1 bg-stone-200/50 rounded-xl mb-8">
        <button
          v-for="tab in ['active', 'past']"
          :key="tab"
          @click="activeTab = tab"
          :class="[
            'flex-1 py-2 text-sm font-medium rounded-lg transition-all duration-200 capitalize',
            activeTab === tab ? 'bg-white text-stone-800 shadow-sm' : 'text-stone-400 hover:text-stone-600',
          ]"
        >
          {{ tab }}
        </button>
      </div>

      <div v-if="loading" class="py-16 text-center text-sm text-stone-400">Loading tickets...</div>

      <div v-else-if="filteredBookings.length" class="space-y-6">
        <div v-for="ticket in filteredBookings" :key="ticket.id" class="bg-white border border-stone-100 rounded-2xl overflow-hidden shadow-sm">
          <div class="flex">
            <div class="w-24 bg-sky-50 shrink-0 relative overflow-hidden flex items-center justify-center">
              <svg class="w-8 h-8 text-sky-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path
                  d="M7 4v16M17 4v16M3 8h4m10 0h4M3 12h18M3 16h4m10 0h4M4 20h16a2 2 0 0 0 2-2V6a2 2 0 0 0-2-2H4a2 2 0 0 0-2 2v12a2 2 0 0 0 2 2z"
                  stroke-width="1.5"
                />
              </svg>
            </div>

            <div class="flex-1 p-5">
              <div class="flex justify-between items-start mb-2 gap-3">
                <h3 class="font-semibold text-stone-800 leading-tight">Showtime {{ shortId(ticket.showtime_id) }}</h3>
                <span
                  :class="[
                    'text-[10px] uppercase tracking-wider px-2 py-0.5 rounded-full font-bold',
                    displayStatus(ticket) === 'SUCCESS'
                      ? 'bg-green-50 text-green-600'
                      : displayStatus(ticket) === 'PENDING'
                        ? 'bg-amber-50 text-amber-600'
                        : 'bg-stone-100 text-stone-400',
                  ]"
                >
                  {{ displayStatus(ticket) }}
                </span>
              </div>

              <div class="space-y-1.5 mb-4">
                <div class="flex items-center gap-2 text-xs text-stone-500">
                  <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" stroke-width="1.5" />
                  </svg>
                  {{ formatDate(ticket.created_at) }}
                </div>
                <div class="flex items-center gap-2 text-xs text-stone-500">
                  <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path d="M12 6v6l4 2" stroke-width="1.5" stroke-linecap="round" />
                    <circle cx="12" cy="12" r="10" stroke-width="1.5" />
                  </svg>
                  ฿{{ ticket.total_price }}
                </div>
              </div>

              <div class="border-t border-dashed border-stone-100 pt-4 flex items-center justify-between gap-3">
                <div class="flex gap-1.5 flex-wrap">
                  <span
                    v-for="seat in ticket.seats"
                    :key="seat"
                    class="text-[11px] font-medium text-stone-600 px-2 py-0.5 bg-stone-50 border border-stone-100 rounded"
                  >
                    {{ seat }}
                  </span>
                </div>
                <div class="flex items-center gap-3 shrink-0">
                  <button
                    v-if="isPayable(ticket)"
                    @click="continuePayment(ticket)"
                    class="px-3 py-1.5 rounded-lg bg-stone-900 text-white text-xs font-medium hover:bg-stone-700 transition-colors"
                  >
                    Continue payment
                  </button>
                  <p class="text-[10px] text-stone-300 font-mono">#{{ shortId(ticket.id) }}</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div v-else class="flex flex-col items-center justify-center py-20 text-center">
        <div class="w-16 h-16 bg-stone-100 rounded-full flex items-center justify-center mb-4">
          <svg class="w-8 h-8 text-stone-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path d="M15 5v2m0 4v2m0 4v2M5 5a2 2 0 00-2 2v3a2 2 0 110 4v3a2 2 0 002 2h14a2 2 0 002-2v-3a2 2 0 110-4V7a2 2 0 00-2-2H5z" stroke-width="1.5" />
          </svg>
        </div>
        <p class="text-stone-400 text-sm">{{ error || `No ${activeTab} tickets found.` }}</p>
        <button @click="$router.push({ name: 'movies' })" class="mt-4 text-sm font-medium text-stone-800 hover:text-stone-600">
          Book a movie now
        </button>
      </div>
    </main>
  </div>
</template>

<script setup>
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { apiFetch } from '../../api/client'

const router = useRouter()
const activeTab = ref('active')
const bookings = ref([])
const loading = ref(true)
const error = ref('')
const now = ref(Date.now())
let timerInterval = null

const filteredBookings = computed(() => {
  return bookings.value.filter((booking) => {
    const isPast = booking.status === 'TIMEOUT' || isExpiredPending(booking)
    return activeTab.value === 'past' ? isPast : !isPast
  })
})

const formatDate = (value) =>
  new Intl.DateTimeFormat('en-US', { dateStyle: 'medium', timeStyle: 'short' }).format(new Date(value))

const shortId = (id) => String(id || '').slice(-6)

const secondsUntil = (value) => {
  if (!value) return 0
  return Math.max(0, Math.floor((new Date(value).getTime() - now.value) / 1000))
}

const isExpiredPending = (booking) => {
  return booking.status === 'PENDING' && secondsUntil(booking.lock_expires_at) <= 0
}

const displayStatus = (booking) => {
  return isExpiredPending(booking) ? 'TIMEOUT' : booking.status
}

const isPayable = (booking) => {
  return booking.status === 'PENDING' && secondsUntil(booking.lock_expires_at) > 0
}

const continuePayment = (booking) => {
  router.push({
    name: 'payment',
    query: { bookingId: booking.id },
    state: {
      bookingId: booking.id,
      movie: `Showtime ${shortId(booking.showtime_id)}`,
      time: formatDate(booking.created_at),
      cinema: 'CinePass Central',
      seats: booking.seats,
      pricePerSeat: booking.seats.length ? booking.total_price / booking.seats.length : 0,
      lockExpiresAt: booking.lock_expires_at,
    },
  })
}

onMounted(async () => {
  timerInterval = setInterval(() => {
    now.value = Date.now()
  }, 1000)

  try {
    bookings.value = await apiFetch('/bookings/my')
  } catch (err) {
    error.value = err.message
  } finally {
    loading.value = false
  }
})

onUnmounted(() => clearInterval(timerInterval))
</script>
