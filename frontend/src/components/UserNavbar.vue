<template>
  <header
    class="bg-white border-b border-stone-100 px-6 py-4 flex items-center justify-between relative z-50"
  >
    <div
      @click="$router.push({ path: '/user' })"
      class="flex items-center gap-2.5 cursor-pointer group"
    >
      <svg
        class="w-5 h-5 text-stone-800 transition-transform group-hover:scale-105"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
      >
        <rect x="2" y="7" width="20" height="13" rx="2" stroke-width="1.5" />
        <path d="M7 7V5a2 2 0 0 1 2-2h6a2 2 0 0 1 2 2v2" stroke-width="1.5" />
        <line x1="12" y1="12" x2="12" y2="16" stroke-width="1.5" stroke-linecap="round" />
        <line x1="10" y1="14" x2="14" y2="14" stroke-width="1.5" stroke-linecap="round" />
      </svg>
      <span class="text-stone-800 font-medium text-sm tracking-wide">CinePass</span>
    </div>

    <div class="relative">
      <div v-if="isOpen" @click="isOpen = false" class="fixed inset-0 z-40"></div>

      <button
        @click="isOpen = !isOpen"
        class="relative z-50 flex items-center gap-2.5 p-1 -mr-1 rounded-lg hover:bg-stone-50 transition-colors outline-none"
      >
        <div
          class="w-7 h-7 rounded-full bg-stone-100 flex items-center justify-center overflow-hidden border border-stone-200 shrink-0"
        >
          <img v-if="user.avatar" :src="user.avatar" class="w-full h-full object-cover" />
          <svg
            v-else
            class="w-4 h-4 text-stone-400"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="1.5"
              d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
            />
          </svg>
        </div>
        <span class="text-sm text-stone-600 font-medium hidden sm:block">{{ user.name }}</span>

        <svg
          class="w-3 h-3 text-stone-400 hidden sm:block transition-transform duration-200"
          :class="isOpen ? 'rotate-180' : ''"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M19 9l-7 7-7-7"
          />
        </svg>
      </button>

      <transition
        enter-active-class="transition ease-out duration-100"
        enter-from-class="transform opacity-0 scale-95"
        enter-to-class="transform opacity-100 scale-100"
        leave-active-class="transition ease-in duration-75"
        leave-from-class="transform opacity-100 scale-100"
        leave-to-class="transform opacity-0 scale-95"
      >
        <div
          v-if="isOpen"
          class="absolute right-0 mt-2 w-48 bg-white border border-stone-100 rounded-xl shadow-[0_4px_20px_-4px_rgba(0,0,0,0.05)] z-50 py-1.5 overflow-hidden"
        >
          <div class="px-4 py-2 border-b border-stone-50 sm:hidden mb-1">
            <p class="text-sm font-medium text-stone-800">{{ user.name }}</p>
          </div>

          <button
            @click="goToHistory"
            class="w-full text-left px-4 py-2.5 text-sm text-stone-600 hover:bg-stone-50 hover:text-stone-800 flex items-center gap-2.5 transition-colors"
          >
            <svg
              class="w-4 h-4 text-stone-400"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="1.5"
                d="M15 5v2m0 4v2m0 4v2M5 5a2 2 0 00-2 2v3a2 2 0 110 4v3a2 2 0 002 2h14a2 2 0 002-2v-3a2 2 0 110-4V7a2 2 0 00-2-2H5z"
              />
            </svg>
            My Tickets
          </button>

          <div class="h-px bg-stone-50 my-1 mx-3"></div>

          <button
            @click="handleLogout"
            class="w-full text-left px-4 py-2.5 text-sm text-red-600 hover:bg-red-50 flex items-center gap-2.5 transition-colors"
          >
            <svg class="w-4 h-4 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="1.5"
                d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"
              />
            </svg>
            Sign out
          </button>
        </div>
      </transition>
    </div>
  </header>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const user = ref({
  name: 'John Doe',
  avatar: '',
})

const isOpen = ref(false)

const goToHistory = () => {
  isOpen.value = false // ปิดเมนูก่อนเปลี่ยนหน้า
  router.push({ name: 'history' }) // สมมติว่าหน้าตั๋วตั้งชื่อ route ว่า 'history'
}

const handleLogout = () => {
  isOpen.value = false
  // TODO: ใส่ Logic เคลียร์ Token หรือ Session ตรงนี้
  console.log('Logging out...')
}
</script>
