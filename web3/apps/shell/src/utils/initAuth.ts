import { getUserFromCookie } from './cookies'
import { useAuthStore } from '../store/authStore'
import { tokenManager } from '../auth/tokenManager'

export const initAuth = async () => {
  const user = getUserFromCookie()
  if (user) {
    useAuthStore.getState().setUser(user)
    await tokenManager.refresh();
  }
}