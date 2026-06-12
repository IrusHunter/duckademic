import { useMemo } from 'react'
import { useQuery } from '@tanstack/react-query'
import { getStudentCourses } from '../../services/courseService'
import CourseCard from '../CourseCard/CourseCard'
import type { CourseCardData } from '../CourseCard/CourseCard'
import Loader from '../Loader/Loader'
import ErrorMessage from '../ErrorMessage/ErrorMessage'
import css from './App.module.css'

const COLORS = ['cardBlue', 'cardMint', 'cardRose'] as const

function Courses() {
  const { data, isLoading, isError, refetch } = useQuery({
    queryKey: ['student-courses'],
    queryFn: getStudentCourses,
  })

  const cards: CourseCardData[] = useMemo(() => {
    return (data ?? []).map((c, i) => ({
      id: c.id,
      title: c.name,
      teacher: c.teacher_name,
      description: c.description,
      assignments: c.assignments_count,
      students: c.student_count,
      averageMark: c.average_mark,
      nearestDeadline: c.upcoming_deadline && !c.upcoming_deadline.startsWith('0001-') ? c.upcoming_deadline : undefined,
      colorClass: COLORS[i % COLORS.length],
    }))
  }, [data])

  return (
    <main className={css.page}>
      <div className={css.header}>
        <h1 className={css.title}>My Courses</h1>
        <p className={css.subtitle}>Manage your enrolled courses</p>
      </div>

      {isLoading && <Loader label="Loading courses…" />}
      {isError && <ErrorMessage onRetry={() => refetch()} />}
      {data && cards.length === 0 && (
        <p className={css.state}>You are not enrolled in any courses yet.</p>
      )}

      {cards.length > 0 && (
        <ul className={css.list}>
          {cards.map((c) => (
            <CourseCard key={c.id} course={c} />
          ))}
        </ul>
      )}
    </main>
  )
}

export default function App() {
  return <Courses />
}