import css from './EmptyState.module.css'

interface EmptyStateProps {
  title: string
  hint?: string
}

export default function EmptyState({ title, hint }: EmptyStateProps) {
  return (
    <div className={css.wrapper}>
      <p className={css.title}>{title}</p>
      {hint && <p className={css.hint}>{hint}</p>}
    </div>
  )
}
