<template>
  <div class="min-h-screen bg-stone-50">
    <main class="max-w-2xl mx-auto px-6 py-10">
      <div class="mb-8">
        <p class="text-xs text-stone-400 uppercase tracking-widest mb-1">Now Showing</p>
        <h1 class="text-xl font-semibold text-stone-800">{{ movie.title }}</h1>
        <p class="text-sm text-stone-400 mt-1">{{ movie.time }} &middot; {{ movie.duration }}</p>
      </div>

      <div class="mb-6 text-center">
        <div class="h-0.75 w-2/3 mx-auto rounded-full bg-stone-200 mb-1"></div>
        <p class="text-[11px] text-stone-300 tracking-widest uppercase">Screen</p>
      </div>

      <div class="bg-white border border-stone-100 rounded-2xl p-6 mb-6">
        <div class="space-y-2.5">
          <div v-for="(row, rowIndex) in seats" :key="rowIndex" class="flex items-center gap-2">
            <span class="text-xs text-stone-300 w-4 text-center shrink-0">{{
              rowLabels[rowIndex]
            }}</span>
            <div class="flex gap-2 flex-1 justify-center">
              <template v-for="(seat, seatIndex) in row" :key="seatIndex">
                <div v-if="seatIndex === 4" class="w-4 shrink-0"></div>
                <button
                  :disabled="seat.status === 'booked' || seat.status === 'locked'"
                  @click="toggleSeat(rowIndex, seatIndex)"
                  :class="[
                    'w-8 h-8 rounded-md text-[11px] font-medium transition-all duration-150',
                    seat.status === 'booked'
                      ? 'bg-stone-100 text-stone-300 cursor-not-allowed'
                      : seat.status === 'locked'
                        ? 'bg-amber-100 text-amber-400 cursor-not-allowed'
                        : seat.selected
                          ? 'bg-stone-800 text-white scale-105'
                          : 'bg-stone-50 text-stone-400 border border-stone-200 hover:border-stone-400 hover:text-stone-700',
                  ]"
                >
                  {{ seatIndex + 1 }}
                </button>
              </template>
            </div>
          </div>
        </div>

        <div class="flex items-center justify-center gap-6 mt-6 pt-5 border-t border-stone-100">
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
        <div class="flex items-center justify-between mb-4">
          <div>
            <p class="text-xs text-stone-400 mb-1">Selected seats</p>
            <p class="text-sm font-medium text-stone-800">
              <span v-if="selectedSeats.length">{{
                selectedSeats.map((s) => s.label).join(', ')
              }}</span>
              <span v-else class="text-stone-300">None</span>
            </p>
          </div>
          <div class="text-right">
            <p class="text-xs text-stone-400 mb-1">Total</p>
            <p class="text-lg font-semibold text-stone-800">฿{{ totalPrice }}</p>
          </div>
        </div>

        <div
          v-if="timerActive"
          class="flex items-center gap-2 mb-4 px-3 py-2 bg-amber-50 rounded-lg border border-amber-100"
        >
          <svg
            class="w-4 h-4 text-amber-400 shrink-0"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <circle cx="12" cy="12" r="10" stroke-width="1.5" />
            <path d="M12 6v6l4 2" stroke-width="1.5" stroke-linecap="round" />
          </svg>
          <p class="text-xs text-amber-600">
            Seats reserved for <span class="font-semibold">{{ formattedTimer }}</span> — complete
            payment to confirm
          </p>
        </div>

        <button
          @click="handleBook"
          :disabled="!selectedSeats.length || loading"
          class="w-full py-3 rounded-xl text-sm font-medium transition-all duration-150 active:scale-[0.99] disabled:opacity-40 disabled:cursor-not-allowed bg-stone-800 text-white hover:bg-stone-700"
        >
          <span
            v-if="loading"
            class="inline-block w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin mr-2 align-middle"
          ></span>
          {{
            loading
              ? 'Reserving seats...'
              : `Book ${selectedSeats.length || ''} seat${selectedSeats.length !== 1 ? 's' : ''}`
          }}
        </button>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, computed, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

// Mock movie data — replace with API call using movieId
const movie = ref({
  title: 'Dune: Part Two',
  time: 'Today 14:30',
  duration: '2h 46m',
  pricePerSeat: 220,
})

const rowLabels = ['A', 'B', 'C', 'D', 'E', 'F']

// Generate seat map: 6 rows x 8 seats
// status: 'available' | 'locked' | 'booked'
// ลบ parameter _, r, s ที่ไม่ได้ใช้ออกไป
const seats = ref(
  rowLabels.map(() =>
    Array.from({ length: 8 }, () => ({
      status: Math.random() < 0.15 ? 'booked' : Math.random() < 0.1 ? 'locked' : 'available',
      selected: false,
    })),
  ),
)

const selectedSeats = computed(() => {
  const result = []
  seats.value.forEach((row, r) =>
    row.forEach((seat, s) => {
      if (seat.selected) result.push({ r, s, label: `${rowLabels[r]}${s + 1}` })
    }),
  )
  return result
})

const totalPrice = computed(() => selectedSeats.value.length * movie.value.pricePerSeat)

const toggleSeat = (r, s) => {
  const seat = seats.value[r][s]
  if (seat.status !== 'available') return
  seat.selected = !seat.selected
}

// Timer
const timerActive = ref(false)
const timeLeft = ref(300) // 5 minutes in seconds
let timerInterval = null

const formattedTimer = computed(() => {
  const m = Math.floor(timeLeft.value / 60)
  const s = timeLeft.value % 60
  return `${m}:${s.toString().padStart(2, '0')}`
})

const startTimer = () => {
  timerActive.value = true
  timeLeft.value = 300
  timerInterval = setInterval(() => {
    timeLeft.value--
    if (timeLeft.value <= 0) {
      clearInterval(timerInterval)
      timerActive.value = false
      // Release selected seats back to available
      seats.value.forEach((row) =>
        row.forEach((seat) => {
          seat.selected = false
        }),
      )
    }
  }, 1000)
}

// Booking
const loading = ref(false)

const handleBook = async () => {
  if (!selectedSeats.value.length) return
  loading.value = true
  try {
    // Backend will acquire Redis distributed lock here
    const res = await new Promise((resolve) =>
      setTimeout(() => resolve({ bookingId: 'BK-0001' }), 800),
    )

    startTimer()

    // Push to payment พร้อมส่ง state ไปด้วย
    router.push({
      name: 'payment',
      state: {
        bookingId: res.bookingId,
        movie: movie.value.title,
        time: movie.value.time,
        cinema: 'CinePass Central',
        seats: selectedSeats.value.map((s) => s.label),
        pricePerSeat: movie.value.pricePerSeat,
        timeLeft: timeLeft.value,
      },
    })
  } catch (err) {
    console.error('Booking failed:', err)
  } finally {
    loading.value = false
  }
}

onUnmounted(() => clearInterval(timerInterval))
</script>
