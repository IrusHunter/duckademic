import { useState, useEffect, useRef } from 'react'
import { useAuthStore } from '../../store/authStore'
import { clearUserCookie } from '../../utils/cookies'
import { useNavigate } from 'react-router-dom'
import { useQuery, useQueryClient } from '@tanstack/react-query'
import css from './Header.module.css'
import { LuSearch, LuBell } from 'react-icons/lu'
import { Link } from 'react-router-dom'
import Avatar from '../avatar/Avatar'
import NotificationsPanel from '../notifications/NotificationsPanel'
import { getNotifications } from '../../services/notificationService'

export default function Header() {
  const { isAuthenticated, user } = useAuthStore()
  const clearUser = useAuthStore((s) => s.clearUser)
  const navigate = useNavigate()
  const qc = useQueryClient()

  const [panelOpen, setPanelOpen] = useState(false)
  const bellRef = useRef<HTMLDivElement>(null)

  const { data: notifications = [] } = useQuery({
    queryKey: ['notifications'],
    queryFn: getNotifications,
    enabled: isAuthenticated,
    refetchInterval: 30_000,
  })

  const hasUnread = notifications.some((n) => !n.is_read)

  // Open panel via custom event from QuickActions
  useEffect(() => {
    const handler = () => setPanelOpen(true)
    window.addEventListener('open-notifications', handler)
    return () => window.removeEventListener('open-notifications', handler)
  }, [])

  // Close on outside click
  useEffect(() => {
    if (!panelOpen) return
    const handler = (e: MouseEvent) => {
      if (bellRef.current && !bellRef.current.contains(e.target as Node)) {
        setPanelOpen(false)
      }
    }
    document.addEventListener('mousedown', handler)
    return () => document.removeEventListener('mousedown', handler)
  }, [panelOpen])

  if (!isAuthenticated) return null

  const handleLogout = () => {
    clearUserCookie()
    clearUser()
    navigate('/login')
  }

  const handleBellClick = () => {
    setPanelOpen((prev) => !prev)
    if (!panelOpen) {
      qc.invalidateQueries({ queryKey: ['notifications'] })
    }
  }

  return (
    <header className={css.header}>
      <nav className={css.navigation}>

        <Link to="/home" className={css.logo}>
          <svg width="28" height="32">
            <use href="../../../icons.svg#icon-Logo-1-1"></use>
          </svg>
          Duckademic
        </Link>

        <div className={css.searchBar}>
          <LuSearch className={css.searchIcon} size={23} />
          <input
            type="text"
            placeholder="Search"
            className={css.searchInput}
          />
        </div>

        <div className={css.actions}>
          <div className={css.bellWrapper} ref={bellRef}>
            <button
              className={css.notificationsBtn}
              aria-label="Notifications"
              onClick={handleBellClick}
            >
              <LuBell size={30} className={css.bellIcon} />
              {hasUnread && <span className={css.notificationDot} />}
            </button>
            {panelOpen && (
              <NotificationsPanel onClose={() => setPanelOpen(false)} />
            )}
          </div>

          <div className={css.userInfo}>
            <span className={css.userEmail}>{user?.email}</span>
            <Avatar name={user?.email ?? '?'} size={48} />
          </div>

          <button className={css.logoutBtn} onClick={handleLogout}>
            Log out
          </button>
        </div>

      </nav>
    </header>
  )
}
