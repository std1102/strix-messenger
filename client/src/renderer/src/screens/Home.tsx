import clsx from 'clsx'
import Sidebar from '../components/Sidebar'
import Chat from '../components/Chat'
import useCallStore from '../stores/useCallStore'
import VideoCall from '../components/VideoCall'
import useConversationStore from '../stores/useConversationStore'
import chatRepository from '../repositories/chat-repository'
import { useCallback, useEffect } from 'react'
import IConversation from '../interfaces/IConversation'
import useAuthStore from '../stores/useAuthStore'
import ReceivingCallModal from '../components/ReceivingCallModal'

const HomeScreen = () => {
  const { status } = useCallStore()
  const { currentConversation, setConversations } = useConversationStore()
  const { userInfo } = useAuthStore()

  const initRatchet = async () => {
    try {
      const ratchetList = await window.api.getRatchetDetailList(userInfo!.userName)
      for (const ratchet of ratchetList) {
        await window.loadRatchet(JSON.stringify(ratchet))
      }
    } catch (error) {
      console.error('ERROR', error)
    }
  }

  const getChatSession = useCallback(async () => {
    try {
      // init chat
      const oldChatSessions = await window.api.getOldChatSessions(userInfo!.userName)
      console.log('oldChatSessions', oldChatSessions)
      const newConversations: IConversation[] = oldChatSessions
        .filter((item) => item.messages.length > 0)
        .map<IConversation>((item) => ({
          lastMessage: item.lastMessage,
          receiver: item.receiver,
          id: item.ratchetId,
          updatedAt: item.updatedAt,
          isReaded: item.isReaded
        }))

      // get pending chat session
      const res = await chatRepository.getPendingChatSession()

      if (res) {
        for (const item of res.data) {
          const ratchetRes = await window.initRatchetFromExternal(
            JSON.stringify(item.senderKeyBundle),
            item.ephemeralKey,
            item.chatSessionId
          )
          await chatRepository.completeChatSession(ratchetRes.ratchetId)
          const ratchetDetail = await window.saveRatchet(ratchetRes.ratchetId)
          await window.api.createRatchetFile(
            item.senderUserName,
            ratchetDetail,
            ratchetRes.ratchetId
          )

          newConversations.push({
            lastMessage: 'Đang chờ phản hồi',
            receiver: item.senderUserName,
            id: item.chatSessionId,
            updatedAt: item.updatedAt || new Date().toISOString(),
            isReaded: false
          })
        }
        setConversations(newConversations)
      }
    } catch (error) {
      console.error('ERROR', error)
    }
  }, [userInfo])

  useEffect(() => {
    Promise.all([initRatchet(), getChatSession()]).then()
  }, [getChatSession])

  return (
    <div className="flex h-screen w-screen relative">
      <div
        className={clsx(
          'bg-yellow-50 w-0 duration-700 relative overflow-hidden',
          status === 'on-call' && 'w-3/4'
        )}
      >
        <VideoCall />
      </div>

      <div
        className={clsx(
          'w-1/4 h-full border-r overflow-hidden',
          status === 'on-call' && 'w-0 hidden'
        )}
      >
        <Sidebar />
      </div>

      <div
        className={clsx(
          'bg-white flex flex-col duration-700 overflow-hidden relative',
          status === 'on-call' ? 'w-1/4' : 'w-3/4'
        )}
      >
        <ReceivingCallModal />

        {currentConversation ? (
          <Chat />
        ) : (
          <div className="absolute top-1/2 left-1/2 -translate-y-1/2 -translate-x-1/2 text-center">
            <h1 className="text-2xl font-bold">Welcome to Chat App</h1>
            <p className="text-gray-500 mt-2">
              Select a conversation from the sidebar or start a new one.
            </p>
          </div>
        )}
      </div>
    </div>
  )
}

export default HomeScreen
