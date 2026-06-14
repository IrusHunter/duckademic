import type { User } from '../store/authStore'

export const getUserFromCookie = (): User | null => {
  try {
    const raw = localStorage.getItem('auth_user')
    if (!raw) return null
    return JSON.parse(raw) as User
  } catch {
    return null
  }
}

// utils/jwt.ts
export const isTokenExpired = (token: string | null, bufferSec = 30): boolean => {
  if (!token) return true
  try {
    const payload = JSON.parse(atob(token.split('.')[1]))
    // bufferSec — оновлюємо за 30с до закінчення
    return payload.exp * 1000 < Date.now() + bufferSec * 1000
  } catch {
    return true
  }
}

export const setAuthCookies = (user: User, refreshToken?: string): void => {
  if (refreshToken) {
    localStorage.setItem('refresh_token', refreshToken)
  }
  localStorage.setItem('auth_user', JSON.stringify({
    id: user.id,
    email: user.email,
    role: user.role,
    is_default_password: user.is_default_password,
  }))
}

export const clearUserCookie = (): void => {
  localStorage.removeItem('refresh_token')
  localStorage.removeItem('auth_user')
}