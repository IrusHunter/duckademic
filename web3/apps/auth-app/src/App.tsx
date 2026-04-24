import { useState } from 'react'
import axios from 'axios'

type User = {
  id: string
  login: string
  role: 'admin' | 'student' | 'teacher'
  is_default_password: boolean
}

type Props = {
  onLoginSuccess?: (user: User) => void
}

// Зберігаємо токени і юзера в cookie
const setAuthCookies = (user: User, accessToken: string, refreshToken: string) => {
  document.cookie = `access_token=${accessToken}; path=/; Max-Age=30; SameSite=Strict`
  document.cookie = `refresh_token=${refreshToken}; path=/; Max-Age=${86400 * 30}; SameSite=Strict`
  const userValue = encodeURIComponent(JSON.stringify({
    id: user.id,
    email: user.login,
    role: mapRole(user.role),
  }))
  document.cookie = `auth_user=${userValue}; path=/; Max-Age=86400; SameSite=Strict`
  localStorage.setItem('access_token', accessToken)
  localStorage.setItem('refresh_token', refreshToken)
}

// Маппінг ролей з бекенду на фронтенд
const mapRole = (role: string): 'admin' | 'student' | 'teacher' => {
  const roleMap: Record<string, 'admin' | 'student' | 'teacher'> = {
    admin: 'admin',
    student: 'student',
    teacher: 'teacher',
    // якщо бекенд повертає інші назви — додай тут
  }
  return roleMap[role.toLowerCase()] ?? 'student'
}

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

      // Бекенд повертає: id, login, role, is_default_password, access_token, refresh_token
      const user: User = {
        id: data.id,
        login: data.login,
        role: mapRole(data.role),
        is_default_password: data.is_default_password,
      }

      setAuthCookies(user, data.access_token, data.refresh_token)

      onLoginSuccess?.({
        id: user.id,
        login: user.login,
        role: user.role,
        is_default_password: user.is_default_password,
      })
    } catch (err: unknown) {
      if (axios.isAxiosError(err)) {
        setError(err.response?.data?.error ?? err.response?.data?.message ?? 'Невірний логін або пароль')
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