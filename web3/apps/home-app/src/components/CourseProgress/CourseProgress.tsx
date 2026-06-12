import Icon from '../Icon/Icon'
import type { CourseProgress as CourseProgressType } from '../../types/dashboard'
import css from '../App/App.module.css'

interface CourseProgressProps {
  items: CourseProgressType[]
}

export default function CourseProgress({ items }: CourseProgressProps) {
  return (
    <section className={css.courseProgress}>
      <div className={css.titleContainer}>
        <Icon id="icon-SVG-10" size={16} />
        <h2 className={css.title}>Course Progress</h2>
      </div>

      {items.length === 0 ? (
        <p className={css.sectionEmpty}>No courses yet</p>
      ) : (
        <ul>
          {items.map((c) => {
            const percent = Math.round((c.complete_rate ?? 0) * 100)
            const accuracy = Math.round((c.complete_accuracy ?? 0) * 100)
            const good = accuracy >= 80
            return (
              <li key={c.id}>
                <div className={css.courseHeader}>
                  <h3>{c.name}</h3>
                  <p className={good ? css.green : ''}>{accuracy}%</p>
                </div>
                <div className={css.progressBar} aria-label={`${c.name} progress`}>
                  <div
                    className={css.progressBarCompleted}
                    style={{ width: `${percent}%` }}
                  />
                </div>
                <p className={css.progressDescription}>{percent}% complete</p>
              </li>
            )
          })}
        </ul>
      )}
    </section>
  )
}
