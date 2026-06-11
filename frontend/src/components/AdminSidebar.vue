<template>
  <aside
    class="w-64 bg-white border-r border-stone-100 flex flex-col justify-between shrink-0 md:flex h-screen sticky top-0"
  >
    <div class="p-6">
      <div class="flex items-center gap-2.5 mb-8">
        <svg class="w-5 h-5 text-stone-800" viewBox="0 0 24 24" fill="none" stroke="currentColor">
          <rect x="2" y="7" width="20" height="13" rx="2" stroke-width="1.5" />
          <path d="M7 7V5a2 2 0 0 1 2-2h6a2 2 0 0 1 2 2v2" stroke-width="1.5" />
          <line x1="12" y1="12" x2="12" y2="16" stroke-width="1.5" stroke-linecap="round" />
          <line x1="10" y1="14" x2="14" y2="14" stroke-width="1.5" stroke-linecap="round" />
        </svg>
        <span class="text-stone-800 font-semibold text-sm tracking-wide">CinePass Admin</span>
      </div>

      <nav class="space-y-1">
        <button
          @click="$router.push({ name: 'admin-dashboard' })"
          :class="[
            'w-full flex items-center gap-3 px-4 py-2.5 text-sm font-medium rounded-xl transition-colors',
            $route.name === 'admin-dashboard'
              ? 'bg-stone-900 text-white'
              : 'text-stone-500 hover:bg-stone-50 hover:text-stone-800',
          ]"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="1.5"
              d="M4 6a2 2 0 012-2h2a2 2 0 012 2v4a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v4a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v4a2 2 0 01-2 2H6a2 2 0 01-2-2v-4zM14 16a2 2 0 012-2h2a2 2 0 012 2v4a2 2 0 01-2 2h-2a2 2 0 01-2-2v-4z"
            />
          </svg>
          Dashboard
        </button>

        <button
          @click="$router.push({ name: 'admin-logs' })"
          :class="[
            'w-full flex items-center gap-3 px-4 py-2.5 text-sm font-medium rounded-xl transition-colors',
            $route.name === 'admin-logs'
              ? 'bg-stone-900 text-white'
              : 'text-stone-500 hover:bg-stone-50 hover:text-stone-800',
          ]"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="1.5"
              d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
            />
          </svg>
          Audit Logs
        </button>
      </nav>
    </div>

    <div>
      <div class="px-6 pb-4">
        <button
          @click="handleLogout"
          class="w-full flex items-center gap-3 px-4 py-2.5 text-sm font-medium rounded-xl text-red-600 hover:bg-red-50 transition-colors"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
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

      <div class="p-4 border-t border-stone-100 flex items-center gap-3">
        <div
          class="w-8 h-8 rounded-full bg-stone-100 flex items-center justify-center text-stone-600 font-semibold text-xs border border-stone-200"
        >
          A
        </div>
        <div>
          <p class="text-xs font-medium text-stone-800">Admin Staff</p>
          <p class="text-[10px] text-stone-400">System Manager</p>
        </div>
      </div>
    </div>
  </aside>
</template>

<script setup>
import { useRouter } from 'vue-router'

const router = useRouter()

const handleLogout = () => {
  // TODO: เพิ่ม Logic การลบ Token หรือเคลียร์ Session ของ Admin ตรงนี้ในภายหลัง

  // เตะกลับไปหน้า Login
  router.push({ name: 'login' })
}
</script>
