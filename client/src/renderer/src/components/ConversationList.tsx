import ConversationItem from './ConversationItem'
import useConversationStore from '../stores/useConversationStore'
import useWebSocketStore from '../stores/useWebSocketStore'
import { useCallback, useEffect } from 'react'
import {
  ACCEPT_CALL_EVENT,
  CALL_VIDEO_EVENT,
  CHAT_AUDIO_EVENT,
  CHAT_CLOSE,
  CHAT_NEW_EVENT,
  FILE_TYPE,
  IMAGE_TYPE,
  TEXT_TYPE,
  VIDEO_TYPE
} from '../configs/consts'
import IConversation from '../interfaces/IConversation'
import { toast } from 'react-toastify'
import chatRepository from '../repositories/chat-repository'
import useCallStore from '../stores/useCallStore'
import userRepository from '../repositories/user-repository'
import IMessage from '../interfaces/IMessage'

const ConversationList = () => {
  const {
    setChatAction,
    setMessages,
    currentConversation,
    currentRatchetId,
    messages,
    conversations,
    setConversations
  } = useConversationStore()
  const { websocket } = useWebSocketStore()
  const {
    setStatus,
    setTypeCall,
    setCaller,
    initWS,
    setVoipToken,
    setEncKey,
    setInitCallType,
    setEnableAudio,
    setEnableVideo,
    turnOffCall,
    typeCall
  } = useCallStore()

  // eslint-disable-next-line
  const createRatchet = useCallback(async (additionalData: any, senderUserName: string) => {
    const ratchetRes = await window.initRatchetFromExternal(
      JSON.stringify(additionalData.senderKeyBundle),
      additionalData.ephemeralKey,
      additionalData.chatSessionId
    )
    await chatRepository.completeChatSession(ratchetRes.ratchetId)
    const ratchetDetail = await window.saveRatchet(ratchetRes.ratchetId)
    await window.api.createRatchetFile(senderUserName, ratchetDetail, ratchetRes.ratchetId)
  }, [])

  useEffect(() => {
    if (!websocket) return
    websocket.onmessage = async (msg) => {
      const data = JSON.parse(msg.data)
      console.log('data', data)
      switch (data.type) {
        case CHAT_NEW_EVENT: {
          const newConversation: IConversation = {
            id: data.chatSessionId,
            receiver: data.senderUsername,
            lastMessage: '',
            updatedAt: new Date().toISOString(),
            isReaded: false
          }
          setConversations([newConversation, ...conversations])
          if (data.plainMessage) toast.info(data.plainMessage)
          createRatchet(data.additionalData, data.senderUsername)
          break
        }

        case TEXT_TYPE:
        case IMAGE_TYPE:
        case VIDEO_TYPE:
        case FILE_TYPE: {
          let mediaContent = ''
          if (data.type === IMAGE_TYPE) mediaContent = 'Người dùng đã gửi ảnh'
          else if (data.type === VIDEO_TYPE) mediaContent = 'Người dùng đã gửi video'

          const content = await window.receiveMessage(
            JSON.stringify({
              chatSessionId: data.chatSessionId,
              index: data.index,
              cipherMessage: data.cipherMessage,
              isBinary: data.isBinary
            })
          )

          const newMessage: IMessage = {
            index: data.index,
            content: mediaContent == '' ? content : mediaContent,
            type: data.type,
            sender: data.senderUsername,
            filePath: content + ':' + data.filePath,
            isDeleted: false
          }

          if (data.senderUsername === currentConversation) {
            setMessages([...messages, newMessage])
            setChatAction('add')
          }

          await window.api.addMessageToRatchet(data.senderUsername!, [newMessage])

          const conversationIndex = conversations.findIndex(
            (item) => item.id === data.chatSessionId
          )
          const newConversations = [...conversations]
          if (!newConversations[conversationIndex]) return
          newConversations[conversationIndex].lastMessage = content

          const temp = newConversations[conversationIndex]
          temp.isReaded = temp.id === currentRatchetId
          newConversations.splice(conversationIndex, 1)
          newConversations.unshift(temp)
          setConversations(newConversations)

          await window.api.changeConversationReaded(data.senderUsername, temp.isReaded)
          // save new ratchet detail
          const ratchetDetail = await window.saveRatchet(data.chatSessionId!)
          await window.api.changeRatchetDetail(data.senderUsername!, ratchetDetail)
          break
        }

        case CHAT_AUDIO_EVENT: {
          console.log('voice call comming')
          setStatus('receiving-call')
          setTypeCall('audio')
          if (typeCall == 'video') {
            setEnableAudio(true)
            setEnableVideo(true)
          } else {
            setEnableAudio(true)
            setEnableVideo(false)
          }
          setCaller(data.senderUsername)
          setVoipToken(data.cipherMessage)
          const res = await userRepository.getExternalUserKey(data.senderUsername)
          const ratchetRes = await window.initRatchetFromExternal(
            JSON.stringify(res.data),
            data.plainMessage,
            ''
          )
          const ratchetDetail = await window.saveRatchet(ratchetRes.ratchetId)
          setEncKey(ratchetDetail.root_key)
          break
        }

        case CALL_VIDEO_EVENT: {
          console.log('video call comming')
          setStatus('receiving-call')
          setTypeCall('video')
          if (typeCall == 'video') {
            setEnableAudio(true)
            setEnableVideo(true)
          } else {
            setEnableAudio(true)
            setEnableVideo(false)
          }
          setCaller(data.senderUsername)
          setVoipToken(data.cipherMessage)
          const res = await userRepository.getExternalUserKey(data.senderUsername)
          const ratchetRes = await window.initRatchetFromExternal(
            JSON.stringify(res.data),
            data.plainMessage,
            ''
          )
          const ratchetDetail = await window.saveRatchet(ratchetRes.ratchetId)
          setEncKey(ratchetDetail.root_key)
          break
        }

        // caller
        case ACCEPT_CALL_EVENT: {
          setStatus('on-call')
          initWS('FROM_CALLER')
          setInitCallType('FROM_CALLER')
          break
        }

        case CHAT_CLOSE: {
          turnOffCall()
          break
        }
      }
    }
  }, [websocket, conversations, createRatchet, messages, setMessages])

  return (
    <div className="w-full h-full overflow-y-auto custom__scroll">
      {conversations.map((conversation, index) => (
        <ConversationItem conversation={conversation} key={index} />
      ))}
    </div>
  )
}

export default ConversationList
