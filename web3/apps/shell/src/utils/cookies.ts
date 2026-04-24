import type { User } from '../store/authStore'

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

export const clearUserCookie = () => {
  document.cookie = 'auth_user=; path=/; Max-Age=0'
  document.cookie = 'access_token=; path=/; Max-Age=0'
  document.cookie = 'refresh_token=; path=/; Max-Age=0'
  localStorage.removeItem('access_token')
  localStorage.removeItem('refresh_token')
}