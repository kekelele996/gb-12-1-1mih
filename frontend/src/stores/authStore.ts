import { create } from 'zustand'
import { login as apiLogin, register as apiRegister, getProfile } from '@/api/user'
import type { UserInfo } from '@/types/api'

interface AuthState {
  token: string
  user: UserInfo | null
  setUser: (u: UserInfo | null) => void
  login: (username: string, password: string) => Promise<void>
  register: (p: { username: string; email: string; password: string; real_name?: string; phone?: string }) => Promise<void>
  fetchProfile: () => Promise<void>
  logout: () => void
}

export const useAuthStore = create<AuthState>((set) => ({
  token: localStorage.getItem('gbstudyapply_token') || '',
  user: null,
  setUser: (u) => set({ user: u }),
  login: async (username, password) => {
    const res = await apiLogin({ username, password })
    localStorage.setItem('gbstudyapply_token', res.token)
    set({ token: res.token, user: res.user })
  },
  register: async (payload) => {
    const res = await apiRegister(payload)
    localStorage.setItem('gbstudyapply_token', res.token)
    set({ token: res.token, user: res.user })
  },
  fetchProfile: async () => set({ user: await getProfile() }),
  logout: () => {
    localStorage.removeItem('gbstudyapply_token')
    set({ token: '', user: null })
  },
}))
