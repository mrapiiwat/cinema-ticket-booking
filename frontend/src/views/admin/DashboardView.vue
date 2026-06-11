<template>
  <div class="min-h-screen bg-stone-50 flex">
    <div class="flex-1 flex flex-col min-w-0">
      <main class="flex-1 p-8 overflow-y-auto max-w-7xl w-full mx-auto">
        <div class="grid grid-cols-1 sm:grid-cols-3 gap-5 mb-8">
          <div class="bg-white border border-stone-100 rounded-2xl p-5 shadow-sm">
            <p class="text-xs text-stone-400 uppercase tracking-wider mb-1">Total Revenue</p>
            <p class="text-2xl font-semibold text-stone-800">
              ฿{{ stats.revenue.toLocaleString() }}
            </p>
          </div>
          <div class="bg-white border border-stone-100 rounded-2xl p-5 shadow-sm">
            <p class="text-xs text-stone-400 uppercase tracking-wider mb-1">Confirmed Bookings</p>
            <p class="text-2xl font-semibold text-stone-800">
              {{ stats.confirmedCount }}
              <span class="text-xs font-normal text-stone-400">orders</span>
            </p>
          </div>
          <div class="bg-white border border-stone-100 rounded-2xl p-5 shadow-sm">
            <p class="text-xs text-amber-600 uppercase tracking-wider mb-1">Active Redis Locks</p>
            <p class="text-2xl font-semibold text-amber-700">
              {{ stats.activeLocks }}
              <span class="text-xs font-normal text-amber-500">holding</span>
            </p>
          </div>
        </div>

        <div
          class="bg-white border border-stone-100 rounded-2xl p-5 mb-6 shadow-sm flex flex-col sm:flex-row gap-4 items-center justify-between"
        >
          <div class="flex flex-col sm:flex-row gap-3 w-full sm:w-auto flex-1">
            <div class="relative flex-1 max-w-xs">
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
                placeholder="Search booking ID / user..."
                class="w-full pl-9 pr-4 py-2 text-sm bg-stone-50 border border-stone-200 rounded-xl text-stone-700 placeholder-stone-400 focus:outline-none focus:border-stone-400 focus:bg-white transition-all"
              />
            </div>

            <select
              v-model="filters.movie"
              class="px-3 py-2 text-sm bg-stone-50 border border-stone-200 rounded-xl text-stone-600 focus:outline-none focus:border-stone-400 focus:bg-white transition-all"
            >
              <option value="">All Movies</option>
              <option v-for="movie in uniqueMovies" :key="movie" :value="movie">{{ movie }}</option>
            </select>

            <select
              v-model="filters.status"
              class="px-3 py-2 text-sm bg-stone-50 border border-stone-200 rounded-xl text-stone-600 focus:outline-none focus:border-stone-400 focus:bg-white transition-all"
            >
              <option value="">All Status</option>
              <option value="confirmed">Confirmed</option>
              <option value="locked">Locked (Pending)</option>
              <option value="timeout">Timeout</option>
            </select>
          </div>

          <p class="text-xs text-stone-400 shrink-0">Found {{ filteredBookings.length }} records</p>
        </div>

        <div class="bg-white border border-stone-100 rounded-2xl shadow-sm overflow-hidden">
          <div class="overflow-x-auto">
            <table class="w-full text-left border-collapse">
              <thead>
                <tr
                  class="border-b border-stone-100 bg-stone-50/50 text-stone-400 text-[11px] uppercase tracking-wider font-semibold"
                >
                  <th class="py-4 px-6">Booking ID</th>
                  <th class="py-4 px-6">Movie & Showtie</th>
                  <th class="py-4 px-6">User</th>
                  <th class="py-4 px-6 text-center">Seats</th>
                  <th class="py-4 px-6 text-right">Total Price</th>
                  <th class="py-4 px-6 text-center">Status</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-stone-50 text-sm text-stone-700">
                <tr
                  v-for="booking in filteredBookings"
                  :key="booking.id"
                  class="hover:bg-stone-50/30 transition-colors"
                >
                  <td class="py-4 px-6 font-mono text-xs text-stone-400">
                    #{{ booking.bookingId }}
                  </td>

                  <td class="py-4 px-6">
                    <p class="font-medium text-stone-800">{{ booking.movie }}</p>
                    <p class="text-xs text-stone-400 mt-0.5">{{ booking.showtime }}</p>
                  </td>

                  <td class="py-4 px-6">
                    <p class="text-stone-800">{{ booking.user.name }}</p>
                    <p class="text-[11px] text-stone-400 font-mono">{{ booking.user.email }}</p>
                  </td>

                  <td class="py-4 px-6 text-center">
                    <div class="flex gap-1 justify-center flex-wrap">
                      <span
                        v-for="seat in booking.seats"
                        :key="seat"
                        class="px-1.5 py-0.5 bg-stone-100 text-stone-600 text-[10px] rounded font-medium"
                      >
                        {{ seat }}
                      </span>
                    </div>
                  </td>

                  <td class="py-4 px-6 text-right font-medium text-stone-800">
                    ฿{{ booking.totalPrice.toLocaleString() }}
                  </td>

                  <td class="py-4 px-6 text-center">
                    <span
                      :class="[
                        'text-[10px] uppercase tracking-wider px-2 py-0.5 rounded-full font-semibold inline-block',
                        booking.status === 'confirmed'
                          ? 'bg-green-50 text-green-600'
                          : booking.status === 'locked'
                            ? 'bg-amber-50 text-amber-600 animate-pulse'
                            : 'bg-stone-100 text-stone-400 line-through',
                      ]"
                    >
                      {{ booking.status }}
                    </span>
                  </td>
                </tr>

                <tr v-if="!filteredBookings.length">
                  <td colspan="6" class="text-center py-12 text-stone-400 text-sm">
                    No booking records found matching the filters.
                  </td>
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
import { ref, computed } from 'vue'

// Controls Filters
const filters = ref({
  search: '',
  movie: '',
  status: '',
})

// Mock Bookings Data (จำลองโครงสร้างข้อมูลที่จะมาจาก MongoDB + Redis Lock)
const bookings = ref([
  {
    id: 1,
    bookingId: 'BK-0001',
    movie: 'Dune: Part Two',
    showtime: 'Today 14:30 · Screen 4',
    user: { name: 'John Doe', email: 'john@example.com' },
    seats: ['E5', 'E6'],
    totalPrice: 470, // (220*2) + 30 service fee
    status: 'confirmed',
  },
  {
    id: 2,
    bookingId: 'BK-0002',
    movie: 'Oppenheimer',
    showtime: 'Today 15:00 · Screen 1',
    user: { name: 'Alice Smith', email: 'alice@gmail.com' },
    seats: ['A1', 'A2'],
    totalPrice: 470,
    status: 'locked', // กำลังติด Redis Lock ค้างไว้ 5 นาทีรอจ่ายเงิน
  },
  {
    id: 3,
    bookingId: 'BK-0003',
    movie: 'Dune: Part Two',
    showtime: 'Today 14:30 · Screen 4',
    user: { name: 'Bob Johnson', email: 'bob.j@outlook.com' },
    seats: ['F1'],
    totalPrice: 250,
    status: 'timeout', // ชำระเงินไม่ทันภายใน 5 นาทีระบบดีดหลุด
  },
  {
    id: 4,
    bookingId: 'BK-0004',
    movie: 'Past Lives',
    showtime: 'Tomorrow 19:00 · Screen 2',
    user: { name: 'Sarah Connor', email: 'sarah@sky.net' },
    seats: ['B3', 'B4', 'B5'],
    totalPrice: 690,
    status: 'confirmed',
  },
])

// ดึงรายชื่อหนังทั้งหมดแบบไม่ซ้ำกันมาทำเป็น Option ในตัวเลือกกรอง
const uniqueMovies = computed(() => {
  return [...new Set(bookings.value.map((b) => b.movie))]
})

// คำนวณสรุปสถิติ (Stats Overview)
const stats = computed(() => {
  let revenue = 0
  let confirmedCount = 0
  let activeLocks = 0

  bookings.value.forEach((b) => {
    if (b.status === 'confirmed') {
      revenue += b.totalPrice
      confirmedCount++
    } else if (b.status === 'locked') {
      activeLocks++
    }
  })

  return { revenue, confirmedCount, activeLocks }
})

// ค้นหาและคัดกรองข้อมูลตามฟิลเตอร์ (เงื่อนไขบังคับของโจทย์)
const filteredBookings = computed(() => {
  return bookings.value.filter((booking) => {
    const matchesSearch =
      booking.bookingId.toLowerCase().includes(filters.value.search.toLowerCase()) ||
      booking.user.name.toLowerCase().includes(filters.value.search.toLowerCase()) ||
      booking.user.email.toLowerCase().includes(filters.value.search.toLowerCase())

    const matchesMovie = !filters.value.movie || booking.movie === filters.value.movie
    const matchesStatus = !filters.value.status || booking.status === filters.value.status

    return matchesSearch && matchesMovie && matchesStatus
  })
})
</script>
