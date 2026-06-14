// shell/src/components/Navigation.tsx
import { useNavigate } from 'react-router-dom'
import { useAuthStore } from '../store/authStore'
import { clearUserCookie } from '../utils/cookies'

export default function Navigation() {
  const navigate = useNavigate()
  const { isAuthenticated, user } = useAuthStore()
  const clearUser = useAuthStore((s) => s.clearUser)

  const handleLogout = () => {
    clearUserCookie()
    clearUser()
    navigate('/login')
  }

  return (
    <nav>
      {isAuthenticated ? (
        <>
          <span>{user?.email} ({user?.role})</span>
          <button onClick={handleLogout}>Вийти</button>
        </>
      ) : null}
    </nav>
  )
}