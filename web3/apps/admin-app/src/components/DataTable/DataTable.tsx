import { useState } from 'react'
import { BiCheck } from "react-icons/bi";
import type { FieldDef } from '../../types/admin'
import { formatCell } from '../../utils/formatters'
import { ActionSelect } from '../ActionSelect/ActionSelect'
import css from './DataTable.module.css'

type Props = {
  data: Record<string, unknown>[]
  columns: FieldDef[]
  editFields: FieldDef[]
  onDelete: (id: string) => void
  onEditClick: (id: string) => void
  readOnly?: boolean
  canDelete: boolean
}

export function DataTable({ data, columns, editFields, onDelete, onEditClick, readOnly, canDelete }: Props) {
  const [selected, setSelected] = useState<string[]>([])
  const [action, setAction] = useState('')

  const toggleSelect = (id: string) =>
    setSelected(prev => prev.includes(id) ? prev.filter(i => i !== id) : [...prev, id])

  const handleGo = () => {
    if (action === 'delete') {
      selected.forEach(id => onDelete(id))
      setSelected([])
    }
  }

  const showActions = !readOnly && (editFields.length > 0 || canDelete)

  const ACTION_OPTIONS = [
    { value: '', label: '---' },
    { value: 'delete', label: 'Delete selected' },
  ]

  return (
    <div>
      {canDelete && (
        <ActionSelect
          value={action}
          options={ACTION_OPTIONS}
          onChange={setAction}
          onGo={handleGo}
        />
      )}
      <div className={css.wrapper}>
        <table className={css.table}>
          <thead>
            <tr className={css.tr}>
              {canDelete && <th className={css.thCheckbox} />}
              {columns.map(col => (
                <th key={col.key} className={css.th}>
                  {col.label}
                </th>
              ))}
              {showActions && <th className={css.th}>Actions</th>}
            </tr>
          </thead>
          <tbody>
            {data.length === 0 && (
              <tr>
                <td
                  colSpan={columns.length + (showActions ? 1 : 0) + (canDelete ? 1 : 0)}
                  className={css.tdEmpty}
                >
                  No data
                </td>
              </tr>
            )}
            {data.map((row, idx) => {
              const id = String(row.id ?? '')
              const rowKey = id || `row-${idx}`
              return (
                <tr
                  key={rowKey}
                  className={selected.includes(id) ? css.rowSelected : css.rowDefault}
                  onClick={() => toggleSelect(id)}
                >
                  {canDelete && (
                    <td className={css.tdCheckbox}>
                      <div className={`${css.checkbox} ${selected.includes(id) ? css.checkboxChecked : ''}`}>
                        {selected.includes(id) && <BiCheck className={css.checkboxIcon} />}
                      </div>
                    </td>
                  )}
                  {columns.map(col => (
                    <td key={col.key} className={css.td}>
                      {formatCell(row[col.key], col.format)}
                    </td>
                  ))}
                  {showActions && (
                    <td className={css.tdActions}>
                      {editFields.length > 0 && (
                        <button
                          onClick={(e) => { e.stopPropagation(); onEditClick(id) }}
                          className={css.btnEdit}
                        >
                          Edit
                        </button>
                      )}
                      {canDelete && (
                        <button
                          onClick={(e) => { e.stopPropagation(); onDelete(id) }}
                          className={css.btnDelete}
                        >
                          Delete
                        </button>
                      )}
                    </td>
                  )}
                </tr>
              )
            })}
          </tbody>
        </table>
      </div>
    </div>
  )
}