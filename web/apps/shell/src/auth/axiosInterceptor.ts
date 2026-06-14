import axios from 'axios'
import { tokenManager } from './tokenManager'

let interceptorSetup = false

// URL які не потребують access_token і не повинні тригерити refresh при 401.
// Shell визначає цей список сам — жоден MFE не повинен про нього знати.
const PUBLIC_URLS = ['/api/auth/login', '/api/auth/refresh', '/api/auth/logout']

const isPublic = (url?: string): boolean =>
  PUBLIC_URLS.some((pub) => url?.includes(pub))

export function setupAxiosInterceptor() {
  if (interceptorSetup) return
  interceptorSetup = true

  // ── Request: додаємо токен лише до захищених запитів ───────────
  axios.interceptors.request.use((config) => {
    if (isPublic(config.url)) return config

    const token = tokenManager.get()
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  })

  // ── Response: 401 → refresh → retry ────────────────────────────
  axios.interceptors.response.use(
    (response) => response,
    async (error) => {
      const original = error.config
      const is401 = error.response?.status === 401
      const alreadyRetried = original._retry

      // Публічні URL не ретраємо — якщо /api/auth/refresh повернув 401,
      // значить refresh_token протух, треба логін, а не нескінченний цикл
      if (is401 && !alreadyRetried && !isPublic(original.url)) {
        original._retry = true
        try {
          const newToken = await tokenManager.refresh()
          original.headers.Authorization = `Bearer ${newToken}`
          return axios(original)
        } catch {
          return Promise.reject(error)
        }
      }

      return Promise.reject(error)
    }
  )
}