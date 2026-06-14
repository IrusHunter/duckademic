import css from './ErrorMessage.module.css'

interface Props { message?: string; onRetry?: () => void }

export default function ErrorMessage({ message = 'Не вдалося завантажити дані.', onRetry }: Props) {
  return (
    <div className={css.wrapper} role="alert">
      <span>{message}</span>
      {onRetry && <button type="button" className={css.retry} onClick={onRetry}>Повторити</button>}
    </div>
  )
}
