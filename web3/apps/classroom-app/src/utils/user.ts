interface AuthUser {
  role: 'admin' | 'teacher' | 'student'
  email: string
}

export function getAuthUser(): AuthUser | null {
  try {
    const raw = localStorage.getItem('auth_user')
    return raw ? (JSON.parse(raw) as AuthUser) : null
  } catch {
    return null
  }
}
