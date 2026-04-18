import { Suspense, lazy, useEffect } from 'react'
import { BrowserRouter, Routes, Route, useNavigate } from 'react-router-dom'
import { useAuthStore } from './store/authStore'
import type { User } from './store/authStore'
import ProtectedRoute from './components/ProtectedRoute'

const AuthApp = lazy(() => import('authApp/AuthApp'))
const ClassroomApp = lazy(() => import('classroomApp/ClassroomApp'))
const HomeApp = lazy(() => import('homeApp/HomeApp'))

function Routes_() {
  const navigate = useNavigate()
  const setUser = useAuthStore((s) => s.setUser)
  const clearUser = useAuthStore((s) => s.clearUser)

  useEffect(() => {
    // Перевіряємо сесію при старті
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
        <Route
          path="/login"
          element={
            <AuthApp
              onLoginSuccess={(user: User) => {
                setUser(user)
                navigate('/')
              }}
            />
          }
        />

        <Route path="/unauthorized" element={<div>Доступ заборонено</div>} />

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