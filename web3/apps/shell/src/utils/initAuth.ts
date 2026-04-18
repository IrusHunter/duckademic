import { getUserFromCookie } from './cookies'
import { useAuthStore } from '../store/authStore'

export const initAuth = () => {
  const user = getUserFromCookie()
  if (user) {
    useAuthStore.getState().setUser(user)
  }
}