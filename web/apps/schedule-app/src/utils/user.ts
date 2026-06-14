export interface AuthUser {
  id: string
  email: string
  role: 'admin' | 'teacher' | 'student'
  is_default_password?: boolean
}

export function getAuthUser(): AuthUser | null {
  try {
    const raw = localStorage.getItem('auth_user')
    return raw ? (JSON.parse(raw) as AuthUser) : null
  } catch {
    return null
  }
}
