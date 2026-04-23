import { create } from 'zustand'

export type User = {
  id: string
  email: string        // зберігаємо login як email для сумісності з хедером
  role: 'admin' | 'student' | 'teacher'
  is_default_password?: boolean
}

type AuthStore = {
  user: User | null
  isAuthenticated: boolean
  setUser: (user: User) => void
  clearUser: () => void
}

export const useAuthStore = create<AuthStore>((set) => ({
  user: null,
  isAuthenticated: false,
  setUser: (user) => set({ user, isAuthenticated: true }),
  clearUser: () => set({ user: null, isAuthenticated: false }),
}))