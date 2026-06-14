import { useNavigate } from 'react-router-dom'
import type { ServiceDef } from '../../types/admin'
import css from '../App/App.module.css'

type Props = {
  service: ServiceDef
}

export function ServiceHome({ service }: Props) {
  const navigate = useNavigate()
  return (
    <div className={css.wrapper}>
      <h2 className={css.serviceItemTitle}>{service.label}</h2>
      <ul className={css.serviceItemList}>
        {service.tables.map(table => (
          <li key={table.key}>
            <button
              onClick={() => navigate(`/admin/${service.key}/${table.key}`)}
              className={css.serviceItemButton}
            >
              {table.label}
              {table.readOnly && (
                <span className={css.span}>read-only</span>
              )}
            </button>
          </li>
        ))}
      </ul>
    </div>
  )
}