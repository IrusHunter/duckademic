import { useNavigate } from 'react-router-dom'
import { LuStar, LuCalendar, LuBell } from 'react-icons/lu'
import css from '../App/App.module.css'

export default function QuickActions() {
  const navigate = useNavigate()

  const openNotifications = () => {
    window.dispatchEvent(new CustomEvent('open-notifications'))
  }

  return (
    <section className={css.quickActions}>
      <h2 className={css.title}>Quick Actions</h2>
      <ul>
        <li>
          <button type="button" onClick={() => navigate('/grades')}>
            <LuStar size={16} />
            Show rating
          </button>
        </li>
        <li>
          <button type="button" onClick={() => navigate('/schedule')}>
            <LuCalendar size={16} />
            Check Schedule
          </button>
        </li>
        <li>
          <button type="button" onClick={openNotifications}>
            <LuBell size={16} />
            Notifications
          </button>
        </li>
      </ul>
    </section>
  )
}
