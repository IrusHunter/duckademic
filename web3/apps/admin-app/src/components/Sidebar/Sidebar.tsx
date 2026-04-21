import { useNavigate } from 'react-router-dom'
import { SERVICES } from '../../config/services'
import css from '../App/App.module.css'

export function Sidebar() {
  const navigate = useNavigate()
  const currentPath = window.location.pathname

  return (
    <aside className={css.aside}>
      <p className={css.servicesHeader}>Services</p>
      <ul className={css.servicesList}>
        {SERVICES.map(service => {
          const isActive = currentPath.includes(`/admin/${service.key}`)
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
    </aside>
  )
}