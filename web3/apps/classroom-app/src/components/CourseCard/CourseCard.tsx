import { formatDeadline, deadlineLevel } from '../../utils/datetime'
import css from '../App/App.module.css'

export interface CourseCardData {
  id: string
  title: string
  ratePercent: number        // 0..100 — проходження
  accuracyPercent: number    // 0..100 — точність
  assignments?: number       // лише якщо вдалося зіставити завдання з курсом
  nearestDeadline?: string   // ISO або undefined
  colorClass: 'cardBlue' | 'cardMint' | 'cardRose'
}

const DEADLINE_CLASS = {
  warning: css.deadlineWarning,
  info: css.deadlineInfo,
  danger: css.deadlineDanger,
} as const

export default function CourseCard({ course }: { course: CourseCardData }) {
  const level = course.nearestDeadline ? deadlineLevel(course.nearestDeadline) : null
  const goodAccuracy = course.accuracyPercent >= 80

  return (
    <li className={`${css.card} ${css[course.colorClass]}`}>
      <div className={css.cardHeader}>
        <div>
          <h2 className={css.courseTitle}>{course.title}</h2>
        </div>
        <span className={`${css.grade} ${goodAccuracy ? css.gradeGreen : ''}`}>
          {course.accuracyPercent}%
        </span>
      </div>

      <p className={css.description}>{course.ratePercent}% complete</p>

      {/* Прогрес-бар (реальний complete_rate студента) */}
      <div className={css.progressBar} aria-label={`${course.title} progress`}>
        <div className={css.progressBarFill} style={{ width: `${course.ratePercent}%` }} />
      </div>

      {(course.assignments !== undefined || course.nearestDeadline) && (
        <ul className={css.meta}>
          {course.assignments !== undefined && (
            <li className={css.metaItem}>
              <svg className={css.metaIcon} width="16" height="16" aria-hidden="true">
                <use href="/img/icons.svg#icon-SVG-5" />
              </svg>
              <span className={css.metaText}>{course.assignments} assignments</span>
            </li>
          )}
        </ul>
      )}

      {course.nearestDeadline && level && (
        <div className={`${css.deadline} ${DEADLINE_CLASS[level]}`}>
          <svg className={css.metaIcon} width="16" height="16" aria-hidden="true" style={{ color: 'inherit' }}>
            <use href="/img/icons.svg#icon-SVG-11" />
          </svg>
          <p className={css.deadlineText}>Deadline: {formatDeadline(course.nearestDeadline)}</p>
        </div>
      )}

      <div className={css.actions}>
        <button className={`${css.btn} ${css.btnPrimary}`} type="button">
          Open Course
        </button>
        <button className={`${css.btn} ${css.btnSecondary}`} type="button">
          Assignments
        </button>
      </div>
    </li>
  )
}
