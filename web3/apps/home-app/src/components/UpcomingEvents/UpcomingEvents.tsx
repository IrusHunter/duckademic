import type { Task } from '../../types/course'
import { daysUntil, formatDeadline } from '../../utils/datetime'
import EmptyState from '../EmptyState/EmptyState'
import css from './UpcomingEvents.module.css'

interface UpcomingEventsProps {
  tasks: Task[]
}

function urgencyClass(deadline: string): string {
  const days = daysUntil(deadline)
  if (days <= 1) return css.urgent
  if (days <= 3) return css.soon
  return css.normal
}

export default function UpcomingEvents({ tasks }: UpcomingEventsProps) {
  if (tasks.length === 0) {
    return <EmptyState title="Найближчих дедлайнів немає" />
  }

  return (
    <ul className={css.list}>
      {tasks.map((task) => (
        <li key={task.id} className={css.item}>
          <span className={`${css.marker} ${urgencyClass(task.deadline)}`} />
          <div className={css.content}>
            <p className={css.title}>{task.title}</p>
            <p className={css.deadline}>{formatDeadline(task.deadline)}</p>
          </div>
          <span className={css.mark}>{task.max_mark} б.</span>
        </li>
      ))}
    </ul>
  )
}
