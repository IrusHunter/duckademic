import axios from 'axios'
import type { ServiceDef } from '../types/admin'
import { SERVICES } from '../config/services'

export function makeApi(baseURL: string) {
  return axios.create({ baseURL, withCredentials: true })
}

export function getServiceByKey(key: string): ServiceDef | undefined {
  return SERVICES.find(s => s.key === key)
}

export function normalizeArray(result: unknown): Record<string, unknown>[] {
  if (Array.isArray(result)) return result
  if (result && typeof result === 'object') {
    const obj = result as Record<string, unknown>
    if (Array.isArray(obj.data)) return obj.data as Record<string, unknown>[]
    if (Array.isArray(obj.items)) return obj.items as Record<string, unknown>[]
    if (Array.isArray(obj.results)) return obj.results as Record<string, unknown>[]
    return [obj]
  }
  return []
}