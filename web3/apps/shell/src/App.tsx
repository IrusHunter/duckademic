import { Suspense, lazy, useEffect } from 'react'
import { BrowserRouter, Routes, Route, useNavigate } from 'react-router-dom'
import { useAuthStore } from './store/authStore'
import ProtectedRoute from './components/ProtectedRoute'
import RoleBasedRedirect from './components/RoleBasedRedirect'
import Header from './components/header/Header'
import { initAuth } from './utils/initAuth'
import { tokenManager } from './auth/tokenManager'
import { isTokenExpired, setAuthCookies } from './utils/cookies'
import css from './App.module.css'
import { setupAxiosInterceptor } from './auth/axiosInterceptor'
import axios from 'axios'

axios.defaults.withCredentials = true
setupAxiosInterceptor()

const AuthApp = lazy(() => import('authApp/AuthApp'))
const ClassroomApp = lazy(() => import('classroomApp/ClassroomApp'))
const HomeApp = lazy(() => import('homeApp/HomeApp'))
const AdminApp = lazy(() => import('adminApp/AdminApp'))

function Routes_() {
  const navigate = useNavigate()

  return (
    <Suspense fallback={<div>Loading...</div>}>
      <Routes>
        <Route path="/" element={<RoleBasedRedirect />} />

        <Route
          path="/login"
          element={
            <AuthApp
              onLoginSuccess={(data) => {
                // Shell — єдиний власник токенів.
                // authApp повернув сирі дані з /api/auth/login,
                // shell зберігає токени і оновлює стор.
                console.log('access_token from login:', data.access_token)
                tokenManager.set(data.access_token)
                setAuthCookies({
                  id: data.id,
                  email: data.login,
                  role: data.role,
                  is_default_password: data.is_default_password,
                }, data.refresh_token)

                useAuthStore.getState().setUser({
                  id: data.id,
                  email: data.login,
                  role: data.role,
                  is_default_password: data.is_default_password,
                })

                navigate(data.role === 'admin' ? '/admin' : '/home')
              }}
            />
          }
        />

        <Route path="/unauthorized" element={<div>Доступ заборонено</div>} />

        <Route
          path="/admin/*"
          element={
            <ProtectedRoute requiredRoles={['admin']}>
              <AdminApp />
            </ProtectedRoute>
          }
        />

        <Route
          path="/home"
          element={
            <ProtectedRoute requiredRoles={['student', 'teacher']}>
              <HomeApp />
            </ProtectedRoute>
          }
        />

        <Route
          path="/classroom/*"
          element={
            <ProtectedRoute requiredRoles={['student', 'teacher']}>
              <ClassroomApp />
            </ProtectedRoute>
          }
        />
      </Routes>
    </Suspense>
  )
}

function App() {
  useEffect(() => {
    const handleVisibility = () => {
      if (document.visibilityState === 'visible') {
        const token = tokenManager.get()
        if (isTokenExpired(token)) {
          tokenManager.refresh().catch(() => { })
        }
      }
    }
    const handleOnline = () => {
  const token = tokenManager.get()
  if (isTokenExpired(token)) {
    tokenManager.refresh().catch(() => {})
  }
}
    const handleLogout = () => { tokenManager.logout() }

    initAuth().finally(() => {
      document.addEventListener('visibilitychange', handleVisibility)
      window.addEventListener('online', handleOnline)
    })

    window.addEventListener('auth:logout', handleLogout)

    return () => {
      document.removeEventListener('visibilitychange', handleVisibility)
      window.removeEventListener('online', handleOnline)
      window.removeEventListener('auth:logout', handleLogout)
    }
  }, [])

  return (
    <BrowserRouter>
      <div className={css.container}>
        <Header />
        <main className={css.main}>
          <Routes_ />
        </main>
      </div>
    </BrowserRouter>
  )
}

export default App