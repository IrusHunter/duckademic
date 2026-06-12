import css from './ErrorMessage.module.css'
interface Props { message?: string; onRetry?: () => void }
export default function ErrorMessage({ message = 'Failed to load courses.', onRetry }: Props) {
  return (
    <div className={css.wrapper} role="alert">
      <span>{message}</span>
      {onRetry && <button type="button" className={css.retry} onClick={onRetry}>Retry</button>}
    </div>
  )
}
