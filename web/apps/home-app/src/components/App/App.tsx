import { useMemo } from 'react'
import { useQuery } from '@tanstack/react-query'
import { getUpcomingEvents, getCoursesProgress, getTeacherCourses } from '../../services/courseService'
import { getAuthUser, roleLabel } from '../../utils/user'
import { nowISO } from '../../utils/datetime'
import ProfileCard from '../ProfileCard/ProfileCard'
import QuickActions from '../QuickActions/QuickActions'
import Feed from '../Feed/Feed'
import Upcoming from '../Upcoming/Upcoming'
import CourseProgress from '../CourseProgress/CourseProgress'
import css from './App.module.css'

function Dashboard() {
  const user = useMemo(() => getAuthUser(), [])
  const startTime = useMemo(() => nowISO(), [])

  const isTeacher = user?.role === 'teacher'

  const upcoming = useQuery({
    queryKey: ['upcoming', startTime],
    queryFn: () => getUpcomingEvents(5, startTime),
    enabled: !isTeacher,
  })

  const progress = useQuery({
    queryKey: ['courses-progress'],
    queryFn: () => getCoursesProgress(),
    enabled: !isTeacher,
  })

  const teacherCourses = useQuery({
    queryKey: ['teacher-courses'],
    queryFn: () => getTeacherCourses(),
    enabled: isTeacher,
  })

  const tasks = upcoming.data ?? []
  const courses = progress.data ?? []

  const coursesCount = isTeacher
    ? (teacherCourses.isLoading ? '—' : (teacherCourses.data ?? []).length)
    : (progress.isLoading ? '—' : courses.length)
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
          hideAssignments={isTeacher}
        />
        <QuickActions />
      </div>

      <div className={css.center}>
        <Feed />
      </div>

      {!isTeacher && (
        <div className={css.right}>
          <Upcoming tasks={tasks} />
          <CourseProgress items={courses} />
        </div>
      )}
    </main>
  )
}

export default function App() {
  return <Dashboard />
}
