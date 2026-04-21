import { useState } from 'react'
import type { FieldDef } from '../../types/admin'
import { RelationSelect } from '../RelationSelect/RelationSelect'

type Props = {
  fields: FieldDef[]
  onSubmit: (data: Record<string, string>) => void
}

export function AddForm({ fields, onSubmit }: Props) {
  const [values, setValues] = useState<Record<string, string>>({})
  const [open, setOpen] = useState(false)

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    onSubmit(values)
    setValues({})
    setOpen(false)
  }

  return (
    <div style={{ marginBottom: 20 }}>
      <button
        onClick={() => setOpen(!open)}
        style={{ padding: '6px 16px', background: '#4A6CF7', color: 'white', border: 'none', borderRadius: 4, cursor: 'pointer', marginBottom: 10 }}
      >
        {open ? '✕ Cancel' : '+ Add'}
      </button>
      {open && (
        <form onSubmit={handleSubmit} style={{ display: 'flex', gap: 12, flexWrap: 'wrap', padding: 16, background: '#f9f9f9', borderRadius: 8 }}>
          {fields.map(field => (
            <div key={field.key} style={{ display: 'flex', flexDirection: 'column', gap: 4 }}>
              <label style={{ fontSize: 12, color: '#666' }}>{field.label}{field.required && ' *'}</label>
              {field.relation ? (
                <RelationSelect
                  field={field}
                  value={values[field.key] ?? ''}
                  onChange={v => setValues(prev => ({ ...prev, [field.key]: v }))}
                />
              ) : (
                <input
                  value={values[field.key] ?? ''}
                  onChange={e => setValues(prev => ({ ...prev, [field.key]: e.target.value }))}
                  required={field.required}
                  style={{ padding: '6px 8px', borderRadius: 4, border: '1px solid #ccc', width: 180 }}
                />
              )}
            </div>
          ))}
          <div style={{ display: 'flex', alignItems: 'flex-end' }}>
            <button type="submit" style={{ padding: '6px 16px', background: '#333', color: 'white', border: 'none', borderRadius: 4, cursor: 'pointer' }}>
              Save
            </button>
          </div>
        </form>
      )}
    </div>
  )
}