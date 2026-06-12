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

  const coursesCount = progress.isLoading ? '—' : courses.length
  const assignmentsCount = upcoming.isLoading ? '—' : tasks.length
  const groupsCount = '—' // окремого ендпоінта "мої групи" в бекенді поки немає

  return (
    <main className={css.dashboard}>
      <div className={css.left}>
        <ProfileCard
          name={user?.email ?? 'User'}
          role={roleLabel(user?.role)}
          courses={coursesCount}
          assignments={assignmentsCount}
          groups={groupsCount}
        />
        <QuickActions />
      </div>

      <div className={css.center}>
        <Feed />
      </div>

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
