import clsx from 'clsx'
import { File, Trash } from 'lucide-react'
import IMessage from '../interfaces/IMessage'
import { useEffect, useState } from 'react'
import useAuthStore from '../stores/useAuthStore'
import { FILE_TYPE, IMAGE_TYPE, TEXT_TYPE, VIDEO_TYPE } from '../configs/consts'
import uploadRepository from '../repositories/upload-repository'
import { decryptblobBrowser } from '../crypto/cryptoLib'
import useConversationStore from '../stores/useConversationStore'
import { toast } from 'react-toastify'

type MessageItemProps = {
  message: IMessage
}

const MessageItem = (props: MessageItemProps) => {
  const { message } = props

  const { userInfo } = useAuthStore()
  const {
    currentConversation,
    messageSearches,
    messages,
    setMessages,
    conversations,
    setConversations
  } = useConversationStore()

  const [content, setContent] = useState(message.type === TEXT_TYPE ? message.content : '')
  const [filename, setFileName] = useState('')
  const [fileSize, setFileSize] = useState(0)
  const [isOpen, setIsOpen] = useState(false)

  const isSender = message.sender === userInfo?.userName

  const handleOpenModal = () => {
    setIsOpen(true)
  }

  const handleDeleteMessage = async () => {
    try {
      await window.api.deleteMessage(currentConversation!, message.index! - 1)
      const newMessages = messages.map((item) => {
        if (item.index === message.index) {
          return {
            ...item,
            isDeleted: true
          }
        }
        return item
      })
      if (message.index! === messages.length) {
        const newConversations = conversations.map((item) => {
          if (item.receiver === currentConversation) {
            return {
              ...item,
              lastMessage: 'Tin nhắn đã bị xóa'
            }
          }
          return item
        })
        setConversations(newConversations)
      }

      setMessages(newMessages)
    } catch (error) {
      console.error('ERROR', error)
      toast.error('Xóa tin nhắn thất bại')
    } finally {
      setIsOpen(false)
    }
  }

  useEffect(() => {
    if (message.type !== TEXT_TYPE && message.filePath) {
      ;(async function () {
        console.log(message)
        const [randomKey, filePath, , filename, fileSize] = message.filePath!.split(':')
        const res = await uploadRepository.downloadFile(filePath!)
        const responseBlob = await res.data
        const decryptedData = await decryptblobBrowser(responseBlob, randomKey)
        const url = URL.createObjectURL(decryptedData)
        setContent(url)

        if (message.type === FILE_TYPE) {
          setFileName(filename)
          setFileSize(parseInt(fileSize))
        }
      })()
    }
  }, [message])

  return (
    <div
      key={content}
      id={message.index!.toString()}
      className={clsx('flex gap-x-3 group', isSender ? 'flex-row-reverse' : 'items-start')}
    >
      {isOpen && (
        <div className="absolute inset-0 z-50">
          <div className="relative border w-96 z-20 bg-white py-4 px-8 rounded-lg overflow-hidden top-36 left-1/2 -translate-x-1/2 shadow-2xl">
            <h4 className="font-medium">Bạn có chắc chắn muốn xóa tin nhắn này?</h4>
            <div className="flex justify-center gap-x-4 mt-4">
              <button
                className="bg-green-500 text-white px-4 py-2 rounded-lg"
                onClick={handleDeleteMessage}
              >
                Xóa
              </button>
              <button
                className="bg-red-500 text-white px-4 py-2 rounded-lg"
                onClick={() => setIsOpen(false)}
              >
                Hủy
              </button>
            </div>
          </div>
        </div>
      )}

      <img src={userInfo!.avatar} alt="image" className="w-10 h-10 rounded-full" />
      <div className={clsx('flex items-center gap-x-4', isSender && 'flex-row-reverse')}>
        <div
          className={clsx(
            'flex flex-col rounded p-2',
            isSender ? 'bg-[#2A5251] text-white' : 'bg-[#F3F2EE]'
          )}
        >
          {message.isDeleted ? (
            <span className="text-xs">Tin nhắn đã bị xóa</span>
          ) : (
            <>
              {message.type === TEXT_TYPE && (
                <p
                  className={clsx(
                    'text-xs ',
                    messageSearches.includes(message.index!.toString()) &&
                      'bg-yellow-200/80 text-black'
                  )}
                >
                  {content}
                </p>
              )}
              {message.type === IMAGE_TYPE && content && (
                <img
                  alt={'loading'}
                  key={content}
                  src={content}
                  className="object-cover max-w-[100px]"
                />
              )}
              {message.type === VIDEO_TYPE && content && (
                <video
                  src={content}
                  key={content}
                  className="object-cover max-w-[300px]"
                  controls
                />
              )}
              {message.type === FILE_TYPE && content && (
                <a href={content} className="flex gap-x-4 max-w-[400px]" download={filename}>
                  <File size={40} />
                  <div className="flex flex-col">
                    <span>{filename}</span>
                    <span>{fileSize}Kb</span>
                  </div>
                </a>
              )}
            </>
          )}
        </div>

        <div className="opacity-0 group-hover:opacity-100 duration-200">
          {isSender && !message.isDeleted && (
            <Trash className="cursor-pointer" size={14} onClick={handleOpenModal} />
          )}
        </div>
      </div>
    </div>
  )
}

export default MessageItem
