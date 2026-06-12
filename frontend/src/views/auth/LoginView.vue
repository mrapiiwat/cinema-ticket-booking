<template>
  <div class="min-h-screen bg-stone-50 flex">
    <div class="hidden lg:flex lg:w-[44%] bg-white border-r border-stone-100 flex-col justify-between p-14">
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
          <div v-for="feature in features" :key="feature" class="flex items-center gap-3">
            <div class="w-1.5 h-1.5 rounded-full bg-emerald-400"></div>
            <span class="text-sm text-stone-500">{{ feature }}</span>
          </div>
        </div>
      </div>

      <p class="text-xs text-stone-300">© 2026 CinePass</p>
    </div>

    <div class="flex-1 flex items-center justify-center px-8">
      <div class="w-full max-w-sm">
        <div class="flex items-center gap-2 mb-12 lg:hidden">
          <svg class="w-5 h-5 text-stone-700" viewBox="0 0 24 24" fill="none" stroke="currentColor">
            <rect x="2" y="7" width="20" height="13" rx="2" stroke-width="1.5" />
            <path d="M7 7V5a2 2 0 0 1 2-2h6a2 2 0 0 1 2 2v2" stroke-width="1.5" />
          </svg>
          <span class="text-stone-700 font-medium text-sm">CinePass</span>
        </div>

        <h1 class="text-xl font-semibold text-stone-800 mb-1">Welcome back</h1>
        <p class="text-sm text-stone-400 mb-8">Sign in to continue booking</p>

        <div class="space-y-3">
          <button
            @click="loginWithGoogle"
            :disabled="loading || loadingConfig || !googleConfigured"
            class="w-full flex items-center justify-center gap-3 py-3 px-4 bg-stone-900 border border-stone-900 rounded-xl text-sm text-white hover:bg-stone-700 transition-all duration-150 active:scale-[0.99] disabled:opacity-50 disabled:cursor-not-allowed shadow-sm"
          >
            <span
              v-if="loadingRole === 'GOOGLE'"
              class="w-4 h-4 border-[1.5px] border-white/30 border-t-white rounded-full animate-spin"
            ></span>
            <span v-else class="w-4 h-4 bg-white rounded-full text-stone-900 flex items-center justify-center font-semibold text-[11px]">G</span>
            <span>Continue with Google</span>
          </button>

          <div v-if="devLoginEnabled" class="pt-3 border-t border-stone-100 space-y-3">
            <button
              @click="devLogin('USER')"
              :disabled="loading"
              class="w-full flex items-center justify-center gap-3 py-3 px-4 bg-white border border-stone-200 rounded-xl text-sm text-stone-700 hover:bg-stone-50 hover:border-stone-300 transition-all duration-150 active:scale-[0.99] disabled:opacity-50 disabled:cursor-not-allowed shadow-sm"
            >
              <span
                v-if="loadingRole === 'USER'"
                class="w-4 h-4 border-[1.5px] border-stone-200 border-t-stone-500 rounded-full animate-spin"
              ></span>
              <span>Demo User</span>
            </button>

            <button
              @click="devLogin('ADMIN')"
              :disabled="loading"
              class="w-full flex items-center justify-center gap-3 py-3 px-4 bg-white border border-stone-200 rounded-xl text-sm text-stone-700 hover:bg-stone-50 hover:border-stone-300 transition-all duration-150 active:scale-[0.99] disabled:opacity-50 disabled:cursor-not-allowed shadow-sm"
            >
              <span
                v-if="loadingRole === 'ADMIN'"
                class="w-4 h-4 border-[1.5px] border-stone-200 border-t-stone-500 rounded-full animate-spin"
              ></span>
              <span>Demo Admin</span>
            </button>
          </div>
        </div>

        <p v-if="!googleConfigured && !loadingConfig" class="text-center text-xs text-amber-600 mt-5">
          Google OAuth is not configured.
        </p>
        <p v-if="error" class="text-center text-xs text-red-500 mt-5">{{ error }}</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { apiFetch, saveSession } from '../../api/client'

const router = useRouter()
const loadingRole = ref('')
const error = ref('')
const loadingConfig = ref(true)
const googleClientId = ref(import.meta.env.VITE_GOOGLE_CLIENT_ID || '')
const googleScope = ref('openid email profile')
const googleConfigured = ref(false)
const devLoginEnabled = ref(false)
let googleCodeClient = null

const features = ['Real-time seat selection', '5-minute reservation hold', 'Concurrent booking protection']

const loading = computed(() => Boolean(loadingRole.value))

const loadAuthConfig = async () => {
  try {
    const config = await apiFetch('/auth/google/config')
    googleClientId.value = googleClientId.value || config.client_id || ''
    googleScope.value = config.scope || googleScope.value
    googleConfigured.value = Boolean(googleClientId.value && config.configured)
    devLoginEnabled.value = Boolean(config.dev_login_enabled)
  } catch (err) {
    error.value = err.message
  } finally {
    loadingConfig.value = false
  }
}

const loadGoogleScript = () => {
  if (window.google?.accounts?.oauth2) return Promise.resolve()

  return new Promise((resolve, reject) => {
    const existing = document.querySelector('script[src="https://accounts.google.com/gsi/client"]')
    if (existing) {
      existing.addEventListener('load', resolve, { once: true })
      existing.addEventListener('error', reject, { once: true })
      return
    }

    const script = document.createElement('script')
    script.src = 'https://accounts.google.com/gsi/client'
    script.async = true
    script.defer = true
    script.onload = resolve
    script.onerror = () => reject(new Error('Unable to load Google sign-in'))
    document.head.appendChild(script)
  })
}

const initGoogleCodeClient = async () => {
  if (googleCodeClient) return googleCodeClient
  await loadGoogleScript()

  googleCodeClient = window.google.accounts.oauth2.initCodeClient({
    client_id: googleClientId.value,
    scope: googleScope.value,
    ux_mode: 'popup',
    callback: handleGoogleCode,
    error_callback: (googleError) => {
      loadingRole.value = ''
      error.value = googleError?.type || 'Google sign-in was cancelled'
    },
  })
  return googleCodeClient
}

const loginWithGoogle = async () => {
  if (!googleConfigured.value) {
    error.value = 'Google OAuth is not configured.'
    return
  }

  loadingRole.value = 'GOOGLE'
  error.value = ''
  try {
    const client = await initGoogleCodeClient()
    client.requestCode()
  } catch (err) {
    loadingRole.value = ''
    error.value = err.message
  }
}

const handleGoogleCode = async (response) => {
  try {
    if (response.error) throw new Error(response.error)
    if (!response.code) throw new Error('Google did not return an authorization code')

    const session = await apiFetch('/auth/google/code', {
      method: 'POST',
      body: {
        code: response.code,
        redirect_uri: 'postmessage',
      },
    })
    saveAndRedirect(session)
  } catch (err) {
    error.value = err.message
  } finally {
    loadingRole.value = ''
  }
}

const devLogin = async (role) => {
  loadingRole.value = role
  error.value = ''
  try {
    const session = await apiFetch('/auth/dev', {
      method: 'POST',
      body: { role },
    })
    saveAndRedirect(session)
  } catch (err) {
    error.value = err.message
  } finally {
    loadingRole.value = ''
  }
}

const saveAndRedirect = (session) => {
  saveSession({ token: session.token, user: session.user })
  router.push({ name: session.user?.role === 'ADMIN' ? 'admin-dashboard' : 'movies' })
}

onMounted(loadAuthConfig)
</script>
