import { Navigate } from 'react-router-dom'
import { useAuthStore } from '../store/authStore'

type Props = {
  children: React.ReactNode
  requiredRoles?: Array<'admin' | 'student' | 'teacher'>
}

export default function ProtectedRoute({ children, requiredRoles }: Props) {
  const { isAuthenticated, user } = useAuthStore()

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />
  }

  if (requiredRoles && user?.role && !requiredRoles.includes(user.role)) {
    return <Navigate to="/unauthorized" replace />
  }

  return <>{children}</>
}