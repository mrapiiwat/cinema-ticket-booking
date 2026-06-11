<template>
  <div class="min-h-screen bg-stone-50">
    <main class="max-w-lg mx-auto px-6 py-10">
      <div
        class="flex items-center gap-2 mb-6 px-4 py-3 bg-amber-50 border border-amber-100 rounded-xl"
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
        <p class="text-xs text-amber-600 flex-1">
          Complete payment within <span class="font-semibold">{{ formattedTimer }}</span> or your
          seats will be released
        </p>
      </div>

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
                >{{ seat }}</span
              >
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
              selected === method.id
                ? 'border-stone-800 bg-stone-50'
                : 'border-stone-100 hover:border-stone-200',
            ]"
          >
            <input type="radio" :value="method.id" v-model="selected" class="accent-stone-800" />
            <div class="flex items-center gap-2.5 flex-1">
              <div
                class="w-8 h-8 rounded-lg bg-stone-50 border border-stone-100 flex items-center justify-center shrink-0"
              >
                <component :is="method.icon" class="w-4 h-4 text-stone-500" />
              </div>
              <div>
                <p class="text-sm font-medium text-stone-700">{{ method.label }}</p>
                <p class="text-xs text-stone-400">{{ method.desc }}</p>
              </div>
            </div>
          </label>
        </div>
      </div>

      <button
        @click="handlePay"
        :disabled="!selected || loading || paid"
        class="w-full py-3.5 rounded-xl text-sm font-medium transition-all duration-150 active:scale-[0.99] disabled:opacity-40 disabled:cursor-not-allowed bg-stone-800 text-white hover:bg-stone-700"
      >
        <span
          v-if="loading"
          class="inline-block w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin mr-2 align-middle"
        ></span>
        <span v-if="paid">✓ Payment Confirmed</span>
        <span v-else-if="loading">Processing...</span>
        <span v-else>Pay ฿{{ grandTotal }}</span>
      </button>

      <p class="text-center text-[11px] text-stone-300 mt-4">
        Secured payment · No real transaction will be made
      </p>
    </main>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, defineComponent, h } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

// Timer — inherit remaining time from booking (mock: 5 min)
const timeLeft = ref(history.state?.timeLeft ?? 300)
let timerInterval = null

onMounted(() => {
  timerInterval = setInterval(() => {
    timeLeft.value--
    if (timeLeft.value <= 0) {
      clearInterval(timerInterval)
      router.replace({ name: 'movies' }) // lock expired → back to movies
    }
  }, 1000)
})
onUnmounted(() => clearInterval(timerInterval))

const formattedTimer = computed(() => {
  const m = Math.floor(timeLeft.value / 60)
  const s = timeLeft.value % 60
  return `${m}:${s.toString().padStart(2, '0')}`
})

// Order — รับข้อมูลจาก BookingView ผ่าน route state
const state = history.state
const order = ref({
  bookingId: state?.bookingId ?? '',
  movie: state?.movie ?? 'Unknown Movie',
  time: state?.time ?? '',
  cinema: state?.cinema ?? '',
  seats: state?.seats ?? [],
  pricePerSeat: state?.pricePerSeat ?? 0,
})

const serviceFee = 30
const grandTotal = computed(() => order.value.seats.length * order.value.pricePerSeat + serviceFee)

// Payment methods
const CreditIcon = defineComponent({
  render: () =>
    h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
      h('rect', { x: '2', y: '5', width: '20', height: '14', rx: '2', 'stroke-width': '1.5' }),
      h('path', { d: 'M2 10h20', 'stroke-width': '1.5' }),
    ]),
})
const QRIcon = defineComponent({
  render: () =>
    h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
      h('rect', { x: '3', y: '3', width: '7', height: '7', rx: '1', 'stroke-width': '1.5' }),
      h('rect', { x: '14', y: '3', width: '7', height: '7', rx: '1', 'stroke-width': '1.5' }),
      h('rect', { x: '3', y: '14', width: '7', height: '7', rx: '1', 'stroke-width': '1.5' }),
      h('path', {
        d: 'M14 14h3v3h-3zM17 17h3v3h-3zM14 17h0M17 14h0',
        'stroke-width': '1.5',
        'stroke-linecap': 'round',
      }),
    ]),
})
const WalletIcon = defineComponent({
  render: () =>
    h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
      h('path', {
        d: 'M20 7H4a2 2 0 0 0-2 2v9a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2V9a2 2 0 0 0-2-2z',
        'stroke-width': '1.5',
      }),
      h('path', {
        d: 'M16 13a1 1 0 1 0 2 0 1 1 0 0 0-2 0z',
        fill: 'currentColor',
        'stroke-width': '0',
      }),
      h('path', { d: 'M4 7V5a2 2 0 0 1 2-2h12', 'stroke-width': '1.5' }),
    ]),
})

const paymentMethods = [
  { id: 'card', label: 'Credit / Debit Card', desc: 'Visa, Mastercard, JCB', icon: CreditIcon },
  { id: 'promptpay', label: 'PromptPay', desc: 'Scan QR to pay', icon: QRIcon },
  { id: 'wallet', label: 'TrueMoney Wallet', desc: 'Pay with e-wallet', icon: WalletIcon },
]

const selected = ref('card')
const loading = ref(false)
const paid = ref(false)

const handlePay = async () => {
  if (!selected.value) return
  loading.value = true
  try {
    // TODO: POST /api/payment { bookingId, method: selected }
    // Backend confirms lock → sets seat status to BOOKED
    // Then publish event to Message Queue (Booking Success)
    await new Promise((r) => setTimeout(r, 1200)) // mock
    paid.value = true
    clearInterval(timerInterval)
    setTimeout(() => router.push({ name: 'movies' }), 1500)
  } catch (err) {
    console.error('Payment failed:', err)
  } finally {
    loading.value = false
  }
}
</script>
