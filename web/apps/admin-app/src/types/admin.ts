export type FieldDef = {
  key: string
  label: string
  required?: boolean
  format?: 'time-ns' | 'duration-ns' | 'weekday-ua'
  relation?: {
    serviceKey: string
    tableKey: string
    labelKey: string
  }
}

export type TableDef = {
  key: string
  label: string
  listEndpoint: string
  itemEndpoint: string
  fields: FieldDef[]
  editFields?: FieldDef[]
  columns: FieldDef[]
  readOnly?: boolean
  numericKeys?: string[]
}

export type ServiceDef = {
  key: string
  label: string
  icon: React.ReactNode | string
  baseURL: string
  tables: TableDef[]
}