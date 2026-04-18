import { Navigate } from 'react-router-dom'
import { useAuthStore } from '../store/authStore'

export default function RoleBasedRedirect() {
  const { isAuthenticated, user } = useAuthStore()

  if (!isAuthenticated) return <Navigate to="/login" replace />
  if (user?.role === 'admin') return <Navigate to="/admin" replace />
  return <Navigate to="/home" replace />
}