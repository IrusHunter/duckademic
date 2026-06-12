import css from './Loader.module.css'

interface LoaderProps {
  label?: string
}

export default function Loader({ label = 'Завантаження…' }: LoaderProps) {
  return (
    <div className={css.wrapper} role="status" aria-live="polite">
      <span className={css.spinner} aria-hidden="true" />
      <span className={css.text}>{label}</span>
    </div>
  )
}
