import type { LessonOccurrence, LessonStatus } from '../../types/schedule'
import { formatTime } from '../../utils/datetime'
import css from './LessonCard.module.css'

interface LessonCardProps {
  lesson: LessonOccurrence
}

const STATUS_LABEL: Record<LessonStatus, string> = {
  scheduled: 'Заплановано',
  completed: 'Проведено',
  canceled: 'Скасовано',
}

export default function LessonCard({ lesson }: LessonCardProps) {
  const load = lesson.study_load

  return (
    <article className={`${css.card} ${css[lesson.status]}`}>
      <div className={css.time}>{formatTime(lesson.date)}</div>

      <div className={css.body}>
        <h3 className={css.discipline}>
          {load?.discipline_name ?? 'Дисципліна'}
        </h3>

        <div className={css.meta}>
          {load?.lesson_type_name && (
            <span className={css.type}>{load.lesson_type_name}</span>
          )}
          {load?.teacher_name && <span>{load.teacher_name}</span>}
          {load?.student_group_name && <span>{load.student_group_name}</span>}
        </div>
      </div>

      <span className={`${css.status} ${css[`status_${lesson.status}`]}`}>
        {STATUS_LABEL[lesson.status]}
      </span>
    </article>
  )
}
