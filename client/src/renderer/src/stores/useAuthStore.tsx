import { create } from 'zustand'

interface IAuthStore {
  userInfo: {
    id: string
    userName: string
    aliasName: string
    avatar: string
  } | null
  authToken: string
  setAuthToken: (authToken: string) => void
  setUserInfo: (
    userInfo: { avatar: string; id: string; userName: string; aliasName: string } | null
  ) => void
  signingIn: boolean
  setSigningIn: (signingIn: boolean) => void
  isAuth: boolean
  setIsAuth: (isAuth: boolean) => void
}

const useAuthStore = create<IAuthStore>((set) => ({
  authToken: '',
  setAuthToken: (authToken: string) => set({ authToken }),
  userInfo: null,
  setUserInfo: (
    userInfo: { avatar: string; userName: string; id: string; aliasName: string } | null
  ) => set({ userInfo }),
  signingIn: false,
  setSigningIn: (signingIn: boolean) => set({ signingIn }),
  isAuth: false,
  setIsAuth: (isAuth: boolean) => set({ isAuth })
}))

export default useAuthStore
