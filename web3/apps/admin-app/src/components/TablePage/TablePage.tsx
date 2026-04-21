import { useNavigate } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { LiaLongArrowAltLeftSolid } from 'react-icons/lia'
import { makeApi, normalizeArray } from '../../api/makeApi'
import { AddForm } from '../AddForm/AddForm'
import { DataTable } from '../DataTable/DataTable'
import type { ServiceDef, TableDef } from '../../types/admin'
import css from '../App/App.module.css'

type Props = {
  service: ServiceDef
  table: TableDef
}

export function TablePage({ service, table }: Props) {
  const navigate = useNavigate()
  const qc = useQueryClient()
  const api = makeApi(service.baseURL)
  const queryKey = [service.key, table.key]

  const { data: rawData, isLoading, error } = useQuery({
    queryKey,
    queryFn: async () => {
      const r = await api.get(table.listEndpoint)
      return normalizeArray(r.data)
    },
  })
  const data = rawData ?? []

  const getErrorMessage = (e: any): string => e?.response?.data?.error || 'Unknown error'

  const convertNumeric = (body: Record<string, string>): Record<string, unknown> => {
    const converted: Record<string, unknown> = { ...body }
    for (const key of table.numericKeys ?? []) {
      if (converted[key] !== undefined && converted[key] !== '') {
        converted[key] = Number(converted[key])
      }
    }
    return converted
  }

  const deleteMutation = useMutation({
    mutationFn: (id: string) => api.delete(`${table.itemEndpoint}/${id}`).then(r => r.data),
    onSuccess: () => qc.invalidateQueries({ queryKey }),
    onError: (e) => console.error('Delete error:', getErrorMessage(e)),
  })

  const createMutation = useMutation({
    mutationFn: (body: Record<string, string>) =>
      api.post(table.listEndpoint, convertNumeric(body)).then(r => r.data),
    onSuccess: () => qc.invalidateQueries({ queryKey }),
    onError: (e) => console.error('Create error:', getErrorMessage(e)),
  })

  const editFields = table.editFields ?? table.fields
  const canDelete = !table.readOnly && table.fields.length > 0

  return (
    <div style={{ padding: 24, flex: 1 }}>
      <button onClick={() => navigate(`/admin/${service.key}`)} className={css.buttonBack}>
        <LiaLongArrowAltLeftSolid size={20} /> <span className={css.spanBack}>Back</span>
      </button>
      <h2 className={css.itemTitle}>{table.label}</h2>

      {!table.readOnly && table.fields.length > 0 && (
        <AddForm fields={table.fields} onSubmit={body => createMutation.mutate(body)} />
      )}

      {isLoading && <p>Loading...</p>}
      {error && <p style={{ color: 'red' }}>Error loading data</p>}
      {!isLoading && !error && (
        <DataTable
          data={data}
          columns={table.columns}
          editFields={editFields}
          onDelete={id => deleteMutation.mutate(id)}
          onEditClick={id => navigate(`/admin/${service.key}/${table.key}/edit/${id}`)}
          readOnly={table.readOnly}
          canDelete={canDelete}
        />
      )}
    </div>
  )
}