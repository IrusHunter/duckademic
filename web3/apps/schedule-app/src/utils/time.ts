// Converts nanoseconds-since-midnight to "HH:MM"
export function nsToTime(ns: number): string {
  const totalSec = Math.floor(ns / 1_000_000_000)
  const h = Math.floor(totalSec / 3600)
  const m = Math.floor((totalSec % 3600) / 60)
  return `${String(h).padStart(2, '0')}:${String(m).padStart(2, '0')}`
}
