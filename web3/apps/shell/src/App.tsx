import { Suspense, lazy, useEffect } from 'react'
import { BrowserRouter, Routes, Route, useNavigate } from 'react-router-dom'
import { useAuthStore } from './store/authStore'
import ProtectedRoute from './components/ProtectedRoute'
import RoleBasedRedirect from './components/RoleBasedRedirect'
import Header from './components/header/Header'
import { initAuth } from './utils/initAuth'
import css from './App.module.css'
import { setupAxiosInterceptor } from './auth/axiosInterceptor'
import { tokenManager } from './auth/tokenManager'

setupAxiosInterceptor();

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
              onLoginSuccess={(user) => {
                useAuthStore.getState().setUser({
                  id: user.id,
                  email: user.login,
                  role: user.role,
                })
                navigate(user.role === 'admin' ? '/admin' : '/home')
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
        tokenManager.refresh().catch(() => {});
      }
    };

    const handleOnline = () => {
      tokenManager.refresh().catch(() => {});
    };
    
    document.addEventListener('visibilitychange', handleVisibility);
    window.addEventListener('online', handleOnline);

    return () => {
      document.removeEventListener('visibilitychange', handleVisibility);
      window.removeEventListener('online', handleOnline);
    };
  }, []);
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