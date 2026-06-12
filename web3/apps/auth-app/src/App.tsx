import { useState } from 'react'
import axios from 'axios'
import css from './App.module.css'

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
function DuckademicLogo() {
  return (
    <div className={css.logo}>
      <svg width="28" height="32" viewBox="0 0 28 32" xmlns="http://www.w3.org/2000/svg">
        <path fill="#253985" d="M1.946 4.397c0-0.061 0.063-0.22 0.314-0.366l7.659-2.992c0.753-0.305 2.361-0.916 2.762-0.916 0.502 0 0.753 0.183 3.076 1.038 1.858 0.684 5.713 2.198 7.408 2.87 0.105 0.102 0.251 0.354 0 0.55s-0.9 0.366-1.193 0.428v4.092c0 0.023-0.001 0.042-0.003 0.059 0.001 0.001 0.002 0.001 0.003 0.002 0.314 0.183 0.502 0.427 0.628 0.733s0.188 1.771 0.188 2.137c0 0.293-0.46 0.529-0.565 0.55-0.167 0.041-0.552 0.122-0.753 0.122-0.251 0-0.565-0.061-0.816-0.244s-0.314-0.428-0.251-1.588c0.050-0.928 0.607-1.567 0.879-1.771 0.042-0.896 0.113-2.773 0.063-3.114-0.063-0.427-0.188-0.366-0.251-0.427-0.050-0.049-0.607 0.102-0.879 0.183-1.046 0.366-3.214 1.136-3.516 1.282-0.377 0.183-3.39 1.527-4.269 1.466-0.703-0.049-1.883-0.509-2.386-0.733l-5.399-2.198c-0.732-0.305-2.235-0.928-2.386-0.977s-0.272-0.143-0.314-0.183z"/>
        <path fill="#ffbc00" d="M5.399 10.137c-0.063-0.611 0.063-1.099 0.628-2.321 0.1-0.098 0.209-0.081 0.251-0.061 0 0 0.364 0.085 2.323 0.916 2.448 1.038 3.139 1.099 4.081 1.099s2.072-0.427 3.578-0.977c1.507-0.55 1.883-0.916 2.135-0.855s0.439 0.672 0.439 0.733 0.439 1.893 0.126 3.542c-0.251 1.319-0.774 2.585-1.005 3.053-0.146 0.305-0.452 0.977-0.502 1.221-0.063 0.305 0 0.366 0.251 0.55s0.942 0.183 3.076 1.038c2.134 0.855 2.7 1.649 2.762 1.649s0.439 0.611 0.816 0.733c0.377 0.122 0.628 0.061 0.816 0 0.151-0.049 0.398-0.224 0.502-0.305l0.502-0.427c0.435-0.37 0.77-0.407 0.883-0.392 0.149-0.052 0.906 0.019 0.749 2.713-0.188 3.237-1.318 4.519-1.381 4.702s-2.134 3.603-6.529 4.641-8.789 0-9.48-0.244c-0.691-0.244-3.014-0.672-4.52-2.626s-1.946-3.542-1.632-5.802c0.251-1.808 1.946-3.725 2.762-4.458 0.167-0.163 0.603-0.623 1.004-1.16s0.293-1.038 0.188-1.221c-0.105-0.285-0.565-0.831-1.569-0.733s-2.218 0.122-2.7 0.122c-0.67-0.020-2.172-0.073-2.825-0.122s-0.9-0.346-0.942-0.489c-0.021-0.102 0-0.232 0-0.428 0-0.183 0.251-0.366 0.439-0.489l3.265-1.527c0.251-0.122 0.829-0.464 1.13-0.855 0.377-0.489 0.439-0.611 0.377-1.221z"/>
        <circle fill="#fff" cx="8.915" cy="11.725" r="1.13"/>
        <path fill="#f96706" d="M6.843 12.458c0.251 0.366 1.005 2.26 0.188 2.687-0.251 0.061-1.256 0.122-1.256 0.122s-1.253 0.080-3.139 0c-1.444-0.061-1.695-0.026-2.072-0.183-0.439-0.183-0.502-0.516-0.502-0.794 0-0.244 0.251-0.489 0.565-0.672 0.215-0.126 1.004-0.448 1.381-0.611 0.502-0.244 0.716-0.269 1.821-0.855 1.381-0.733 1.004-0.916 1.256-0.855l0 0c0.311 0.075 0.753 0.183 1.004 0.366s0.502 0.428 0.753 0.794z"/>
      </svg>
      <span className={css.logoText}>Duckademic</span>
    </div>
  )
}

export default function App({ onLoginSuccess }: Props) {
  const [error, setError] = useState('')
  const [isLoading, setIsLoading] = useState(false)

  const handleSubmit = async (e: React.SyntheticEvent<HTMLFormElement>) => {
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
    <div className={css.page}>
      <form onSubmit={handleSubmit} className={css.form}>
        <DuckademicLogo />
        <h1 className={css.title}>Вхід</h1>
        <input className={css.input} name="login" type="text" placeholder="Логін" required />
        <input className={css.input} name="password" type="password" placeholder="Пароль" required />
        <button className={css.button} type="submit" disabled={isLoading}>
          {isLoading ? 'Завантаження...' : 'Увійти'}
        </button>
        {error && <p className={css.error}>{error}</p>}
      </form>
    </div>
  )
}