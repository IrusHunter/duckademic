import { formatDeadline, deadlineLevel } from '../../utils/datetime'
import css from '../App/App.module.css'

export interface CourseCardData {
  id: string
  title: string
  teacher: string
  description: string
  assignments: number
  students: number
  averageMark: number
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

  return (
    <li className={`${css.card} ${css[course.colorClass]}`}>
      <div className={css.cardHeader}>
        <div>
          <h2 className={css.courseTitle}>{course.title}</h2>
          <p className={css.teacher}>{course.teacher}</p>
        </div>
        <span className={css.grade}>
          {course.averageMark.toFixed(1)}
        </span>
      </div>

      <p className={css.description}>{course.description}</p>

      <ul className={css.meta}>
        <li className={css.metaItem}>
          <svg className={css.metaIcon} width="16" height="16" aria-hidden="true">
            <use href="/img/icons.svg#icon-SVG-5" />
          </svg>
          <span className={css.metaText}>{course.assignments} assignments</span>
        </li>
        <li className={css.metaItem}>
          <svg className={css.metaIcon} width="16" height="16" aria-hidden="true">
            <use href="/img/icons.svg#icon-SVG-9" />
          </svg>
          <span className={css.metaText}>{course.students} students</span>
        </li>
      </ul>

      {level ? (
        <div className={`${css.deadline} ${DEADLINE_CLASS[level]}`}>
          <svg className={css.metaIcon} width="16" height="16" aria-hidden="true" style={{ color: 'inherit' }}>
            <use href="/img/icons.svg#icon-SVG-11" />
          </svg>
          <p className={css.deadlineText}>Deadline: {formatDeadline(course.nearestDeadline!)}</p>
        </div>
      ) : (
        <div className={css.deadlineNone}>
          <svg width="16" height="16" aria-hidden="true" style={{ color: 'inherit' }}>
            <use href="/img/icons.svg#icon-SVG-11" />
          </svg>
          <p className={css.deadlineText}>No deadline</p>
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