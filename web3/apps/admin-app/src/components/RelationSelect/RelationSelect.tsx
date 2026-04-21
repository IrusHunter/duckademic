import { useQuery } from '@tanstack/react-query'
import { makeApi, normalizeArray, getServiceByKey } from '../../api/makeApi'
import type { FieldDef } from '../../types/admin'

type Props = {
  field: FieldDef
  value: string
  onChange: (v: string) => void
}

export function RelationSelect({ field, value, onChange }: Props) {
  const rel = field.relation!
  const service = getServiceByKey(rel.serviceKey)
  const table = service?.tables.find(t => t.key === rel.tableKey)

  const { data: raw, isLoading } = useQuery({
    queryKey: ['relation', rel.serviceKey, rel.tableKey],
    queryFn: async () => {
      const api = makeApi(service!.baseURL)
      const r = await api.get(table!.listEndpoint)
      return normalizeArray(r.data)
    },
    enabled: !!service && !!table,
  })

  const options = raw ?? []

  return (
    <select
      value={value}
      onChange={e => onChange(e.target.value)}
      required={field.required}
      style={{ padding: '6px 8px', borderRadius: 4, border: '1px solid #ccc', width: 194 }}
    >
      <option value="">— select —</option>
      {isLoading && <option disabled>Loading...</option>}
      {options.map((item, idx) => {
        const id = String(item.id ?? idx)
        const label = String(item[rel.labelKey] ?? id)
        return (
          <option key={id || `opt-${idx}`} value={id}>{label}</option>
        )
      })}
    </select>
  )
}