import { Suspense, lazy } from 'react'
import { BrowserRouter, Routes, Route, useNavigate } from 'react-router-dom'
import { useAuthStore } from './store/authStore'
import type { User } from './store/authStore'
import ProtectedRoute from './components/ProtectedRoute'
import { getUserFromCookie } from './utils/cookies'
import Navigation from './components/Navigation'

const AuthApp = lazy(() => import('authApp/AuthApp'))
const ClassroomApp = lazy(() => import('classroomApp/ClassroomApp'))
const HomeApp = lazy(() => import('homeApp/HomeApp'))
const AdminApp = lazy(() => import('adminApp/AdminApp'))

// Читаємо куку одразу при ініціалізації store — до будь-якого рендеру
const initAuth = () => {
  const user = getUserFromCookie()
  if (user) {
    useAuthStore.getState().setUser(user)
  }
}

initAuth() // викликаємо одразу при імпорті модуля

function Routes_() {
  const navigate = useNavigate()

  return (
    <Suspense fallback={<div>Loading...</div>}>
      <Routes>
        <Route
          path="/login"
          element={
            <AuthApp
              onLoginSuccess={(user: User) => {
                useAuthStore.getState().setUser(user)
                if (user.role === 'admin') {
                  navigate('/admin')
                } else {
                  navigate('/')
                }
              }}
            />
          }
        />

        <Route path="/unauthorized" element={<div>Доступ заборонено</div>} />

        <Route
          path="/admin/*"
          element={
            <ProtectedRoute requiredRole="admin">
              <AdminApp />
            </ProtectedRoute>
          }
        />

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
      <Navigation />
    </BrowserRouter>
  )
}

export default App