<template>
  <div class="min-h-screen bg-stone-50 flex">
    <div
      class="hidden lg:flex lg:w-[44%] bg-white border-r border-stone-100 flex-col justify-between p-14"
    >
      <div class="flex items-center gap-2.5">
        <svg class="w-5 h-5 text-stone-800" viewBox="0 0 24 24" fill="none" stroke="currentColor">
          <rect x="2" y="7" width="20" height="13" rx="2" stroke-width="1.5" />
          <path d="M7 7V5a2 2 0 0 1 2-2h6a2 2 0 0 1 2 2v2" stroke-width="1.5" />
          <line x1="12" y1="12" x2="12" y2="16" stroke-width="1.5" stroke-linecap="round" />
          <line x1="10" y1="14" x2="14" y2="14" stroke-width="1.5" stroke-linecap="round" />
        </svg>
        <span class="text-stone-800 font-medium text-sm tracking-wide">CinePass</span>
      </div>

      <div>
        <p class="text-3xl font-light text-stone-800 leading-snug tracking-tight">
          Great stories<br />deserve great seats.
        </p>
        <div class="mt-10 space-y-4">
          <div v-for="f in features" :key="f.text" class="flex items-center gap-3">
            <div class="w-1 h-1 rounded-full bg-stone-300"></div>
            <span class="text-sm text-stone-400">{{ f.text }}</span>
          </div>
        </div>
      </div>

      <p class="text-xs text-stone-300">© 2026 CinePass</p>
    </div>

    <div class="flex-1 flex items-center justify-center px-8">
      <div class="w-full max-w-85">
        <div class="flex items-center gap-2 mb-12 lg:hidden">
          <svg class="w-5 h-5 text-stone-700" viewBox="0 0 24 24" fill="none" stroke="currentColor">
            <rect x="2" y="7" width="20" height="13" rx="2" stroke-width="1.5" />
            <path d="M7 7V5a2 2 0 0 1 2-2h6a2 2 0 0 1 2 2v2" stroke-width="1.5" />
          </svg>
          <span class="text-stone-700 font-medium text-sm">CinePass</span>
        </div>

        <h1 class="text-xl font-semibold text-stone-800 mb-1">Welcome back</h1>
        <p class="text-sm text-stone-400 mb-8">Sign in to continue booking</p>

        <button
          @click="loginWithGoogle"
          :disabled="loading"
          class="w-full flex items-center justify-center gap-3 py-3 px-4 bg-white border border-stone-200 rounded-xl text-sm text-stone-700 hover:bg-stone-50 hover:border-stone-300 transition-all duration-150 active:scale-[0.99] disabled:opacity-50 disabled:cursor-not-allowed shadow-sm"
        >
          <span
            v-if="loading"
            class="w-4 h-4 border-[1.5px] border-stone-200 border-t-stone-500 rounded-full animate-spin"
          ></span>
          <svg v-else width="18" height="18" viewBox="0 0 24 24" class="shrink-0">
            <path
              fill="#4285F4"
              d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
            />
            <path
              fill="#34A853"
              d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
            />
            <path
              fill="#FBBC05"
              d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l3.66-2.84z"
            />
            <path
              fill="#EA4335"
              d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
            />
          </svg>
          <span>{{ loading ? 'Signing in...' : 'Continue with Google' }}</span>
        </button>

        <p class="text-center text-[11px] text-stone-300 mt-8 leading-relaxed">
          By continuing you agree to our
          <a
            href="#"
            class="text-stone-400 underline underline-offset-2 hover:text-stone-600 transition-colors"
            >Terms</a
          >
          and
          <a
            href="#"
            class="text-stone-400 underline underline-offset-2 hover:text-stone-600 transition-colors"
            >Privacy Policy</a
          >
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const loading = ref(false)

const features = [
  { text: 'Real-time seat selection' },
  { text: '5-minute seat reservation hold' },
  { text: 'Concurrent booking protection' },
]

const loginWithGoogle = async () => {
  loading.value = true
  try {
    console.log('Login with Google')
  } catch (err) {
    console.error(err)
  } finally {
    loading.value = false
  }
}
</script>
