import { useMemo } from 'react'
import { QueryClient, QueryClientProvider, useQuery } from '@tanstack/react-query'
import { getCoursesProgress, getUpcomingEvents } from '../../services/courseService'
import CourseCard from '../CourseCard/CourseCard'
import type { CourseCardData } from '../CourseCard/CourseCard'
import Loader from '../Loader/Loader'
import ErrorMessage from '../ErrorMessage/ErrorMessage'
import css from './App.module.css'

const queryClient = new QueryClient({
  defaultOptions: { queries: { staleTime: 60_000, retry: 1 } },
})

const COLORS = ['cardBlue', 'cardMint', 'cardRose'] as const

function Courses() {
  const startTime = useMemo(() => new Date().toISOString(), [])

  // Список моїх курсів — дозволений студенту ендпоінт
  const progress = useQuery({ queryKey: ['courses-progress'], queryFn: getCoursesProgress })
  // Мої завдання — для assignments/дедлайнів (best-effort зіставлення за курсом)
  const upcoming = useQuery({
    queryKey: ['upcoming', startTime],
    queryFn: () => getUpcomingEvents(50, startTime),
  })

  const cards: CourseCardData[] = useMemo(() => {
    const list = progress.data ?? []
    const tasks = upcoming.data ?? []
    const now = Date.now()

    return list.map((c, i) => {
      // Зіставляємо завдання з курсом за course_id (best-effort)
      const courseTasks = tasks.filter((t) => t.course_id === c.id)
      const hasTasks = courseTasks.length > 0

      const upcomingDeadlines = courseTasks
        .map((t) => t.deadline)
        .filter((d) => new Date(d).getTime() >= now)
        .sort((a, b) => new Date(a).getTime() - new Date(b).getTime())

      return {
        id: c.id,
        title: c.name,
        ratePercent: Math.round((c.complete_rate ?? 0) * 100),
        accuracyPercent: Math.round((c.complete_accuracy ?? 0) * 100),
        assignments: hasTasks ? courseTasks.length : undefined,
        nearestDeadline: upcomingDeadlines[0],
        colorClass: COLORS[i % COLORS.length],
      }
    })
  }, [progress.data, upcoming.data])

  return (
    <main className={css.page}>
      <div className={css.header}>
        <h1 className={css.title}>My Courses</h1>
        <p className={css.subtitle}>Manage your enrolled courses</p>
      </div>

      {progress.isLoading && <Loader label="Loading courses…" />}
      {progress.isError && <ErrorMessage onRetry={() => progress.refetch()} />}
      {progress.data && cards.length === 0 && (
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
  return (
    <QueryClientProvider client={queryClient}>
      <Courses />
    </QueryClientProvider>
  )
}
