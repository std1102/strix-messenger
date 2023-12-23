import Input from './Input'
import { Image, Paperclip, SendHorizonal } from 'lucide-react'
import { FormEvent, useRef, useState } from 'react'
import useConversationStore from '../stores/useConversationStore'
import useWebSocketStore from '../stores/useWebSocketStore'
import useAuthStore from '../stores/useAuthStore'
import { FILE_TYPE, IMAGE_TYPE, TEXT_TYPE, VIDEO_TYPE } from '../configs/consts'
import uploadRepository from '../repositories/upload-repository'
import { encryptblobBrowser } from '../crypto/cryptoLib'
import IMessage from '../interfaces/IMessage'

type MessageFormProps = {
  handleScroll: (content: string) => void
}

const MessageForm = (props: MessageFormProps) => {
  const { handleScroll } = props

  const {
    messages,
    setConversations,
    setMessages,
    currentRatchetId,
    currentConversation,
    conversations
  } = useConversationStore()
  const { websocket } = useWebSocketStore()
  const { userInfo } = useAuthStore()

  const imageRef = useRef<HTMLInputElement>(null)
  const fileRef = useRef<HTMLInputElement>(null)
  const [content, setContent] = useState('')

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    const value = content.trim()
    if (!value) return

    const res = await window.sendMessage(currentRatchetId!, false, content)
    websocket?.send(
      JSON.stringify({
        type: TEXT_TYPE,
        senderUsername: userInfo?.userName,
        plainMessage: '',
        chatSessionId: currentRatchetId,
        index: res.index,
        cipherMessage: res.cipherMessage,
        isBinary: res.isBinary
      })
    )
    const newMessage: IMessage = {
      index: res.index,
      content: value,
      type: TEXT_TYPE,
      sender: userInfo!.userName,
      isDeleted: false
    }

    await window.api.addMessageToRatchet(currentConversation!, [newMessage])

    setMessages([...messages, newMessage])

    const conversationIndex = conversations.findIndex(
      (item) => item.receiver === currentConversation
    )
    const newConversations = [...conversations]
    newConversations[conversationIndex].lastMessage = value
    const temp = newConversations[conversationIndex]
    temp.updatedAt = new Date().toISOString()
    newConversations.splice(conversationIndex, 1)
    newConversations.unshift(temp)
    setConversations(newConversations)

    setContent('')
    handleScroll(value)
  }

  const handleSendImage = async (e: FormEvent<HTMLInputElement>) => {
    const file = e.currentTarget.files?.[0]
    if (!file) return
    const mimeType = file.type
    const randomKey = await window.api.random32Bytes()
    const reader = new FileReader()
    reader.onload = async (e) => {
      const encryptedBlob = await encryptblobBrowser(
        new Blob([e!.target!.result as ArrayBuffer]),
        randomKey
      )

      const formData = new FormData()
      formData.set('upload', encryptedBlob)
      const res = await uploadRepository.uploadFile(formData)
      const chatType = mimeType.split('/')[0] === 'image' ? IMAGE_TYPE : VIDEO_TYPE

      const resMsg = await window.sendMessage(currentRatchetId!, false, randomKey)
      websocket?.send(
        JSON.stringify({
          type: chatType,
          senderUsername: userInfo?.userName,
          filePath: res.data.filePath + ':' + mimeType,
          plainMessage: res.data.filePath,
          chatSessionId: currentRatchetId,
          index: resMsg.index,
          cipherMessage: resMsg.cipherMessage,
          isBinary: false
        })
      )
      const newMessage: IMessage = {
        index: resMsg.index,
        content: chatType === IMAGE_TYPE ? 'Bạn đã gửi một ảnh' : 'Bạn đã gửi một video',
        filePath: randomKey + ':' + res.data.filePath + ':' + mimeType,
        type: chatType,
        sender: userInfo!.userName,
        isDeleted: false
      }

      await window.api.addMessageToRatchet(currentConversation!, [newMessage])

      setMessages([...messages, newMessage])

      const conversationIndex = conversations.findIndex(
        (item) => item.receiver === currentConversation
      )
      const newConversations = [...conversations]
      newConversations[conversationIndex].lastMessage =
        chatType === IMAGE_TYPE ? 'Người dùng đã gửi một ảnh' : 'Người dùng đã gửi một video'
      const temp = newConversations[conversationIndex]
      newConversations.splice(conversationIndex, 1)
      newConversations.unshift(temp)
      setConversations(newConversations)

      setContent('')
      handleScroll(res.data.filePath)
    }
    reader.readAsArrayBuffer(file)
  }

  const handleSendFile = async (e: FormEvent<HTMLInputElement>) => {
    const file = e.currentTarget.files?.[0]
    if (!file) return
    const mimeType = file.type
    const fileName = file.name
    const fileSize = file.size
    const randomKey = await window.api.random32Bytes()
    const reader = new FileReader()
    reader.onload = async (e) => {
      const encryptedBlob = await encryptblobBrowser(
        new Blob([e!.target!.result as ArrayBuffer]),
        randomKey
      )

      const formData = new FormData()
      formData.set('upload', encryptedBlob)
      const res = await uploadRepository.uploadFile(formData)

      const resMsg = await window.sendMessage(currentRatchetId!, false, randomKey)
      websocket?.send(
        JSON.stringify({
          type: FILE_TYPE,
          senderUsername: userInfo?.userName,
          filePath: res.data.filePath + ':' + mimeType + ':' + fileName + ':' + fileSize / 1024,
          plainMessage: res.data.filePath,
          chatSessionId: currentRatchetId,
          index: resMsg.index,
          cipherMessage: resMsg.cipherMessage,
          isBinary: false
        })
      )
      const newMessage: IMessage = {
        index: resMsg.index,
        content: 'Người dùng đã gửi một tệp tin',
        filePath:
          randomKey +
          ':' +
          res.data.filePath +
          ':' +
          mimeType +
          ':' +
          fileName +
          ':' +
          fileSize / 1024,
        type: FILE_TYPE,
        sender: userInfo!.userName,
        isDeleted: false
      }

      await window.api.addMessageToRatchet(currentConversation!, [newMessage])

      setMessages([...messages, newMessage])

      const conversationIndex = conversations.findIndex(
        (item) => item.receiver === currentConversation
      )
      const newConversations = [...conversations]
      newConversations[conversationIndex].lastMessage = 'Bạn dùng đã gửi một tệp tin'
      const temp = newConversations[conversationIndex]
      newConversations.splice(conversationIndex, 1)
      newConversations.unshift(temp)
      setConversations(newConversations)

      setContent('')
      handleScroll(res.data.filePath)
    }
    reader.readAsArrayBuffer(file)
  }

  return (
    <form className="flex relative h-full px-4" onSubmit={handleSubmit}>
      <Input
        wrapperClass="h-full w-full"
        inputClass="pr-20"
        placeholder="Enter message..."
        value={content}
        onChange={(e) => setContent(e.target.value)}
      />
      <input
        ref={imageRef}
        type="file"
        hidden
        onChange={handleSendImage}
        accept="image/*,video/*"
      />
      <input ref={fileRef} type="file" hidden onChange={handleSendFile} accept="application/*" />
      <Paperclip
        className="absolute -translate-y-1/2 top-1/2 right-[5.5rem] cursor-pointer"
        onClick={() => {
          fileRef.current?.click()
        }}
      />
      <Image
        className="absolute -translate-y-1/2 top-1/2 right-14 cursor-pointer"
        onClick={() => {
          imageRef.current?.click()
        }}
      />
      <SendHorizonal className="absolute -translate-y-1/2 top-1/2 right-6 cursor-pointer" />

      <button type="submit" hidden>
        submit
      </button>
    </form>
  )
}

export default MessageForm
