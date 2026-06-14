import axios from 'axios'
import { clearUserCookie } from '../utils/cookies'
import { useAuthStore } from '../store/authStore'

// __refreshPromise живе у window — єдиний на всю вкладку, спільний між
// усіма MFE-бандлами. Promise — жива JS-сутність, її не можна покласти в
// localStorage, тому саме window гарантує, що всі MFE бачать один і той
// самий «refresh у процесі» і не шлють паралельних запитів.
// access_token зберігається в localStorage — він теж спільний для всього
// origin і доступний будь-якому MFE без додаткової координації.
declare global {
  interface Window {
    __refreshPromise?: Promise<string> | null
  }
}

const getAccessToken = (): string | null => 
  localStorage.getItem('access_token') ?? null
const setAccessToken = (token: string | null): void => {
  if (token) {
    localStorage.setItem('access_token', token)
  } else {
    localStorage.removeItem('access_token')
  }
}
const getRefreshPromise = (): Promise<string> | null => window.__refreshPromise ?? null
const setRefreshPromise = (p: Promise<string> | null): void => { window.__refreshPromise = p }

// tokenManager.ts
const refreshAxios = axios.create()

const doRefreshRequest = (): Promise<string> => {
  const existing = getRefreshPromise()
  if (existing) return existing

  const refreshToken = localStorage.getItem('refresh_token') ?? ''
  if (!refreshToken) {
    return Promise.reject(new Error('No refresh token'))
  }
  const accessToken = getAccessToken() ?? ''

  const promise = refreshAxios  // ← не axios, а refreshAxios
    .post<{ access_token: string; refresh_token: string }>('/api/auth/refresh', {
      ...(accessToken ? { access_token: accessToken } : {}),
      refresh_token: refreshToken,
    }, {
      withCredentials: true,
      headers: { 'Content-Type': 'application/json' },
    })
    .then(({ data }) => {
      setAccessToken(data.access_token)
      if (data.refresh_token) {
        localStorage.setItem('refresh_token', data.refresh_token)
      }
      return data.access_token
    })
    .finally(() => {
      setRefreshPromise(null)
    })

  setRefreshPromise(promise)
  return promise
}

export const tokenManager = {
  get: (): string | null => getAccessToken(),
  set: (token: string): void => setAccessToken(token),
  clear: (): void => setAccessToken(null),

  initRefresh: async (): Promise<string | null> => {
    try {
      return await doRefreshRequest()
    } catch {
      setAccessToken(null)
      return null
    }
  },

  refresh: (): Promise<string> => {
    return doRefreshRequest().catch((err) => {
      setAccessToken(null)
      const isServerRejection =
        err.response?.status === 401 || err.response?.status === 403
      if (isServerRejection) {
        window.dispatchEvent(new CustomEvent('auth:logout'))
      }
      return Promise.reject(err)
    })
  },

  logout: async (): Promise<void> => {
    try {
      await axios.post('/api/auth/logout', null, {
        withCredentials: true,
        headers: { 'Content-Type': 'application/json' },
      })
    } catch {
      // сервер недоступний — очищаємо локальний стан
    } finally {
      setAccessToken(null)
      clearUserCookie()
      useAuthStore.getState().clearUser()
    }
  },
}