const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api/v1'
const SESSION_KEY = 'cinema_session'

export function getSession() {
  try {
    return JSON.parse(localStorage.getItem(SESSION_KEY) || 'null')
  } catch {
    return null
  }
}

export function saveSession(session) {
  localStorage.setItem(SESSION_KEY, JSON.stringify(session))
}

export function clearSession() {
  localStorage.removeItem(SESSION_KEY)
}

export async function apiFetch(path, options = {}) {
  const session = getSession()
  const headers = new Headers(options.headers || {})

  if (!(options.body instanceof FormData)) {
    headers.set('Content-Type', 'application/json')
  }
  if (session?.token) {
    headers.set('Authorization', `Bearer ${session.token}`)
  }

  const response = await fetch(`${API_BASE_URL}${path}`, {
    ...options,
    headers,
    body:
      options.body && !(options.body instanceof FormData)
        ? JSON.stringify(options.body)
        : options.body,
  })

  const text = await response.text()
  const data = text ? JSON.parse(text) : null
  if (!response.ok) {
    throw new Error(data?.error || data?.message || `Request failed with ${response.status}`)
  }
  return data
}

export function websocketUrl() {
  if (import.meta.env.VITE_WS_URL) {
    return import.meta.env.VITE_WS_URL
  }

  if (API_BASE_URL.startsWith('http')) {
    const url = new URL(API_BASE_URL)
    url.protocol = url.protocol === 'https:' ? 'wss:' : 'ws:'
    url.pathname = `${url.pathname.replace(/\/$/, '')}/ws`
    return url.toString()
  }

  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  return `${protocol}//${window.location.host}${API_BASE_URL}/ws`
}
