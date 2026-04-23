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

const getTokenFromCookie = (name: string): string | null => {
  const cookie = document.cookie.split(';').find(c => c.trim().startsWith(`${name}=`))
  return cookie ? cookie.split('=')[1]?.trim() ?? null : null
}

const setTokens = (access: string, refresh: string) => {
  localStorage.setItem('access_token', access)
  localStorage.setItem('refresh_token', refresh)
  
  // Коротке життя для access token cookie — 15 хв
  document.cookie = `access_token=${access}; path=/; Max-Age=900; SameSite=Strict`
  // Refresh token — 30 днів
  document.cookie = `refresh_token=${refresh}; path=/; Max-Age=${86400 * 30}; SameSite=Strict`
  
  // Зберігаємо auth_user щоб initAuth не скидав стан після рефрешу
  const existingUserCookie = document.cookie
    .split(';')
    .find(c => c.trim().startsWith('auth_user='))
  if (existingUserCookie) {
    const value = existingUserCookie.split('=').slice(1).join('=').trim()
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

  api.interceptors.request.use((config) => {
    const token = localStorage.getItem('access_token') ?? getTokenFromCookie('access_token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  })

  api.interceptors.response.use(
    (response) => response,
    async (error) => {
      const originalRequest = error.config

      if (error.response?.status !== 401 || originalRequest._retry) {
        return Promise.reject(error)
      }

      originalRequest._retry = true

      if (isRefreshing) {
        // Чекаємо поки інший запит зрефрешить токен
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject })
        }).then(token => {
          originalRequest.headers.Authorization = `Bearer ${token}`
          return api(originalRequest)
        })
      }

      isRefreshing = true

      try {
        const refreshToken = localStorage.getItem('refresh_token') ?? getTokenFromCookie('refresh_token')
        
        if (!refreshToken) {
          throw new Error('No refresh token')
        }

        const res = await axios.post(
          '/api/auth/refresh',
          {},
          {
            headers: {
              Authorization: `Bearer ${refreshToken}`
            }
          }
        )
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