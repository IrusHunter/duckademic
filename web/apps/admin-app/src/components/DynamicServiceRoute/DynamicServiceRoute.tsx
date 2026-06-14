import { useParams } from 'react-router-dom'
import { SERVICES } from '../../config/services'
import { ServiceHome } from '../ServiceHome/ServiceHome'
import { TablePage } from '../TablePage/TablePage'
import { EditPage } from '../EditPage/EditPage'

export function DynamicServiceRoute() {
  const { serviceKey, tableKey, itemId } = useParams<{ serviceKey: string; tableKey: string; itemId: string }>()
  const service = SERVICES.find(s => s.key === serviceKey)
  if (!service) return <div style={{ padding: 24 }}>Service not found</div>
  if (!tableKey) return <ServiceHome service={service} />
  const table = service.tables.find(t => t.key === tableKey)
  if (!table) return <div style={{ padding: 24 }}>Table not found</div>
  if (itemId) return <EditPage service={service} table={table} />
  return <TablePage service={service} table={table} />
}