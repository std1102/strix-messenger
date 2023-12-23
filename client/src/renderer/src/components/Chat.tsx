import MessageList from './MessageList'
import MessageForm from './MessageForm'
import MessageSearch from './MessageSearch'
import { useCallback, useRef } from 'react'

const Chat = () => {
  const lastMessageRef = useRef<HTMLDivElement>(null)

  const handleScrollBottom = useCallback(() => {
    lastMessageRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [])

  const handleScrollToMessage = useCallback((messageIndex: number) => {
    const messageElement = document.getElementById(messageIndex.toString())
    messageElement?.scrollIntoView({ behavior: 'smooth' })
  }, [])

  const handleScroll = () => {
    requestIdleCallback(() => {
      lastMessageRef.current?.scrollIntoView({ behavior: 'smooth' })
    })
  }

  return (
    <>
      <MessageSearch handleScrollToMessage={handleScrollToMessage} />

      <div className="pt-4 flex-1 mb-16 overflow-auto">
        <MessageList lastMessageRef={lastMessageRef} handleScrollBottom={handleScrollBottom} />
      </div>

      <div className="mt-5 absolute bottom-2 inset-x-0">
        <MessageForm handleScroll={handleScroll} />
      </div>
    </>
  )
}

export default Chat
