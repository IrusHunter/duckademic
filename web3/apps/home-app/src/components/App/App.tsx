import { useMemo } from 'react'
import { QueryClient, QueryClientProvider, useQuery } from '@tanstack/react-query'
import { getUpcomingEvents, getCoursesProgress } from '../../services/courseService'
import { getAuthUser, roleLabel } from '../../utils/user'
import { nowISO } from '../../utils/datetime'
import ProfileCard from '../ProfileCard/ProfileCard'
import QuickActions from '../QuickActions/QuickActions'
import Feed from '../Feed/Feed'
import Upcoming from '../Upcoming/Upcoming'
import CourseProgress from '../CourseProgress/CourseProgress'
import css from './App.module.css'

const queryClient = new QueryClient({
  defaultOptions: { queries: { staleTime: 60_000, retry: 1 } },
})

function Dashboard() {
  const user = useMemo(() => getAuthUser(), [])
  const startTime = useMemo(() => nowISO(), [])

  const upcoming = useQuery({
    queryKey: ['upcoming', startTime],
    queryFn: () => getUpcomingEvents(5, startTime),
  })

  const progress = useQuery({
    queryKey: ['courses-progress'],
    queryFn: () => getCoursesProgress(),
  })

  const tasks = upcoming.data ?? []
  const courses = progress.data ?? []

  // Лічильники профілю — з реальних даних (де можливо)
  const coursesCount = progress.isLoading ? '—' : courses.length
  const assignmentsCount = upcoming.isLoading ? '—' : tasks.length
  // Окремого ендпоінта "мої навчальні групи" в бекенді поки немає → показуємо прочерк
  const groupsCount = '—'

  return (
    <main className={css.dashboard}>
      {/* ── Ліва колонка ── */}
      <div className={css.left}>
        <ProfileCard
          name={user?.email ?? 'Користувач'}
          role={roleLabel(user?.role)}
          avatar="/img/profile_pic.png"
          courses={coursesCount}
          assignments={assignmentsCount}
          groups={groupsCount}
        />
        <QuickActions />
      </div>

      {/* ── Центр: стрічка (моки) ── */}
      <div className={css.center}>
        <Feed />
      </div>

      {/* ── Права колонка: реальні дані ── */}
      <div className={css.right}>
        <Upcoming tasks={tasks} />
        <CourseProgress items={courses} />
      </div>
    </main>
  )
}

export default function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <Dashboard />
    </QueryClientProvider>
  )
}
