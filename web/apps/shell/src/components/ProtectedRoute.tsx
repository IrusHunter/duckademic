import { Navigate } from 'react-router-dom'
import { useAuthStore } from '../store/authStore'

type Props = {
  children: React.ReactNode
  requiredRoles?: Array<'admin' | 'student' | 'teacher'>
}

// ✅ FIX #1: перевіряємо isInitialized перед будь-яким redirect.
// Без цього після hard reload юзер миттєво летить на /login,
// поки initAuth ще виконує refresh у фоні.
export default function ProtectedRoute({ children, requiredRoles }: Props) {
  const { isAuthenticated, isInitialized, user } = useAuthStore()

  if (!isInitialized) {
    return null // або <Spinner /> — чекаємо завершення initAuth
  }

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />
  }

  if (requiredRoles && user?.role && !requiredRoles.includes(user.role)) {
    return <Navigate to="/unauthorized" replace />
  }

  return <>{children}</>
}