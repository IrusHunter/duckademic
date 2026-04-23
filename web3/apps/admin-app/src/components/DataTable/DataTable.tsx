import { useState } from 'react'
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
      <div style={{ overflowX: 'auto' }}>
        <table style={{ width: '100%', borderCollapse: 'collapse', fontSize: 14 }}>
          <thead>
            <tr className={css.tr}>
              {canDelete && <th style={{ padding: '8px 12px', borderBottom: '2px solid #e0e0e0', width: 40 }} />}
              {columns.map(col => (
                <th key={col.key} style={{ padding: '8px 12px', textAlign: 'left', borderBottom: '2px solid #e0e0e0', color: '#555', whiteSpace: 'nowrap' }}>
                  {col.label}
                </th>
              ))}
              {showActions && <th style={{ padding: '8px 12px', borderBottom: '2px solid #e0e0e0', color: '#555' }}>Actions</th>}
            </tr>
          </thead>
          <tbody>
            {data.length === 0 && (
              <tr>
                <td colSpan={columns.length + (showActions ? 1 : 0) + (canDelete ? 1 : 0)} style={{ padding: 24, textAlign: 'center', color: '#aaa' }}>
                  No data
                </td>
              </tr>
            )}
            {data.map((row, idx) => {
              const id = String(row.id ?? '')
              const rowKey = id || `row-${idx}`
              return (
                <tr key={rowKey} style={{ background: selected.includes(id) ? '#e8f0fe' : 'white' }}>
                  {canDelete && (
                    <td style={{ padding: '8px 12px', borderBottom: '1px solid #eee' }}>
                      <input type="checkbox" checked={selected.includes(id)} onChange={() => toggleSelect(id)} />
                    </td>
                  )}
                  {columns.map(col => (
                    <td key={col.key} style={{ padding: '8px 12px', borderBottom: '1px solid #eee', maxWidth: 220, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
                      {formatCell(row[col.key], col.format)}
                    </td>
                  ))}
                  {showActions && (
                    <td style={{ padding: '8px 12px', borderBottom: '1px solid #eee', whiteSpace: 'nowrap' }}>
                      {editFields.length > 0 && (
                        <button
                          onClick={() => onEditClick(id)}
                          style={{ padding: '3px 10px', color: '#4A6CF7', border: '1px solid #4A6CF7', borderRadius: 4, background: 'white', cursor: 'pointer', fontSize: 13, marginRight: 6 }}
                        >
                          Edit
                        </button>
                      )}
                      {canDelete && (
                        <button
                          onClick={() => onDelete(id)}
                          style={{ padding: '3px 10px', color: 'red', border: '1px solid red', borderRadius: 4, background: 'white', cursor: 'pointer', fontSize: 13 }}
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