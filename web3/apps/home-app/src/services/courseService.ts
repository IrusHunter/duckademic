import { api, asArray } from './api'
import type { CourseProgress, Task } from '../types/dashboard'

// GET /api/course/get-upcoming-events?count=&start-time=
export async function getUpcomingEvents(count: number, startTime: string): Promise<Task[]> {
  const res = await api.get('/api/course/get-upcoming-events', {
    params: { count, 'start-time': startTime },
  })
  return asArray<Task>(res.data)
}

// GET /api/course/get-courses-progress
export async function getCoursesProgress(): Promise<CourseProgress[]> {
  const res = await api.get('/api/course/get-courses-progress')
  return asArray<CourseProgress>(res.data)
}
