<template>
  <div class="min-h-screen bg-stone-50">
    <main class="max-w-3xl mx-auto px-6 py-10">
      <div class="mb-8">
        <p class="text-xs text-stone-400 uppercase tracking-widest mb-1">Now Showing</p>
        <h1 class="text-xl font-semibold text-stone-800">{{ movie?.title || 'Loading...' }}</h1>
        <p class="text-sm text-stone-400 mt-1">
          {{ selectedShowtimeLabel }} <span v-if="movie">&middot; {{ formatDuration(movie.duration_mins) }}</span>
        </p>
      </div>

      <div v-if="showtimes.length > 1" class="flex gap-2 overflow-x-auto pb-2 mb-6">
        <button
          v-for="showtime in showtimes"
          :key="showtime.id"
          @click="selectShowtime(showtime.id)"
          :class="[
            'px-4 py-2 rounded-lg border text-sm shrink-0 transition-colors',
            selectedShowtimeId === showtime.id
              ? 'bg-stone-900 text-white border-stone-900'
              : 'bg-white text-stone-500 border-stone-200 hover:border-stone-400',
          ]"
        >
          {{ formatDate(showtime.start_time) }}
        </button>
      </div>

      <div class="mb-6 text-center">
        <div class="h-1 w-2/3 mx-auto rounded-full bg-stone-200 mb-1"></div>
        <p class="text-[11px] text-stone-300 tracking-widest uppercase">Screen</p>
      </div>

      <div class="bg-white border border-stone-100 rounded-2xl p-6 mb-6">
        <div v-if="loading" class="py-16 text-center text-sm text-stone-400">Loading seats...</div>
        <div v-else-if="seatRows.length" class="space-y-2.5">
          <div v-for="row in seatRows" :key="row.label" class="flex items-center gap-2">
            <span class="text-xs text-stone-300 w-4 text-center shrink-0">{{ row.label }}</span>
            <div class="flex gap-2 flex-1 justify-center">
              <template v-for="(seat, index) in row.seats" :key="seat.seat_no">
                <div v-if="index === 4" class="w-4 shrink-0"></div>
                <button
                  :disabled="seat.status !== 'AVAILABLE'"
                  @click="toggleSeat(seat.seat_no)"
                  :class="[
                    'w-9 h-9 rounded-md text-[11px] font-medium transition-all duration-150',
                    seat.status === 'BOOKED'
                      ? 'bg-stone-100 text-stone-300 cursor-not-allowed'
                      : seat.status === 'LOCKED'
                        ? 'bg-amber-100 text-amber-500 cursor-not-allowed'
                        : selectedSeatNos.includes(seat.seat_no)
                          ? 'bg-stone-800 text-white scale-105'
                          : 'bg-stone-50 text-stone-400 border border-stone-200 hover:border-stone-400 hover:text-stone-700',
                  ]"
                >
                  {{ seat.seat_no.slice(1) }}
                </button>
              </template>
            </div>
          </div>
        </div>
        <div v-else class="py-16 text-center text-sm text-stone-400">No showtime available.</div>

        <div class="flex items-center justify-center gap-6 mt-6 pt-5 border-t border-stone-100 flex-wrap">
          <div class="flex items-center gap-2">
            <div class="w-4 h-4 rounded bg-stone-50 border border-stone-200"></div>
            <span class="text-xs text-stone-400">Available</span>
          </div>
          <div class="flex items-center gap-2">
            <div class="w-4 h-4 rounded bg-stone-800"></div>
            <span class="text-xs text-stone-400">Selected</span>
          </div>
          <div class="flex items-center gap-2">
            <div class="w-4 h-4 rounded bg-amber-100"></div>
            <span class="text-xs text-stone-400">Locked</span>
          </div>
          <div class="flex items-center gap-2">
            <div class="w-4 h-4 rounded bg-stone-100"></div>
            <span class="text-xs text-stone-400">Booked</span>
          </div>
        </div>
      </div>

      <div class="bg-white border border-stone-100 rounded-2xl p-5">
        <div class="flex items-center justify-between mb-4 gap-4">
          <div>
            <p class="text-xs text-stone-400 mb-1">Selected seats</p>
            <p class="text-sm font-medium text-stone-800">
              <span v-if="selectedSeatNos.length">{{ selectedSeatNos.join(', ') }}</span>
              <span v-else class="text-stone-300">None</span>
            </p>
          </div>
          <div class="text-right">
            <p class="text-xs text-stone-400 mb-1">Total</p>
            <p class="text-lg font-semibold text-stone-800">฿{{ totalPrice }}</p>
          </div>
        </div>

        <p v-if="error" class="mb-4 text-xs text-red-500">{{ error }}</p>

        <button
          @click="handleBook"
          :disabled="!selectedSeatNos.length || loadingBooking || !selectedShowtime"
          class="w-full py-3 rounded-xl text-sm font-medium transition-all duration-150 active:scale-[0.99] disabled:opacity-40 disabled:cursor-not-allowed bg-stone-800 text-white hover:bg-stone-700"
        >
          <span
            v-if="loadingBooking"
            class="inline-block w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin mr-2 align-middle"
          ></span>
          {{ loadingBooking ? 'Reserving seats...' : `Book ${selectedSeatNos.length || ''} seat${selectedSeatNos.length !== 1 ? 's' : ''}` }}
        </button>
      </div>
    </main>
  </div>
</template>

<script setup>
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { apiFetch, websocketUrl } from '../../api/client'

const route = useRoute()
const router = useRouter()

const movie = ref(null)
const showtimes = ref([])
const selectedShowtimeId = ref('')
const selectedSeatNos = ref([])
const loading = ref(true)
const loadingBooking = ref(false)
const error = ref('')
let socket = null

const selectedShowtime = computed(() =>
  showtimes.value.find((showtime) => showtime.id === selectedShowtimeId.value),
)

const selectedShowtimeLabel = computed(() =>
  selectedShowtime.value ? formatDate(selectedShowtime.value.start_time) : 'Choose a showtime',
)

const seatRows = computed(() => {
  if (!selectedShowtime.value?.seats) return []
  const rows = new Map()

  selectedShowtime.value.seats.forEach((seat) => {
    const rowLabel = seat.seat_no.slice(0, 1)
    if (!rows.has(rowLabel)) rows.set(rowLabel, [])
    rows.get(rowLabel).push(seat)
  })

  return [...rows.entries()]
    .sort(([a], [b]) => a.localeCompare(b))
    .map(([label, seats]) => ({
      label,
      seats: seats.sort((a, b) => Number(a.seat_no.slice(1)) - Number(b.seat_no.slice(1))),
    }))
})

const totalPrice = computed(() => {
  return selectedSeatNos.value.length * (selectedShowtime.value?.price_per_seat || 0)
})

const formatDuration = (minutes) => {
  const hours = Math.floor(minutes / 60)
  const mins = minutes % 60
  return hours ? `${hours}h ${mins}m` : `${mins}m`
}

const formatDate = (value) => {
  return new Intl.DateTimeFormat('en-US', {
    dateStyle: 'medium',
    timeStyle: 'short',
  }).format(new Date(value))
}

const selectShowtime = (id) => {
  selectedShowtimeId.value = id
  selectedSeatNos.value = []
}

const toggleSeat = (seatNo) => {
  if (selectedSeatNos.value.includes(seatNo)) {
    selectedSeatNos.value = selectedSeatNos.value.filter((seat) => seat !== seatNo)
  } else {
    selectedSeatNos.value = [...selectedSeatNos.value, seatNo].sort()
  }
}

const applySeatUpdates = (message) => {
  if (message.type !== 'SEAT_UPDATED' || message.showtime_id !== selectedShowtimeId.value) return
  const updates = message.payload?.seats || []
  const showtime = selectedShowtime.value
  if (!showtime) return

  updates.forEach((updatedSeat) => {
    const target = showtime.seats.find((seat) => seat.seat_no === updatedSeat.seat_no)
    if (target) {
      target.status = updatedSeat.status
      target.locked_by = updatedSeat.locked_by
      target.locked_until = updatedSeat.locked_until
    }
    if (updatedSeat.status !== 'AVAILABLE') {
      selectedSeatNos.value = selectedSeatNos.value.filter((seat) => seat !== updatedSeat.seat_no)
    }
  })
}

const connectWebSocket = () => {
  socket?.close()
  socket = new WebSocket(websocketUrl())
  socket.onmessage = (event) => {
    try {
      applySeatUpdates(JSON.parse(event.data))
    } catch (err) {
      console.error(err)
    }
  }
}

const handleBook = async () => {
  if (!selectedSeatNos.value.length || !selectedShowtime.value) return
  loadingBooking.value = true
  error.value = ''
  try {
    const response = await apiFetch('/bookings/lock', {
      method: 'POST',
      body: {
        showtime_id: selectedShowtime.value.id,
        seats: selectedSeatNos.value,
      },
    })
    const booking = response.booking
    router.push({
      name: 'payment',
      query: { bookingId: booking.id },
      state: {
        bookingId: booking.id,
        movie: movie.value.title,
        time: selectedShowtimeLabel.value,
        cinema: 'CinePass Central',
        seats: booking.seats,
        pricePerSeat: selectedShowtime.value.price_per_seat,
        lockExpiresAt: booking.lock_expires_at,
      },
    })
  } catch (err) {
    error.value = err.message
    await loadShowtimes()
  } finally {
    loadingBooking.value = false
  }
}

const loadShowtimes = async () => {
  showtimes.value = await apiFetch(`/showtimes/movie/${route.params.movieId}`)
  if (!selectedShowtimeId.value && showtimes.value.length) {
    selectedShowtimeId.value = showtimes.value[0].id
  }
}

onMounted(async () => {
  try {
    const [movieData] = await Promise.all([
      apiFetch(`/movies/${route.params.movieId}`),
      loadShowtimes(),
    ])
    movie.value = movieData
    connectWebSocket()
  } catch (err) {
    error.value = err.message
  } finally {
    loading.value = false
  }
})

watch(selectedShowtimeId, () => {
  selectedSeatNos.value = []
})

onUnmounted(() => {
  socket?.close()
})
</script>
