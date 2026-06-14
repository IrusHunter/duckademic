import css from './Loader.module.css'
export default function Loader({ label = 'Loading…' }: { label?: string }) {
  return (
    <div className={css.wrapper} role="status">
      <span className={css.spinner} aria-hidden="true" />
      <span>{label}</span>
    </div>
  )
}
