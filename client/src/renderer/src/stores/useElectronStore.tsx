import { create } from 'zustand'
import Store from 'electron-store'

interface IElectronStore {
  store: Store
  setStore: (store: Store) => void
}

const useElectronStore = create<IElectronStore>((set) => ({
  store: new Store(),
  setStore: (store: Store) => set({ store })
}))

export default useElectronStore
