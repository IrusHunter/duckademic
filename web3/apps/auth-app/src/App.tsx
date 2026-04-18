import { useState } from 'react'
import './App.css'

type User = {
  id: string
  email: string
  role: 'admin' | 'user'
}

type Props = {
  onLoginSuccess?: (user: User) => void
}

type AuthMode = 'login' | 'register'

export default function App({ onLoginSuccess }: Props) {
  const [mode, setMode] = useState<AuthMode>('login')
  const [error, setError] = useState('')

  const handleLogin = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    const form = new FormData(e.currentTarget)

    try {
      const res = await fetch('http://your-api/auth/login', {
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          email: form.get('email'),
          password: form.get('password'),
        }),
      })

      if (res.ok) {
        const user = await res.json()
        onLoginSuccess?.(user)
      } else {
        setError('Невірний email або пароль')
      }
    } catch {
      setError('Помилка сервера')
    }
  }

  const handleRegister = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    const form = new FormData(e.currentTarget)

    const password = form.get('password') as string
    const confirmPassword = form.get('confirmPassword') as string

    if (password !== confirmPassword) {
      setError('Паролі не збігаються')
      return
    }

    try {
      const res = await fetch('http://your-api/auth/register', {
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          email: form.get('email'),
          password,
        }),
      })

      if (res.ok) {
        const user = await res.json()
        onLoginSuccess?.(user)
      } else {
        const data = await res.json()
        setError(data.error ?? 'Помилка реєстрації')
      }
    } catch {
      setError('Помилка сервера')
    }
  }

  return (
    <div>
      <div>
        <button onClick={() => { setMode('login'); setError('') }}>
          Вхід
        </button>
        <button onClick={() => { setMode('register'); setError('') }}>
          Реєстрація
        </button>
      </div>

      {mode === 'login' ? (
        <form onSubmit={handleLogin}>
          <h1>Вхід</h1>
          <input name="email" type="email" placeholder="Email" required />
          <input name="password" type="password" placeholder="Пароль" required />
          <button type="submit">Увійти</button>
        </form>
      ) : (
        <form onSubmit={handleRegister}>
          <h1>Реєстрація</h1>
          <input name="email" type="email" placeholder="Email" required />
          <input name="password" type="password" placeholder="Пароль" required />
          <input name="confirmPassword" type="password" placeholder="Повторіть пароль" required />
          <button type="submit">Зареєструватись</button>
        </form>
      )}

      {error && <p>{error}</p>}
    </div>
  )
}