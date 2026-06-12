import { api } from './api'
import type { LessonOccurrence } from '../types/schedule'

interface PersonalScheduleParams {
  startTime: string // ISO — початок періоду
  endTime: string   // ISO — кінець періоду
}

// POST /api/schedule/get-personal-schedule
// Бекенд визначає користувача за токеном і повертає його заняття за період.
// body: { start_time, end_time }  →  LessonOccurrence[] (повна модель)
export async function getPersonalSchedule(
  { startTime, endTime }: PersonalScheduleParams,
): Promise<LessonOccurrence[]> {
  const res = await api.post<LessonOccurrence[]>(
    '/api/schedule/get-personal-schedule',
    { start_time: startTime, end_time: endTime },
  )
  return res.data ?? []
}
