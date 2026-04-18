import { Suspense, lazy, useEffect } from 'react'
import { BrowserRouter, Routes, Route, useNavigate } from 'react-router-dom'
import { useAuthStore } from './store/authStore'
import type { User } from './store/authStore'
import ProtectedRoute from './components/ProtectedRoute'

const AuthApp = lazy(() => import('authApp/AuthApp'))
const ClassroomApp = lazy(() => import('classroomApp/ClassroomApp'))
const HomeApp = lazy(() => import('homeApp/HomeApp'))
const AdminApp = lazy(() => import('adminApp/AdminApp'))

function Routes_() {
  const navigate = useNavigate()
  const setUser = useAuthStore((s) => s.setUser)
  const clearUser = useAuthStore((s) => s.clearUser)

  useEffect(() => {
    fetch('/api/auth/session')
      .then((res) => res.ok ? res.json() : null)
      .then((user) => {
        if (user) setUser(user)
        else clearUser()
      })
      .catch(() => clearUser())
  }, [])

  return (
    <Suspense fallback={<div>Loading...</div>}>
      <Routes>
        {/* Публічний маршрут */}
        <Route
          path="/login"
          element={
            <AuthApp
              onLoginSuccess={(user: User) => {
                setUser(user)
                // Адмін іде в /admin, решта на /
                if (user.role === 'admin') {
                  navigate('/admin')
                } else {
                  navigate('/')
                }
              }}
            />
          }
        />

        <Route
          path="/unauthorized"
          element={<div>Доступ заборонено</div>}
        />

        {/* Тільки для адміна */}
        <Route
          path="/admin/*"
          element={
            <ProtectedRoute requiredRole="admin">
              <AdminApp />
            </ProtectedRoute>
          }
        />

        {/* Для студента і викладача */}
        <Route
          path="/"
          element={
            <ProtectedRoute>
              <HomeApp />
            </ProtectedRoute>
          }
        />

        <Route
          path="/classroom/*"
          element={
            <ProtectedRoute>
              <ClassroomApp />
            </ProtectedRoute>
          }
        />
      </Routes>
    </Suspense>
  )
}

function App() {
  return (
    <BrowserRouter>
      <Routes_ />
    </BrowserRouter>
  )
}

export default App