import { useMemo } from 'react'
import { Routes, Route, useNavigate } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { getStudentCourses, getTeacherCourses } from '../../services/courseService'
import { getAuthUser } from '../../utils/user'
import CourseCard from '../CourseCard/CourseCard'
import type { CourseCardData } from '../CourseCard/CourseCard'
import Loader from '../Loader/Loader'
import ErrorMessage from '../ErrorMessage/ErrorMessage'
import CoursePage from '../CoursePage/CoursePage'
import css from './App.module.css'
import type { CourseInfo } from '../../types/course'

const COLORS = ['cardBlue', 'cardMint', 'cardRose'] as const

function StudentCourses() {
  const navigate = useNavigate()
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

  const courseInfoMap = useMemo(() => {
    const m = new Map<string, CourseInfo>()
    ;(data ?? []).forEach((c, i) => {
      m.set(c.id, { id: c.id, name: c.name, description: c.description, teacher_name: c.teacher_name, colorIndex: i })
    })
    return m
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
            <CourseCard
              key={c.id}
              course={c}
              onOpen={() => navigate(c.id, { state: courseInfoMap.get(c.id) })}
            />
          ))}
        </ul>
      )}
    </main>
  )
}

function TeacherCourses() {
  const navigate = useNavigate()
  const { data, isLoading, isError, refetch } = useQuery({
    queryKey: ['teacher-courses'],
    queryFn: getTeacherCourses,
  })

  const cards: CourseCardData[] = useMemo(() => {
    return (data ?? []).map((c, i) => ({
      id: c.id,
      title: c.name,
      teacher: '',
      description: c.description,
      assignments: c.assignments_count,
      students: c.student_count,
      averageMark: 0,
      nearestDeadline: undefined,
      colorClass: COLORS[i % COLORS.length],
      hideGrade: true,
      hideDeadline: true,
    }))
  }, [data])

  const courseInfoMap = useMemo(() => {
    const m = new Map<string, CourseInfo>()
    ;(data ?? []).forEach((c, i) => {
      m.set(c.id, { id: c.id, name: c.name, description: c.description, colorIndex: i })
    })
    return m
  }, [data])

  return (
    <main className={css.page}>
      <div className={css.header}>
        <h1 className={css.title}>My Courses</h1>
        <p className={css.subtitle}>Courses you teach</p>
      </div>

      {isLoading && <Loader label="Loading courses…" />}
      {isError && <ErrorMessage onRetry={() => refetch()} />}
      {data && cards.length === 0 && (
        <p className={css.state}>You are not assigned to any courses yet.</p>
      )}

      {cards.length > 0 && (
        <ul className={css.list}>
          {cards.map((c) => (
            <CourseCard
              key={c.id}
              course={c}
              onOpen={() => navigate(c.id, { state: courseInfoMap.get(c.id) })}
            />
          ))}
        </ul>
      )}
    </main>
  )
}

function CourseList() {
  const user = getAuthUser()
  return user?.role === 'teacher' ? <TeacherCourses /> : <StudentCourses />
}

export default function App() {
  return (
    <Routes>
      <Route index element={<CourseList />} />
      <Route path=":courseId" element={<CoursePage />} />
    </Routes>
  )
}
