// Невеликі помічники для роботи з датами.
// Винесено окремо, щоб компоненти не дублювали логіку форматування.

export function startOfToday(): Date {
  const d = new Date()
  d.setHours(0, 0, 0, 0)
  return d
}

export function addDays(date: Date, days: number): Date {
  const d = new Date(date)
  d.setDate(d.getDate() + days)
  return d
}

const timeFmt = new Intl.DateTimeFormat('uk-UA', {
  hour: '2-digit',
  minute: '2-digit',
})

const dayFmt = new Intl.DateTimeFormat('uk-UA', {
  weekday: 'long',
  day: 'numeric',
  month: 'long',
})

const deadlineFmt = new Intl.DateTimeFormat('uk-UA', {
  day: 'numeric',
  month: 'short',
  hour: '2-digit',
  minute: '2-digit',
})

export const formatTime = (iso: string): string => timeFmt.format(new Date(iso))
export const formatDay = (iso: string): string => dayFmt.format(new Date(iso))
export const formatDeadline = (iso: string): string => deadlineFmt.format(new Date(iso))

// Скільки днів лишилось до дедлайну (для підсвічування "горить / сьогодні").
export function daysUntil(iso: string): number {
  const ms = new Date(iso).getTime() - Date.now()
  return Math.ceil(ms / (1000 * 60 * 60 * 24))
}

// Групує елементи за календарним днем. Ключ — локальна дата (toDateString),
// значення — елементи цього дня. Зберігає порядок появи днів.
export function groupByDay<T>(items: T[], getIso: (item: T) => string): [string, T[]][] {
  const map = new Map<string, T[]>()
  for (const item of items) {
    const key = getIso(item)
    const bucket = map.get(new Date(key).toDateString())
    if (bucket) {
      bucket.push(item)
    } else {
      map.set(new Date(key).toDateString(), [item])
    }
  }
  return [...map.entries()]
}
