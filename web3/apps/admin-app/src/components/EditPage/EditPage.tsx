import { useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { LiaLongArrowAltLeftSolid } from 'react-icons/lia'
import { makeApi } from '../../api/makeApi'
import { formatCell } from '../../utils/formatters'
import { RelationSelect } from '../RelationSelect/RelationSelect'
import type { ServiceDef, TableDef } from '../../types/admin'
import css from '../App/App.module.css'

type Props = {
  service: ServiceDef
  table: TableDef
}

export function EditPage({ service, table }: Props) {
  const navigate = useNavigate()
  const { itemId } = useParams<{ itemId: string }>()
  const qc = useQueryClient()
  const api = makeApi(service.baseURL)
  const queryKey = [service.key, table.key]

  const editFields = table.editFields ?? table.fields

  const { data: raw, isLoading } = useQuery({
    queryKey: [service.key, table.key, itemId],
    queryFn: async () => {
      const r = await api.get(`${table.itemEndpoint}/${itemId}`)
      return r.data as Record<string, unknown>
    },
    enabled: !!itemId,
  })

  const [values, setValues] = useState<Record<string, string>>({})
  const initialized = Object.keys(values).length > 0
  if (raw && !initialized) {
    const initial = Object.fromEntries(editFields.map(f => [f.key, String(raw[f.key] ?? '')]))
    setValues(initial)
  }

  const getErrorMessage = (e: any): string => e?.response?.data?.error || 'Unknown error'

  const updateMutation = useMutation({
    mutationFn: (body: Record<string, string>) => {
      const converted: Record<string, unknown> = { ...body }
      for (const key of table.numericKeys ?? []) {
        if (converted[key] !== undefined && converted[key] !== '') {
          converted[key] = Number(converted[key])
        }
      }
      return api.put(`${table.itemEndpoint}/${itemId}`, converted).then(r => r.data)
    },
    onSuccess: () => {
      qc.invalidateQueries({ queryKey })
      navigate(-1)
    },
    onError: (e) => console.error('Update error:', getErrorMessage(e)),
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    updateMutation.mutate(values)
  }

  return (
    <div className={css.wrapper}>
      <button onClick={() => navigate(-1)} className={css.buttonBack}>
        <LiaLongArrowAltLeftSolid size={20} />Back
      </button>
      <h2 className={css.itemTitle}>Edit — {table.label}</h2>
      <p style={{ color: '#999', fontSize: 13, marginBottom: 24 }}>
        {service.label} · {service.baseURL}{table.itemEndpoint}/{itemId}
      </p>

      {isLoading && <p>Loading...</p>}

      {!isLoading && raw && (
        <>
          <div style={{ marginBottom: 24, padding: 16, background: '#f9f9f9', borderRadius: 8, fontSize: 13 }}>
            <p style={{ fontWeight: 600, marginBottom: 8, color: '#555' }}>Record info</p>
            {table.columns
              .filter(col => !editFields.find(f => f.key === col.key))
              .map(col => (
                <div key={col.key} style={{ display: 'flex', gap: 8, marginBottom: 4 }}>
                  <span style={{ color: '#999', minWidth: 140 }}>{col.label}:</span>
                  <span style={{ color: '#333' }}>{formatCell(raw[col.key], col.format)}</span>
                </div>
              ))
            }
          </div>

          <form onSubmit={handleSubmit}>
            <div style={{ display: 'flex', flexDirection: 'column', gap: 16, marginBottom: 24 }}>
              {editFields.map(field => (
                <div key={field.key}>
                  <label style={{ display: 'block', fontSize: 13, fontWeight: 600, color: '#444', marginBottom: 6 }}>
                    {field.label}{field.required && ' *'}
                  </label>
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
                      style={{ padding: '8px 12px', borderRadius: 4, border: '1px solid #ccc', width: '100%', boxSizing: 'border-box', fontSize: 14 }}
                    />
                  )}
                </div>
              ))}
            </div>

            <div style={{ display: 'flex', gap: 10 }}>
              <button
                type="submit"
                disabled={updateMutation.isPending}
                style={{ padding: '8px 24px', background: '#4A6CF7', color: 'white', border: 'none', borderRadius: 4, cursor: 'pointer', fontSize: 14 }}
              >
                {updateMutation.isPending ? 'Saving...' : '✓ Save'}
              </button>
              <button
                type="button"
                onClick={() => navigate(-1)}
                style={{ padding: '8px 24px', background: 'white', color: '#666', border: '1px solid #ccc', borderRadius: 4, cursor: 'pointer', fontSize: 14 }}
              >
                Cancel
              </button>
            </div>

            {updateMutation.isError && (
              <p style={{ color: 'red', marginTop: 12, fontSize: 13 }}>
                Error: {(updateMutation.error as any)?.response?.data?.error || 'Unknown error'}
              </p>
            )}
          </form>
        </>
      )}
    </div>
  )
}