import { ACCEPT_CALL_EVENT, AVATAR_DEFAULT } from '../configs/consts'
import useCallStore from '../stores/useCallStore'
import { Phone, X } from 'lucide-react'
import useAuthStore from '../stores/useAuthStore'
import useWebSocketStore from '../stores/useWebSocketStore'
import IMessage from '../interfaces/IMessage'
import useConversationStore from '../stores/useConversationStore'

const ReceivingCallModal = () => {
  const {
    status,
    setEnableAudio,
    setEnableVideo,
    caller,
    setStatus,
    typeCall,
    setTypeCall,
    setCaller,
    initWS,
    setInitCallType
  } = useCallStore()
  const {
    setCurrentIdxSearch,
    setMessageSearches,
    setCurrentRatchetId,
    setMessages,
    conversations,
    setCurrentConversation,
    setConversations
  } = useConversationStore()
  const { userInfo } = useAuthStore()
  const { websocket } = useWebSocketStore()

  const handleJoinConversation = async () => {
    try {
      setCurrentConversation(caller)
      setConversations(
        conversations.map((item) => {
          if (item.receiver === caller) {
            return {
              ...item,
              isReaded: true
            }
          }
          return item
        })
      )

      const oldMessages = await window.api.getMessagesByUsername(caller!)

      const ratchetId = await window.api.getRatchetId(caller!)

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
      await window.api.changeConversationReaded(caller!, true)

      setCurrentIdxSearch(0)
      setMessageSearches([])
      setMessages(newMessages)
      setCurrentRatchetId(ratchetId)
    } catch (error) {
      console.error('ERROR', error)
    }
  }

  const handleAccept = () => {
    websocket!.send(
      JSON.stringify({
        type: ACCEPT_CALL_EVENT,
        plainMessage: caller,
        senderUsername: userInfo?.userName
      })
    )
    handleJoinConversation()
    setStatus('on-call')
    if (typeCall == 'video') {
      setEnableAudio(true)
      setEnableVideo(true)
    } else {
      setEnableAudio(true)
      setEnableVideo(false)
    }
    initWS('FROM_RECIEVER')
    setInitCallType('FROM_RECIEVER')
  }

  const handleReject = () => {
    setStatus('idle')
    setTypeCall(null)
    setCaller(null)
  }

  return (
    (status === 'receiving-call' || status === 'calling') && (
      <div className="absolute inset-0">
        <span
          className="absolute z-20 inset-0 bg-black/50"
          onClick={() => {
            setCaller(null)
            setTypeCall(null)
            setStatus('idle')
          }}
        ></span>
        <div className="relative w-80 z-20 bg-white py-4 px-8 rounded-lg overflow-hidden top-36 left-1/2 -translate-x-1/2">
          <div className="flex justify-center mb-3">
            <img className="w-14 rounded-full h-14" src={AVATAR_DEFAULT} />
          </div>
          {caller === userInfo!.userName ? (
            <div>Đang gọi...</div>
          ) : (
            <>
              <p className="text-center">
                <span className="font-medium">{caller} </span>
                <span>muốn gọi {typeCall === 'video' ? 'video' : ''} cho bạn</span>
              </p>

              <ul className="mt-4 flex justify-center items-center gap-5 list-none mb-0">
                <li>
                  <Phone className="cursor-pointer" onClick={handleAccept} />
                </li>

                <li className="cursor-pointer" onClick={handleReject}>
                  <X />
                </li>
              </ul>
            </>
          )}
        </div>
      </div>
    )
  )
}

export default ReceivingCallModal
