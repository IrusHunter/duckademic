import css from './Avatar.module.css'

const PALETTE = [
  '#DB4437', '#E67E22', '#F4B400', '#0F9D58', '#16A085',
  '#4285F4', '#2980B9', '#8E44AD', '#9B59B6', '#E91E63',
  '#FF5722', '#795548', '#607D8B', '#009688',
]

function hashString(str: string): number {
  let hash = 0
  for (let i = 0; i < str.length; i++) {
    hash = str.charCodeAt(i) + ((hash << 5) - hash)
    hash |= 0
  }
  return Math.abs(hash)
}

function colorFor(value: string): string {
  return PALETTE[hashString(value || '?') % PALETTE.length]
}

function initialsFor(name: string): string {
  const base = (name || '').split('@')[0].trim()
  const parts = base.split(/[\s._-]+/).filter(Boolean)
  if (parts.length >= 2) return (parts[0][0] + parts[1][0]).toUpperCase()
  return (base.slice(0, 1) || '?').toUpperCase()
}

interface AvatarProps {
  name: string
  size?: number
  className?: string
}

export default function Avatar({ name, size = 48, className }: AvatarProps) {
  return (
    <div
      className={`${css.avatar} ${className ?? ''}`}
      style={{
        width: size,
        height: size,
        backgroundColor: colorFor(name),
        fontSize: Math.round(size * 0.4),
      }}
      aria-hidden="true"
      title={name}
    >
      {initialsFor(name)}
    </div>
  )
}
