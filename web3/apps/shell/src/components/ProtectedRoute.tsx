import { Navigate } from 'react-router-dom'
import { useAuthStore } from '../store/authStore'

type Props = {
  children: React.ReactNode
  requiredRole?: 'admin' | 'student' | 'teacher'
}

export default function ProtectedRoute({ children, requiredRole }: Props) {
  const { isAuthenticated, user } = useAuthStore()

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />
  }

  if (requiredRole && user?.role !== requiredRole) {
    return <Navigate to="/unauthorized" replace />
  }

  return <>{children}</>
}