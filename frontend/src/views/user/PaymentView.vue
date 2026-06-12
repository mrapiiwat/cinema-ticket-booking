<template>
  <div class="min-h-screen bg-stone-50">
    <main class="max-w-lg mx-auto px-6 py-10">
      <div
        class="flex items-center gap-2 mb-6 px-4 py-3 bg-amber-50 border border-amber-100 rounded-xl"
      >
        <svg class="w-4 h-4 text-amber-400 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <circle cx="12" cy="12" r="10" stroke-width="1.5" />
          <path d="M12 6v6l4 2" stroke-width="1.5" stroke-linecap="round" />
        </svg>
        <p class="text-xs text-amber-600 flex-1">
          Complete payment within <span class="font-semibold">{{ formattedTimer }}</span> or your seats will be released
        </p>
      </div>

      <div v-if="loadingOrder" class="py-20 text-center text-sm text-stone-400">Loading reservation...</div>

      <template v-else>
        <div class="bg-white border border-stone-100 rounded-2xl p-5 mb-4">
          <p class="text-xs text-stone-400 uppercase tracking-widest mb-4">Order Summary</p>

          <div class="flex items-start justify-between gap-4 pb-4 border-b border-stone-100">
            <div>
              <p class="text-sm font-medium text-stone-800">{{ order.movie }}</p>
              <p class="text-xs text-stone-400 mt-0.5">{{ order.time }} · {{ order.cinema }}</p>
              <div class="flex gap-1.5 mt-2 flex-wrap">
                <span
                  v-for="seat in order.seats"
                  :key="seat"
                  class="px-2 py-0.5 bg-stone-100 text-stone-500 text-[11px] rounded-md font-medium"
                >
                  {{ seat }}
                </span>
              </div>
            </div>
            <div class="text-right shrink-0">
              <p class="text-xs text-stone-400 mb-0.5">
                {{ order.seats.length }} seat{{ order.seats.length > 1 ? 's' : '' }}
              </p>
              <p class="text-sm font-semibold text-stone-800">
                ฿{{ order.seats.length * order.pricePerSeat }}
              </p>
            </div>
          </div>

          <div class="pt-4 space-y-2">
            <div class="flex justify-between text-sm">
              <span class="text-stone-400">Subtotal</span>
              <span class="text-stone-700">฿{{ order.seats.length * order.pricePerSeat }}</span>
            </div>
            <div class="flex justify-between text-sm">
              <span class="text-stone-400">Service fee</span>
              <span class="text-stone-700">฿{{ serviceFee }}</span>
            </div>
            <div class="flex justify-between text-sm font-semibold pt-2 border-t border-stone-100">
              <span class="text-stone-800">Total</span>
              <span class="text-stone-800">฿{{ grandTotal }}</span>
            </div>
          </div>
        </div>

        <div class="bg-white border border-stone-100 rounded-2xl p-5 mb-4">
          <p class="text-xs text-stone-400 uppercase tracking-widest mb-4">Payment Method</p>

          <div class="space-y-2.5">
            <label
              v-for="method in paymentMethods"
              :key="method.id"
              :class="[
                'flex items-center gap-3 p-3.5 rounded-xl border cursor-pointer transition-all duration-150',
                selected === method.id ? 'border-stone-800 bg-stone-50' : 'border-stone-100 hover:border-stone-200',
              ]"
            >
              <input type="radio" :value="method.id" v-model="selected" class="accent-stone-800" />
              <div>
                <p class="text-sm font-medium text-stone-700">{{ method.label }}</p>
                <p class="text-xs text-stone-400">{{ method.desc }}</p>
              </div>
            </label>
          </div>
        </div>

        <p v-if="error" class="mb-4 text-xs text-red-500">{{ error }}</p>

        <button
          @click="handlePay"
          :disabled="!selected || loading || paid || expired"
          class="w-full py-3.5 rounded-xl text-sm font-medium transition-all duration-150 active:scale-[0.99] disabled:opacity-40 disabled:cursor-not-allowed bg-stone-800 text-white hover:bg-stone-700"
        >
          <span
            v-if="loading"
            class="inline-block w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin mr-2 align-middle"
          ></span>
          <span v-if="paid">Payment Confirmed</span>
          <span v-else-if="expired">Reservation Expired</span>
          <span v-else-if="loading">Processing...</span>
          <span v-else>Pay ฿{{ grandTotal }}</span>
        </button>
      </template>
    </main>
  </div>
</template>

<script setup>
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { apiFetch } from '../../api/client'

const route = useRoute()
const router = useRouter()
const state = history.state

const order = ref({
  bookingId: state?.bookingId ?? '',
  movie: state?.movie ?? 'Unknown Movie',
  time: state?.time ?? '',
  cinema: state?.cinema ?? '',
  seats: state?.seats ?? [],
  pricePerSeat: state?.pricePerSeat ?? 0,
  lockExpiresAt: state?.lockExpiresAt ?? '',
})

const serviceFee = 30
const selected = ref('card')
const loadingOrder = ref(true)
const loading = ref(false)
const paid = ref(false)
const error = ref('')
const timeLeft = ref(0)
let timerInterval = null

const expired = computed(() => timeLeft.value <= 0)
const grandTotal = computed(() => order.value.seats.length * order.value.pricePerSeat + serviceFee)
const formattedTimer = computed(() => {
  const safe = Math.max(0, timeLeft.value)
  const minutes = Math.floor(safe / 60)
  const seconds = safe % 60
  return `${minutes}:${seconds.toString().padStart(2, '0')}`
})

const paymentMethods = [
  { id: 'card', label: 'Credit / Debit Card', desc: 'Visa, Mastercard, JCB' },
  { id: 'promptpay', label: 'PromptPay', desc: 'Scan QR to pay' },
  { id: 'wallet', label: 'TrueMoney Wallet', desc: 'Pay with e-wallet' },
]

function secondsUntil(value) {
  if (!value) return 0
  return Math.max(0, Math.floor((new Date(value).getTime() - Date.now()) / 1000))
}

const hydrateOrderFromBooking = (booking) => {
  order.value = {
    bookingId: booking.id,
    movie: `Showtime ${shortId(booking.showtime_id)}`,
    time: new Intl.DateTimeFormat('en-US', {
      dateStyle: 'medium',
      timeStyle: 'short',
    }).format(new Date(booking.created_at)),
    cinema: 'CinePass Central',
    seats: booking.seats,
    pricePerSeat: booking.seats.length ? booking.total_price / booking.seats.length : 0,
    lockExpiresAt: booking.lock_expires_at,
  }
}

const shortId = (id) => String(id || '').slice(-6)

const loadOrder = async () => {
  const queryBookingId = Array.isArray(route.query.bookingId)
    ? route.query.bookingId[0]
    : route.query.bookingId
  const bookingId = order.value.bookingId || queryBookingId
  if (!bookingId) {
    router.replace({ name: 'movies' })
    return
  }

  const bookings = await apiFetch('/bookings/my')
  const booking = bookings.find((item) => item.id === bookingId)
  if (!booking || booking.status !== 'PENDING') {
    router.replace({ name: 'history' })
    return
  }

  if (!order.value.bookingId || !order.value.lockExpiresAt) {
    hydrateOrderFromBooking(booking)
  } else {
    order.value = {
      ...order.value,
      bookingId: booking.id,
      seats: booking.seats,
      pricePerSeat: booking.seats.length ? booking.total_price / booking.seats.length : 0,
      lockExpiresAt: booking.lock_expires_at,
    }
  }

  timeLeft.value = secondsUntil(order.value.lockExpiresAt)
}

const startTimer = () => {
  clearInterval(timerInterval)
  timerInterval = setInterval(() => {
    timeLeft.value = secondsUntil(order.value.lockExpiresAt)
    if (timeLeft.value <= 0) {
      clearInterval(timerInterval)
      error.value = 'Reservation expired. The seats have been released.'
    }
  }, 1000)
}

const handlePay = async () => {
  if (!selected.value || expired.value || !order.value.bookingId) return
  loading.value = true
  error.value = ''
  try {
    await apiFetch('/bookings/confirm', {
      method: 'POST',
      body: { booking_id: order.value.bookingId },
    })
    paid.value = true
    clearInterval(timerInterval)
    setTimeout(() => router.push({ name: 'history' }), 900)
  } catch (err) {
    error.value = err.message
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  try {
    await loadOrder()
    startTimer()
  } catch (err) {
    error.value = err.message
  } finally {
    loadingOrder.value = false
  }
})

onUnmounted(() => clearInterval(timerInterval))
</script>
