import type { AuthUser } from '../types/dashboard'

// Дані користувача кладе shell у localStorage під ключем 'auth_user'.
export function getAuthUser(): AuthUser | null {
  try {
    const raw = localStorage.getItem('auth_user')
    return raw ? (JSON.parse(raw) as AuthUser) : null
  } catch {
    return null
  }
}

const ROLE_LABEL: Record<AuthUser['role'], string> = {
  admin: 'Administrator',
  teacher: 'Teacher',
  student: 'Student',
}

export function roleLabel(role: AuthUser['role'] | undefined): string {
  return role ? ROLE_LABEL[role] : ''
}
