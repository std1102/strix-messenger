import { create } from 'zustand'
import IMessage from '../interfaces/IMessage'
import IConversation from '../interfaces/IConversation'

type IConversationStore = {
  chatAction: 'add' | 'update' | 'delete' | null
  setChatAction: (action: 'add' | 'update' | 'delete' | null) => void
  currentRatchetId: string | null
  setCurrentRatchetId: (id: string | null) => void
  conversations: IConversation[]
  setConversations: (conversations: IConversation[]) => void
  currentConversation: string | null
  setCurrentConversation: (id: string | null) => void
  currentIdxSearch: number
  setCurrentIdxSearch: (idx: number) => void
  messageSearches: string[]
  setMessageSearches: (messageSearches: string[]) => void
  messages: IMessage[]
  setMessages: (messages: IMessage[]) => void
}

const useConversationStore = create<IConversationStore>((set) => ({
  chatAction: null,
  setChatAction: (action) => set({ chatAction: action }),
  currentRatchetId: null,
  setCurrentRatchetId: (id) => set({ currentRatchetId: id }),
  conversations: [],
  setConversations: (conversations) => set({ conversations }),
  currentConversation: null,
  setCurrentConversation: (id) => set({ currentConversation: id }),
  currentIdxSearch: 0,
  setCurrentIdxSearch: (idx) => set({ currentIdxSearch: idx }),
  messageSearches: [],
  setMessageSearches: (messageSearches) => set({ messageSearches }),
  messages: [],
  setMessages: (messages) => set({ messages })
}))

export default useConversationStore
