import Icon from '../Icon/Icon'
import css from '../App/App.module.css'

export default function QuickActions() {
  return (
    <section className={css.quickActions}>
      <h2 className={css.title}>Quick Actions</h2>
      <ul>
        <li>
          <button type="button">
            <Icon id="icon-SVG-4" size={16} />
            Show rating
          </button>
        </li>
        <li>
          <button type="button">
            <Icon id="icon-SVG-3" size={16} />
            Check Schedule
          </button>
        </li>
        <li>
          <button type="button">
            <Icon id="icon-bell-1" size={16} />
            Notifications
          </button>
        </li>
      </ul>
    </section>
  )
}
