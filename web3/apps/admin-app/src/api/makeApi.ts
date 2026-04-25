import axios from 'axios'
import type { ServiceDef } from '../types/admin'
import { SERVICES } from '../config/services'

let isRefreshing = false
let failedQueue: Array<{
  resolve: (token: string) => void
  reject: (err: unknown) => void
}> = []

const processQueue = (error: unknown, token: string | null) => {
  failedQueue.forEach(p => error ? p.reject(error) : p.resolve(token!))
  failedQueue = []
}

const getToken = (name: string): string | null => {
  // Спочатку localStorage, потім cookie як fallback
  const fromStorage = localStorage.getItem(name)
  if (fromStorage) return fromStorage
  const cookie = document.cookie.split(';').find(c => c.trim().startsWith(`${name}=`))
  return cookie ? cookie.split('=')[1]?.trim() ?? null : null
}

const setTokens = (access: string, refresh: string) => {
  localStorage.setItem('access_token', access)
  localStorage.setItem('refresh_token', refresh)
  document.cookie = `access_token=${access}; path=/; Max-Age=890; SameSite=Strict`
  document.cookie = `refresh_token=${refresh}; path=/; Max-Age=${86400 * 30}; SameSite=Strict`
  // Оновлюємо auth_user куку щоб initAuth не скидав стан
  const existingUser = document.cookie.split(';').find(c => c.trim().startsWith('auth_user='))
  if (existingUser) {
    const value = existingUser.split('=').slice(1).join('=').trim()
    document.cookie = `auth_user=${value}; path=/; Max-Age=86400; SameSite=Strict`
  }
}

const clearTokens = () => {
  localStorage.removeItem('access_token')
  localStorage.removeItem('refresh_token')
  document.cookie = 'access_token=; path=/; Max-Age=0'
  document.cookie = 'refresh_token=; path=/; Max-Age=0'
  document.cookie = 'auth_user=; path=/; Max-Age=0'
}

export function makeApi(baseURL: string) {
  const api = axios.create({ baseURL })

  // Додаємо access_token до кожного запиту
  api.interceptors.request.use((config) => {
    const token = getToken('access_token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  })

  // При 401 — оновлюємо токени через /refresh і повторюємо запит
  api.interceptors.response.use(
    (response) => response,
    async (error) => {
      const originalRequest = error.config

      if (error.response?.status !== 401 || originalRequest._retry) {
        return Promise.reject(error)
      }

      originalRequest._retry = true

      // Якщо вже рефрешимо — ставимо в чергу
      if (isRefreshing) {
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject })
        }).then(token => {
          originalRequest.headers.Authorization = `Bearer ${token}`
          return api(originalRequest)
        })
      }

      isRefreshing = true

      try {
        const accessToken = getToken('access_token')
        const refreshToken = getToken('refresh_token')

        if (!refreshToken) throw new Error('No refresh token')

        // Згідно з документацією: POST /refresh з body { access_token, refresh_token }
        const res = await axios.post('/api/auth/refresh', {
          access_token: accessToken ?? '',
          refresh_token: refreshToken,
        })

        const { access_token, refresh_token } = res.data
        setTokens(access_token, refresh_token)
        processQueue(null, access_token)

        originalRequest.headers.Authorization = `Bearer ${access_token}`
        return api(originalRequest)
      } catch (refreshError) {
        processQueue(refreshError, null)
        clearTokens()
        window.location.href = '/login'
        return Promise.reject(refreshError)
      } finally {
        isRefreshing = false
      }
    }
  )

  return api
}

export function getServiceByKey(key: string): ServiceDef | undefined {
  return SERVICES.find(s => s.key === key)
}

export function normalizeArray(result: unknown): Record<string, unknown>[] {
  if (Array.isArray(result)) return result
  if (result && typeof result === 'object') {
    const obj = result as Record<string, unknown>
    if (Array.isArray(obj.data)) return obj.data as Record<string, unknown>[]
    if (Array.isArray(obj.items)) return obj.items as Record<string, unknown>[]
    if (Array.isArray(obj.results)) return obj.results as Record<string, unknown>[]
    return [obj]
  }
  return []
}