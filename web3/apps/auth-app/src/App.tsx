import { useState } from 'react'
import './App.css'

type User = {
  id: string
  email: string
  role: 'admin' | 'student' | 'teacher'
}

type Props = {
  onLoginSuccess?: (user: User) => void
}

const MOCK_USERS: Record<string, User> = {
  'admin@gmail.com': { id: '1', email: 'admin@gmail.com', role: 'admin' },
  'student@gmail.com': { id: '2', email: 'student@gmail.com', role: 'student' },
  'teacher@gmail.com': { id: '3', email: 'teacher@gmail.com', role: 'teacher' },
}

const MOCK_PASSWORDS: Record<string, string> = {
  'admin@gmail.com': 'admin',
  'student@gmail.com': 'student',
  'teacher@gmail.com': 'teacher',
}

// Зберігаємо user в cookie
const setUserCookie = (user: User) => {
  const value = encodeURIComponent(JSON.stringify(user))
  // Max-Age=86400 — cookie живе 1 день
  document.cookie = `auth_user=${value}; path=/; Max-Age=86400; SameSite=Strict`
}

export default function App({ onLoginSuccess }: Props) {
  const [error, setError] = useState('')

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    const form = new FormData(e.currentTarget)

    const email = form.get('email') as string
    const password = form.get('password') as string

    const user = MOCK_USERS[email]
    const correctPassword = MOCK_PASSWORDS[email]

    if (user && password === correctPassword) {
      setUserCookie(user) // ← зберігаємо в cookie
      onLoginSuccess?.(user)
    } else {
      setError('Невірний email або пароль')
    }
  }

  return (
    <form onSubmit={handleSubmit}>
      <h1>Вхід</h1>
      <input name="email" type="email" placeholder="Email" required />
      <input name="password" type="password" placeholder="Пароль" required />
      <button type="submit">Увійти</button>
      {error && <p>{error}</p>}
      <div style={{ marginTop: '20px', fontSize: '12px', color: '#888' }}>
        <p>Тестові акаунти:</p>
        <p>admin@gmail.com / admin</p>
        <p>student@gmail.com / student</p>
        <p>teacher@gmail.com / teacher</p>
      </div>
    </form>
  )
}