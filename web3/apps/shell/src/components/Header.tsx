import { useAuthStore } from '../store/authStore'
import { clearUserCookie } from '../utils/cookies'
import { useNavigate } from 'react-router-dom'

export default function Header() {
  const { isAuthenticated, user } = useAuthStore()
  const clearUser = useAuthStore((s) => s.clearUser)
  const navigate = useNavigate()

  if (!isAuthenticated) return null

  const handleLogout = () => {
    clearUserCookie()
    clearUser()
    navigate('/login')
  }

  return (
    <header style={{
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'space-between',
      padding: '0 24px',
      height: '60px',
      borderBottom: '1px solid #eee',
      backgroundColor: '#fff',
      position: 'sticky',
      top: 0,
      zIndex: 100,
    }}>
      <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
        <span>🎓</span>
        <strong>Duckademic</strong>
      </div>
      <input
        type="text"
        placeholder="Search"
        style={{ width: '300px', padding: '8px 12px', borderRadius: '8px', border: '1px solid #ddd' }}
      />
      <div style={{ display: 'flex', alignItems: 'center', gap: '16px' }}>
        <span>🔔</span>
        <span style={{ fontSize: '14px', color: '#666' }}>{user?.email}</span>
        <button
          onClick={handleLogout}
          style={{ padding: '8px 16px', borderRadius: '8px', border: '1px solid #ddd', cursor: 'pointer' }}
        >
          Вийти
        </button>
      </div>
    </header>
  )
}