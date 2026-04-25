import { getUserFromCookie } from './cookies'
import { useAuthStore } from '../store/authStore'
import { tokenManager } from '../auth/tokenManager'

// Відновлення сесії після перезагрузки сторінки:
// 1. Читаємо auth_user з куки → одразу ставимо в стор (Header рендериться без мигання)
// 2. Викликаємо initRefresh → сервер перевіряє refresh_token з куки і видає новий access_token
//    - OK  → accessToken отримано, сесія активна
//    - null → refresh_token протух або відсутній → clearUser() → юзер іде на /login
// 3. setInitialized() у finally → завжди розблоковує ProtectedRoute
export const initAuth = async (): Promise<void> => {
  try {
    const user = getUserFromCookie() // тепер читає з localStorage
    if (!user) return

    useAuthStore.getState().setUser(user)

    const token = await tokenManager.initRefresh()
    if (!token) {
      useAuthStore.getState().clearUser()
    }
  } finally {
    useAuthStore.getState().setInitialized()
  }
}