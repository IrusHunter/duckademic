import { useNavigate, useLocation } from 'react-router-dom'
import { LuCalendarClock } from 'react-icons/lu'
import { SERVICES } from '../../config/services'
import css from '../App/App.module.css'

interface SidebarProps {
  onOpenGenerator: () => void
}

export function Sidebar({ onOpenGenerator }: SidebarProps) {
  const navigate = useNavigate()
  const location = useLocation()

  // Сегмент сервісу в URL: /admin/<serviceKey>/...
  // Точне порівняння сегмента, а не includes — інакше '/admin/student'
  // підсвічувало б і 'student', і 'student-group'.
  const activeKey = location.pathname.split('/')[2] ?? ''

  return (
    <aside className={css.aside}>
      <p className={css.servicesHeader}>Services</p>
      <ul className={css.servicesList}>
        {SERVICES.map(service => {
          const isActive = activeKey === service.key
          return (
            <li key={service.key}>
              <button
                onClick={() => navigate(`/admin/${service.key}`)}
                className={`${css.button} ${isActive ? css.active : ''}`}
              >
                <span className={css.serviceIcon}>{service.icon}</span>
                <span className={css.serviceTitle}>{service.label}</span>
              </button>
            </li>
          )
        })}
      </ul>

      {/* Кнопка генерації розкладу — під списком сервісів */}
      <button
        type="button"
        className={css.generatorButton}
        onClick={onOpenGenerator}
      >
        <span className={css.serviceIcon}><LuCalendarClock /></span>
        <span className={css.serviceTitle}>Schedule Generation</span>
      </button>
    </aside>
  )
}
