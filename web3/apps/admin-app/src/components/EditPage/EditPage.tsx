import { useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { LiaLongArrowAltLeftSolid } from 'react-icons/lia'
import { makeApi } from '../../api/makeApi'
import { formatCell } from '../../utils/formatters'
import { RelationSelect } from '../RelationSelect/RelationSelect'
import type { ServiceDef, TableDef } from '../../types/admin'
import css from './EditPage.module.css'
import appCss from '../App/App.module.css'

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
      <button onClick={() => navigate(-1)} className={appCss.buttonBack}>
        <LiaLongArrowAltLeftSolid size={20} />Back
      </button>
      <h2 className={appCss.itemTitle}>Edit — {table.label}</h2>
      <p className={css.subtitle}>
        {service.label} · {service.baseURL}{table.itemEndpoint}/{itemId}
      </p>

      {isLoading && <p>Loading...</p>}

      {!isLoading && raw && (
        <>
          <div className={css.recordInfo}>
            <p className={css.recordInfoTitle}>Record info</p>
            {table.columns
              .filter(col => !editFields.find(f => f.key === col.key))
              .map(col => (
                <div key={col.key} className={css.recordInfoRow}>
                  <span className={css.recordInfoLabel}>{col.label}:</span>
                  <span className={css.recordInfoValue}>{formatCell(raw[col.key], col.format)}</span>
                </div>
              ))
            }
          </div>

          <form onSubmit={handleSubmit}>
            <div className={css.fieldsWrapper}>
              {editFields.map(field => (
                <div key={field.key}>
                  <label className={css.fieldLabel}>
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
                      className={css.fieldInput}
                      value={values[field.key] ?? ''}
                      onChange={e => setValues(prev => ({ ...prev, [field.key]: e.target.value }))}
                      required={field.required}
                    />
                  )}
                </div>
              ))}
            </div>

            <div className={css.actions}>
              <button
                type="submit"
                disabled={updateMutation.isPending}
                className={css.buttonSave}
              >
                {updateMutation.isPending ? 'Saving...' : '✓ Save'}
              </button>
              <button
                type="button"
                onClick={() => navigate(-1)}
                className={css.buttonCancel}
              >
                Cancel
              </button>
            </div>

            {updateMutation.isError && (
              <p className={css.errorMessage}>
                Error: {(updateMutation.error as any)?.response?.data?.error || 'Unknown error'}
              </p>
            )}
          </form>
        </>
      )}
    </div>
  )
}