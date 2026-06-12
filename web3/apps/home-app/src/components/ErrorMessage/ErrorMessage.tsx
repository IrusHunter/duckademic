import css from './ErrorMessage.module.css'

interface ErrorMessageProps {
  message?: string
  onRetry?: () => void
}

export default function ErrorMessage({
  message = 'Не вдалося завантажити дані. Спробуйте ще раз.',
  onRetry,
}: ErrorMessageProps) {
  return (
    <div className={css.wrapper} role="alert">
      <p className={css.text}>{message}</p>
      {onRetry && (
        <button type="button" className={css.retry} onClick={onRetry}>
          Повторити
        </button>
      )}
    </div>
  )
}
