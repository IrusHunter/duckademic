import { Suspense, lazy } from 'react'
import { BrowserRouter, Routes, Route, useNavigate, Navigate } from 'react-router-dom'
import { useAuthStore } from './store/authStore'
import type { User } from './store/authStore'
import ProtectedRoute from './components/ProtectedRoute'
import { getUserFromCookie } from './utils/cookies'
import Navigation from './components/Navigation'

const AuthApp = lazy(() => import('authApp/AuthApp'))
const ClassroomApp = lazy(() => import('classroomApp/ClassroomApp'))
const HomeApp = lazy(() => import('homeApp/HomeApp'))
const AdminApp = lazy(() => import('adminApp/AdminApp'))

const initAuth = () => {
  const user = getUserFromCookie()
  if (user) {
    useAuthStore.getState().setUser(user)
  }
}

initAuth()

// Редіректить на потрібну сторінку залежно від ролі
function RoleBasedRedirect() {
  const { isAuthenticated, user } = useAuthStore()

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />
  }

  if (user?.role === 'admin') {
    return <Navigate to="/admin" replace />
  }

  return <Navigate to="/home" replace />
}

function Routes_() {
  const navigate = useNavigate()

  return (
    <Suspense fallback={<div>Loading...</div>}>
      <Routes>
        {/* / — розподіляє по ролях */}
        <Route path="/" element={<RoleBasedRedirect />} />

        <Route
          path="/login"
          element={
            <AuthApp
              onLoginSuccess={(user: User) => {
                useAuthStore.getState().setUser(user)
                if (user.role === 'admin') {
                  navigate('/admin')
                } else {
                  navigate('/home')
                }
              }}
            />
          }
        />

        <Route path="/unauthorized" element={<div>Доступ заборонено</div>} />

        {/* Тільки адмін */}
        <Route
          path="/admin/*"
          element={
            <ProtectedRoute requiredRole="admin">
              <AdminApp />
            </ProtectedRoute>
          }
        />

        {/* Тільки student і teacher */}
        <Route
          path="/home"
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
      <Navigation />
      <Routes_ />
    </BrowserRouter>
  )
}

export default App