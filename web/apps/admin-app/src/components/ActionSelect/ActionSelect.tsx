import { useEffect, useRef, useState } from 'react'
import css from './ActionSelect.module.css'
import { LuChevronDown } from "react-icons/lu";

type Option = {
  value: string
  label: string
}

type Props = {
  value: string
  options: Option[]
  onChange: (value: string) => void
  onGo: () => void
}

export function ActionSelect({ value, options, onChange, onGo }: Props) {
  const [isOpen, setIsOpen] = useState(false)
  const ref = useRef<HTMLDivElement>(null)

  const selectedLabel = options.find(o => o.value === value)?.label ?? '---'

  useEffect(() => {
    function onDocClick(e: MouseEvent) {
      if (ref.current && !ref.current.contains(e.target as Node)) {
        setIsOpen(false)
      }
    }
    function onKeyDown(e: KeyboardEvent) {
      if (e.key === 'Escape') setIsOpen(false)
    }
    document.addEventListener('click', onDocClick)
    document.addEventListener('keydown', onKeyDown)
    return () => {
      document.removeEventListener('click', onDocClick)
      document.removeEventListener('keydown', onKeyDown)
    }
  }, [])

  function select(val: string) {
    onChange(val)
    setIsOpen(false)
  }

  return (
    <div className={css.wrapper}>
      <span className={css.label}>Action:</span>

      <div ref={ref} className={`${css.dropdownWrap} ${isOpen ? css.open : ''}`}>
        <button
          type="button"
          className={css.trigger}
          aria-expanded={isOpen}
          onClick={e => { e.stopPropagation(); setIsOpen(v => !v) }}
        >
          <span className={css.triggerValue}>{selectedLabel}</span>
          <span className={css.caret}>
            <LuChevronDown size={24} />
          </span>
        </button>

        <div className={css.menu}>
          {options.map(opt => (
            <button
              key={opt.value}
              type="button"
              className={`${css.menuItem} ${opt.value === value ? css.menuItemActive : ''}`}
              onClick={() => select(opt.value)}
            >
              {opt.label}
            </button>
          ))}
        </div>
      </div>

      <button
        type="button"
        className={css.goBtn}
        onClick={onGo}
      >
        Go
      </button>
    </div>
  )
}