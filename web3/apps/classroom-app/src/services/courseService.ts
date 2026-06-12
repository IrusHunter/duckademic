import { api, asArray } from './api'
import type { CourseProgress, Task } from '../types/course'

// Список МОЇХ курсів (дозволено студенту)
export async function getCoursesProgress(): Promise<CourseProgress[]> {
  const res = await api.get('/api/course/get-courses-progress')
  return asArray<CourseProgress>(res.data)
}

// МОЇ найближчі завдання (для підрахунку assignments + дедлайнів на курс)
export async function getUpcomingEvents(count: number, startTime: string): Promise<Task[]> {
  const res = await api.get('/api/course/get-upcoming-events', {
    params: { count, 'start-time': startTime },
  })
  return asArray<Task>(res.data)
}
