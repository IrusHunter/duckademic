// Розширення типів axios для shell
import 'axios'

declare module 'axios' {
  interface InternalAxiosRequestConfig {
    _retry?: boolean  // прапор що цей запит вже ретраївся після 401
  }
}