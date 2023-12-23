import Input from './Input'
import { ChevronDown, ChevronUp, Phone, Video } from 'lucide-react'
import { FormEvent, useEffect, useState } from 'react'
import useCallStore from '../stores/useCallStore'
import clsx from 'clsx'
import useConversationStore from '../stores/useConversationStore'
import { CALL_VIDEO_EVENT, CHAT_AUDIO_EVENT } from '../configs/consts'
import callRepository from '../repositories/call-repository'
import useAuthStore from '../stores/useAuthStore'
import userRepository from '../repositories/user-repository'

type MessageSearchProps = {
  handleScrollToMessage: (messageIndex: number) => void
}

const MessageSearch = (props: MessageSearchProps) => {
  const { handleScrollToMessage } = props

  const { userInfo } = useAuthStore()
  const {
    typeCall,
    setEnableAudio,
    setEnableVideo,
    setStatus,
    setCaller,
    setTypeCall,
    setVoipToken,
    setEncKey
  } = useCallStore()
  const {
    messages,
    currentConversation,
    messageSearches,
    currentIdxSearch,
    setMessageSearches,
    setCurrentIdxSearch
  } = useConversationStore()

  const [searchValue, setSearchValue] = useState('')

  const handleSearch = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault()

    const _messages = messages.filter((item) =>
      item.content.toLowerCase().includes(searchValue.trim().toLowerCase())
    )

    if (messages.length) {
      handleScrollToMessage(messages[0].index!)
      setMessageSearches(_messages.map((item) => item.index!.toString()))
    } else {
      setMessageSearches([])
    }
  }

  const handleAudioCall = async () => {
    try {
      const otherUserKeyBundle = await userRepository.getExternalUserKey(currentConversation!)
      const initRatchetRes = await window.initRatchetFromInternal(
        JSON.stringify(otherUserKeyBundle.data)
      )
      const res = await callRepository.initVOIP(
        currentConversation!,
        CHAT_AUDIO_EVENT,
        initRatchetRes.ephemeralKey
      )
      const ratchetInfo = await window.saveRatchet(initRatchetRes.ratchetId)
      setVoipToken(res.data.voipSession)
      setEncKey(ratchetInfo.root_key)
      setCaller(userInfo!.userName)
      setStatus('calling')
      setEnableAudio(true)
      setTypeCall('audio')
    } catch (error) {
      setTypeCall(null)
      setCaller(null)
      console.error('ERROR', error)
    }
  }

  const handleVideoCall = async () => {
    try {
      const otherUserKeyBundle = await userRepository.getExternalUserKey(currentConversation!)
      const initRatchetRes = await window.initRatchetFromInternal(
        JSON.stringify(otherUserKeyBundle.data)
      )
      const res = await callRepository.initVOIP(
        currentConversation!,
        CALL_VIDEO_EVENT,
        initRatchetRes.ephemeralKey
      )
      const ratchetInfo = await window.saveRatchet(initRatchetRes.ratchetId)
      setVoipToken(res.data.voipSession)
      setEncKey(ratchetInfo.root_key)
      setCaller(userInfo!.userName)
      setStatus('calling')
      setEnableAudio(true)
      setEnableVideo(true)
      setTypeCall('video')
    } catch (error) {
      setTypeCall(null)
      setCaller(null)
      console.error('ERROR', error)
    }
  }

  useEffect(() => {
    setSearchValue('')
  }, [currentConversation])

  return (
    <div
      className={clsx(
        'p-4 border-b flex justify-between gap-4',
        typeCall ? 'flex-col items-start' : 'items-center'
      )}
    >
      <div className="flex gap-x-3 items-center">
        <img
          src="https://source.unsplash.com/RZrIJ8C0860"
          alt="image"
          className="w-10 h-10 rounded-full"
        />
        <div className="flex flex-col">
          <h4 className="font-semibold text-sm">{currentConversation}</h4>
          <p className="text-xs text-gray-400">Active 1 hour ago</p>
        </div>
      </div>
      <div className={clsx('flex gap-x-4 items-center', typeCall && 'w-full')}>
        {!!messageSearches.length && (
          <span className="text-xs text-gray-400">
            {currentIdxSearch + 1}/{messageSearches.length}
          </span>
        )}

        <form onSubmit={handleSearch} className={clsx('relative', typeCall && 'w-full')}>
          <Input
            wrapperClass="h-full"
            inputSize={'sm'}
            placeholder="Search..."
            value={searchValue}
            onChange={(e) => setSearchValue(e.target.value)}
          />
          {!!messageSearches.length && (
            <>
              <ChevronUp
                size={14}
                className="cursor-pointer absolute top-1/2 right-2 -translate-y-1/2"
                onClick={() => {
                  const newIdx =
                    currentIdxSearch === 0 ? messageSearches.length - 1 : currentIdxSearch - 1
                  setCurrentIdxSearch(newIdx)
                  handleScrollToMessage(parseInt(messageSearches[newIdx]))
                }}
              />
              <ChevronDown
                size={14}
                className="cursor-pointer absolute top-1/2 right-6 -translate-y-1/2"
                onClick={() => {
                  const newIdx =
                    currentIdxSearch === messageSearches.length - 1 ? 0 : currentIdxSearch + 1
                  setCurrentIdxSearch(newIdx)
                  handleScrollToMessage(parseInt(messageSearches[newIdx]))
                }}
              />
            </>
          )}
          <button type="submit" hidden>
            submit
          </button>
        </form>

        {!typeCall && (
          <>
            <Phone className="cursor-pointer" onClick={handleAudioCall} />
            <Video className="cursor-pointer" onClick={handleVideoCall} />
          </>
        )}
      </div>
    </div>
  )
}

export default MessageSearch
