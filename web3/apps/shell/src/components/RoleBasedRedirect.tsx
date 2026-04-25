import { Navigate } from 'react-router-dom'
import { useAuthStore } from '../store/authStore'

// ✅ FIX #1: чекаємо isInitialized так само, як ProtectedRoute.
// Інакше при першому рендері isAuthenticated=false → негайний redirect на /login.
export default function RoleBasedRedirect() {
  const { isAuthenticated, isInitialized, user } = useAuthStore()

  if (!isInitialized) {
    return null // або <Spinner />
  }

  if (!isAuthenticated) return <Navigate to="/login" replace />
  if (user?.role === 'admin') return <Navigate to="/admin" replace />
  return <Navigate to="/home" replace />
}