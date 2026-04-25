import { useState } from 'react'
import axios from 'axios'

type User = {
  id: string
  login: string
  role: 'admin' | 'student' | 'teacher'
  is_default_password: boolean
}

type Props = {
  onLoginSuccess?: (data: User & { access_token: string; refresh_token: string }) => void
}

const mapRole = (role: string): 'admin' | 'student' | 'teacher' => {
  const roleMap: Record<string, 'admin' | 'student' | 'teacher'> = {
    admin: 'admin',
    student: 'student',
    teacher: 'teacher',
  }
  return roleMap[role.toLowerCase()] ?? 'student'
}

// authApp — незалежний MFE. Він нічого не знає про shell:
// жодних спільних стор, жодних флагів (skipAuthInterceptor тощо).
// Єдина відповідальність: зібрати credentials → POST /api/auth/login → передати результат shell через onLoginSuccess.
// Shell сам вирішує що робити з токенами.
export default function App({ onLoginSuccess }: Props) {
  const [error, setError] = useState('')
  const [isLoading, setIsLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    setIsLoading(true)
    setError('')

    const form = new FormData(e.currentTarget)
    const login = form.get('login') as string
    const password = form.get('password') as string

    try {
      const res = await axios.post('/api/auth/login', { login, password })
      const data = res.data

      onLoginSuccess?.({
        id: data.id,
        login: data.login,
        role: mapRole(data.role),
        is_default_password: data.is_default_password,
        access_token: data.access_token,
        refresh_token: data.refresh_token,
      })
    } catch (err: unknown) {
      if (axios.isAxiosError(err)) {
        setError(
          err.response?.data?.error ??
          err.response?.data?.message ??
          'Невірний логін або пароль'
        )
      } else {
        setError('Помилка сервера')
      }
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <form onSubmit={handleSubmit}>
      <h1>Вхід</h1>
      <input name="login" type="text" placeholder="Логін" required />
      <input name="password" type="password" placeholder="Пароль" required />
      <button type="submit" disabled={isLoading}>
        {isLoading ? 'Завантаження...' : 'Увійти'}
      </button>
      {error && <p style={{ color: 'red' }}>{error}</p>}
    </form>
  )
}