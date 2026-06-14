import type { FieldDef } from '../types/admin'

export function nsToTime(ns: unknown): string {
  if (ns === null || ns === undefined || ns === '') return ''
  const totalSeconds = Math.floor(Number(ns) / 1_000_000_000)
  const hours = Math.floor(totalSeconds / 3600)
  const minutes = Math.floor((totalSeconds % 3600) / 60)
  return `${String(hours).padStart(2, '0')}:${String(minutes).padStart(2, '0')}`
}

export function nsToDuration(ns: unknown): string {
  if (ns === null || ns === undefined || ns === '') return ''
  const totalMinutes = Math.floor(Number(ns) / 60_000_000_000)
  const hours = Math.floor(totalMinutes / 60)
  const minutes = totalMinutes % 60
  if (hours === 0) return `${minutes}хв`
  if (minutes === 0) return `${hours} год`
  return `${hours} год ${minutes} хв`
}

export function nsToWeekday(n: unknown): string {
  const days: Record<number, string> = {
    1: 'Понеділок',
    2: 'Вівторок',
    3: 'Середа',
    4: 'Четвер',
    5: 'Пʼятниця',
    6: 'Субота',
    7: 'Неділя',
  }
  return days[Number(n)] ?? String(n ?? '')
}

export function formatCell(value: unknown, format?: FieldDef['format']): string {
  if (format === 'time-ns') return nsToTime(value)
  if (format === 'duration-ns') return nsToDuration(value)
  if (format === 'weekday-ua') return nsToWeekday(value)
  return String(value ?? '')
}