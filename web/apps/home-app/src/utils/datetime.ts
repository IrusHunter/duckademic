const fmt = new Intl.DateTimeFormat('uk-UA', {
  weekday: 'short',
  hour: '2-digit',
  minute: '2-digit',
})

export function formatDeadline(iso: string): string {
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return iso
  return fmt.format(d)
}

export function nowISO(): string {
  return new Date().toISOString()
}
