import { create } from 'zustand'
import { WS_CALL_URL } from '../configs/consts'

type CallStatus = 'calling' | 'receiving-call' | 'on-call' | 'idle'

type IVideoCallStore = {
  myWS: WebSocket | null
  initWS: (type: 'FROM_CALLER' | 'FROM_RECIEVER') => void
  voipToken: string | null
  setVoipToken: (value: string | null) => void
  status: CallStatus
  setStatus: (value: CallStatus) => void
  caller: string | null
  setCaller: (value: string | null) => void
  typeCall: 'video' | 'audio' | null
  setTypeCall: (value: 'video' | 'audio' | null) => void
  enableVideo: boolean
  setEnableVideo: (value: boolean) => void
  enableAudio: boolean
  setEnableAudio: (value: boolean) => void
  encKey: string | null
  setEncKey: (value: string | null) => void
  turnOffCall: () => void
  initCallType: string | null
  setInitCallType: (type: 'FROM_CALLER' | 'FROM_RECIEVER') => void
}

const useCallStore = create<IVideoCallStore>((set, getState) => ({
  myWS: null,
  initWS: (type: 'FROM_CALLER' | 'FROM_RECIEVER') => {
    const ws = new WebSocket(
      WS_CALL_URL.replace('{{voipToken}}', getState().voipToken!).replace('{{connType}}', type)
    )
    set(() => ({ myWS: ws }))
    ws.onopen = () => {
      console.log('SenderVOIP Connected')
    }
    ws.onclose = () => {
      console.log('SenderVOIP Disconnected')
    }
  },
  voipToken: null,
  setVoipToken: (value: string | null) => set(() => ({ voipToken: value })),
  status: 'idle',
  setStatus: (value: CallStatus) => set(() => ({ status: value })),
  caller: null,
  setCaller: (value: string | null) => set(() => ({ caller: value })),
  typeCall: null,
  setTypeCall: (value: 'video' | 'audio' | null) => set(() => ({ typeCall: value })),
  enableAudio: false,
  setEnableAudio: (value: boolean) => set(() => ({ enableAudio: value })),
  enableVideo: false,
  setEnableVideo: (value: boolean) => set(() => ({ enableVideo: value })),
  encKey: null,
  setEncKey: (value: string | null) => set(() => ({ encKey: value })),
  initCallType: null,
  setInitCallType: (value: 'FROM_CALLER' | 'FROM_RECIEVER') => set(() => ({ initCallType: value })),
  turnOffCall: () => {
    if (getState().myWS) {
      console.log('getState().myWS', getState().myWS)
      getState().myWS!.close()
      set(() => ({ myWS: null }))
    }
    set(() => ({
      voipToken: null,
      status: 'idle',
      caller: null,
      typeCall: null,
      enableAudio: false,
      enableVideo: false,
      encKey: null,
      initCallType: null
    }))
  }
}))

export default useCallStore
