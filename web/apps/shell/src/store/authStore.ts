import { create } from 'zustand'

export type User = {
  id: string
  email: string
  role: 'admin' | 'student' | 'teacher'
  is_default_password?: boolean
}

type AuthStore = {
  user: User | null
  isAuthenticated: boolean
  isInitialized: boolean   // ✅ FIX #1: запобігає флікеру на /login до завершення initAuth
  setUser: (user: User) => void
  clearUser: () => void
  setInitialized: () => void
}

export const useAuthStore = create<AuthStore>((set) => ({
  user: null,
  isAuthenticated: false,
  isInitialized: false,
  setUser: (user) => set({ user, isAuthenticated: true }),
  clearUser: () => set({ user: null, isAuthenticated: false }),
  setInitialized: () => set({ isInitialized: true }),
}))