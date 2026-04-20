import { useAuthStore } from '../../store/authStore'
import { clearUserCookie } from '../../utils/cookies'
import { useNavigate } from 'react-router-dom'
import css from './Header.module.css'
import { LuSearch, LuBell } from 'react-icons/lu'
import { Link } from "react-router-dom";

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
          <button className={css.notificationsBtn} aria-label="Notifications">
            <LuBell size={30} className={css.bellIcon} />
            <span className={css.notificationDot} />
          </button>

          <div className={css.userInfo}>
            <span className={css.userEmail}>{user?.email}</span>
            <div className={css.avatar}>
              {user?.email?.[0]?.toUpperCase() ?? '?'}
            </div>
          </div>

          <button className={css.logoutBtn} onClick={handleLogout}>
            Вийти
          </button>
        </div>

      </nav>
    </header>
  )
}