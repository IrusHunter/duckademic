import { Suspense, lazy } from 'react'
import { BrowserRouter, Routes, Route, useNavigate } from 'react-router-dom'
import { useAuthStore } from './store/authStore'
import type { User } from './store/authStore'
import ProtectedRoute from './components/ProtectedRoute'
import RoleBasedRedirect from './components/RoleBasedRedirect'
import Header from './components/Header'
import { initAuth } from './utils/initAuth'

const AuthApp = lazy(() => import('authApp/AuthApp'))
const ClassroomApp = lazy(() => import('classroomApp/ClassroomApp'))
const HomeApp = lazy(() => import('homeApp/HomeApp'))
const AdminApp = lazy(() => import('adminApp/AdminApp'))

initAuth()

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
              onLoginSuccess={(user: User) => {
                useAuthStore.getState().setUser(user)
                navigate(user.role === 'admin' ? '/admin' : '/home')
              }}
            />
          }
        />
        <Route path="/unauthorized" element={<div>Доступ заборонено</div>} />
        <Route path="/admin/*" element={
          <ProtectedRoute requiredRoles={['admin']}>
            <AdminApp />
          </ProtectedRoute>
        } />
        <Route path="/home" element={
          <ProtectedRoute requiredRoles={['student', 'teacher']}>
            <HomeApp />
          </ProtectedRoute>
        } />
        <Route path="/classroom/*" element={
          <ProtectedRoute requiredRoles={['student', 'teacher']}>
            <ClassroomApp />
          </ProtectedRoute>
        } />
      </Routes>
    </Suspense>
  )
}

function App() {
  return (
    <BrowserRouter>
      <Header />
      <Routes_ />
    </BrowserRouter>
  )
}

export default App