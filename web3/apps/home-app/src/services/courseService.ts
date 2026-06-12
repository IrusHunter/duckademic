import { api } from './api'
import type { Task } from '../types/course'

interface UpcomingEventsParams {
  count: number     // максимум подій
  startTime: string // ISO — від якого моменту фільтрувати
}

// GET /api/course/get-upcoming-events?count=&start-time=
// Повертає найближчі завдання/дедлайни поточного користувача.
// Увага: параметр називається саме 'start-time' (через дефіс).
export async function getUpcomingEvents(
  { count, startTime }: UpcomingEventsParams,
): Promise<Task[]> {
  const res = await api.get<Task[]>('/api/course/get-upcoming-events', {
    params: { count, 'start-time': startTime },
  })
  return res.data ?? []
}
