import { api, asArray } from './api'
import type { LessonOccurrence } from '../types/schedule'

function weekBounds() {
  const now = new Date()
  const day = now.getUTCDay()
  const daysFromMon = day === 0 ? 6 : day - 1

  const mon = new Date(now)
  mon.setUTCDate(now.getUTCDate() - daysFromMon)
  mon.setUTCHours(0, 0, 0, 0)

  const sun = new Date(mon)
  sun.setUTCDate(mon.getUTCDate() + 6)
  sun.setUTCHours(23, 59, 59, 999)

  return { start: mon, end: sun }
}

export async function getPersonalSchedule(): Promise<LessonOccurrence[]> {
  const { start, end } = weekBounds()
  const res = await api.post('/api/schedule/get-personal-schedule', {
    start_time: start.toISOString(),
    end_time: end.toISOString(),
  })
  return asArray<LessonOccurrence>(res.data)
}

export async function getAllPersonalSchedule(): Promise<LessonOccurrence[]> {
  const res = await api.get('/api/schedule/get-all-personal-schedule')
  return asArray<LessonOccurrence>(res.data)
}
