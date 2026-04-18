import type { User } from '../store/authStore'

// Читаємо user з cookie
export const getUserFromCookie = (): User | null => {
  const cookies = document.cookie.split(';')
  const authCookie = cookies.find((c) => c.trim().startsWith('auth_user='))
  if (!authCookie) return null

  try {
    const value = authCookie.split('=')[1]
    return JSON.parse(decodeURIComponent(value)) as User
  } catch {
    return null
  }
}

// Видаляємо cookie при логауті
export const clearUserCookie = () => {
  document.cookie = 'auth_user=; path=/; Max-Age=0'
}