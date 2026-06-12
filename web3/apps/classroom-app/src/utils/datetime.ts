const fmt = new Intl.DateTimeFormat('en-US', {
  weekday: 'short',
  hour: '2-digit',
  minute: '2-digit',
})

export function formatDeadline(iso: string): string {
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return iso
  return fmt.format(d)
}

// Класифікація терміновості для кольору дедлайну
export type DeadlineLevel = 'warning' | 'info' | 'danger'

export function deadlineLevel(iso: string): DeadlineLevel {
  const diffDays = (new Date(iso).getTime() - Date.now()) / 86_400_000
  if (diffDays < 0) return 'danger'
  if (diffDays <= 2) return 'warning'
  return 'info'
}
