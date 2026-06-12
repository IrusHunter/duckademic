import Icon from '../Icon/Icon'
import { formatDeadline } from '../../utils/datetime'
import type { Task } from '../../types/dashboard'
import css from '../App/App.module.css'

interface UpcomingProps {
  tasks: Task[]
}

export default function Upcoming({ tasks }: UpcomingProps) {
  return (
    <section className={css.upcoming}>
      <div className={css.titleContainer}>
        <Icon id="icon-SVG-3" size={16} />
        <h2 className={css.title}>Upcoming</h2>
      </div>

      {tasks.length === 0 ? (
        <p className={css.sectionEmpty}>Найближчих подій немає</p>
      ) : (
        <ul className={css.upcomingList}>
          {tasks.map((task) => (
            <li className={css.upcomingItem} key={task.id}>
              <h3 className={css.title}>{task.title}</h3>
              {/* Завдання = Assignment → помаранчевий час, без бейджа типу */}
              <div className={`${css.timeContainer} ${css.timeAssignment}`}>
                <p>{formatDeadline(task.deadline)}</p>
              </div>
            </li>
          ))}
        </ul>
      )}
    </section>
  )
}
