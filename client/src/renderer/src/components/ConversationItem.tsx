import clsx from 'clsx'
import { memo } from 'react'
import useConversationStore from '../stores/useConversationStore'
import IConversation from '../interfaces/IConversation'
import chatRepository from '../repositories/chat-repository'
import userRepository from '../repositories/user-repository'
import IMessage from '../interfaces/IMessage'

type ConversationItemProps = {
  conversation: IConversation
}

const ConversationItem = (props: ConversationItemProps) => {
  const { conversation } = props

  const {
    currentConversation,
    setConversations,
    conversations,
    setCurrentRatchetId,
    setMessages,
    setCurrentConversation,
    setMessageSearches,
    setCurrentIdxSearch
  } = useConversationStore()

  const createRatchet = async () => {
    const res = await userRepository.getExternalUserKey(conversation.receiver)
    const initRatchetRes = await window.initRatchetFromInternal(JSON.stringify(res.data))
    await chatRepository.initChatSession({
      chatSessionId: initRatchetRes.ratchetId,
      ephemeralKey: initRatchetRes.ephemeralKey,
      receiverUserName: conversation.receiver
    })
    const ratchetDetail = await window.saveRatchet(initRatchetRes.ratchetId)
    await window.api.createRatchetFile(
      conversation.receiver,
      ratchetDetail,
      initRatchetRes.ratchetId
    )
    return initRatchetRes.ratchetId
  }

  const handleSelectConversation = async () => {
    try {
      setCurrentConversation(conversation.receiver)
      setConversations(
        conversations.map((item) => {
          if (item.receiver === conversation.receiver) {
            return {
              ...item,
              isReaded: true
            }
          }
          return item
        })
      )

      const oldMessages = await window.api.getMessagesByUsername(conversation.receiver)
      const res = await chatRepository.getPendingChatMessage(conversation.id)

      let ratchetId = await window.api.getRatchetId(conversation.receiver)
      let isNewConversation = false

      if (ratchetId) {
        const isExist = await window.isRatchetExist(ratchetId)

        if (!isExist) {
          ratchetId = await createRatchet()
          isNewConversation = true
        }
      } else {
        ratchetId = await createRatchet()
        isNewConversation = true
      }

      const newMessages: IMessage[] = oldMessages.map(
        (item) =>
          ({
            content: item.content,
            index: item.index,
            sender: item.sender,
            filePath: item.filePath,
            type: item.type,
            isDeleted: item.isDeleted
          }) as IMessage
      )

      for (let i = 0; i < res?.data?.length; i++) {
        const item = res.data[i]
        const content = await window.receiveMessage(
          JSON.stringify({
            chatSessionId: item.chatSessionId,
            index: item.index,
            cipherMessage: item.cipherMessage,
            isBinary: item.isBinary
          })
        )

        console.log(content)

        const msg = {
          content: content,
          index: item.index,
          sender: item.senderUsername,
          type: item.type,
          filePath: content + ':' + item.filePath,
          isDeleted: item.isDeleted
        }

        newMessages.push(msg)

        // save new ratchet detail
        const ratchetDetail = await window.saveRatchet(item.chatSessionId!)
        await window.api.changeRatchetDetail(item.senderUsername!, ratchetDetail)
      }

      /*isNewConversation &&
        (await window.api.addMessageToRatchet(conversation.receiver, newMessages))*/

      console.log('newMessage', newMessages)

      await window.api.addMessageToRatchet(conversation.receiver, newMessages)
      await window.api.changeConversationReaded(conversation.receiver, true)

      setCurrentIdxSearch(0)
      setMessageSearches([])

      setCurrentIdxSearch(0)
      setMessageSearches([])
      setMessages(newMessages)
      setCurrentRatchetId(ratchetId)
    } catch (error) {
      console.error('ERROR', error)
    }
  }

  return (
    <div
      className={clsx(
        'border-b cursor-pointer flex gap-x-3 justify-between hover:bg-black/10 w-full p-3 items-center',
        conversation.isReaded === false && 'bg-black/10 !font-medium',
        currentConversation === conversation.receiver && 'bg-yellow-50/80'
      )}
      onClick={handleSelectConversation}
    >
      <div className="flex gap-2 items-center flex-1 ">
        <img
          src="https://source.unsplash.com/RZrIJ8C0860"
          alt="image"
          className="w-10 h-10 rounded-full"
        />
        <div className="flex flex-col overflow-hidden">
          <h4 className="font-semibold text-sm">{conversation.receiver}</h4>
          <p className="line-clamp-1 text-xs text-gray-400">{conversation.lastMessage}</p>
        </div>
      </div>

      <div className="flex flex-col items-end">
        <span className="w-2 h-2 rounded-full bg-green-500"></span>
        <span className="text-[10px] text-gray-500">
          {conversation.updatedAt && new Date(conversation.updatedAt).toLocaleTimeString()}
        </span>
      </div>
    </div>
  )
}

export default memo(ConversationItem)
