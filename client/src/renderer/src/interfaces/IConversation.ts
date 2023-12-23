interface IConversation {
  id: string
  receiver: string
  lastMessage: string
  updatedAt?: string
  isReaded?: boolean
}

export default IConversation
