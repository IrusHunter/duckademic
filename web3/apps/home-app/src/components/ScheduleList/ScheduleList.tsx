import type { LessonOccurrence } from '../../types/schedule'
import { formatDay, groupByDay } from '../../utils/datetime'
import LessonCard from '../LessonCard/LessonCard'
import EmptyState from '../EmptyState/EmptyState'
import css from './ScheduleList.module.css'

interface ScheduleListProps {
  lessons: LessonOccurrence[]
}

export default function ScheduleList({ lessons }: ScheduleListProps) {
  if (lessons.length === 0) {
    return (
      <EmptyState
        title="Занять на цей тиждень немає"
        hint="Коли з'являться нові заняття, вони відобразяться тут."
      />
    )
  }

  // Сортуємо за часом і групуємо за днями.
  const sorted = [...lessons].sort(
    (a, b) => new Date(a.date).getTime() - new Date(b.date).getTime(),
  )
  const days = groupByDay(sorted, (l) => l.date)

  return (
    <div className={css.list}>
      {days.map(([dayKey, dayLessons]) => (
        <section key={dayKey} className={css.day}>
          <h2 className={css.dayTitle}>{formatDay(dayLessons[0].date)}</h2>
          <div className={css.lessons}>
            {dayLessons.map((lesson) => (
              <LessonCard key={lesson.id} lesson={lesson} />
            ))}
          </div>
        </section>
      ))}
    </div>
  )
}
