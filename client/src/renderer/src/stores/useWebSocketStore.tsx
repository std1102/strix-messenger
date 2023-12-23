import { create } from 'zustand'
import useAuthStore from './useAuthStore'
import { SOCKET_URL } from '../configs/consts'

interface ISocketStoreType {
  websocket: WebSocket | null
  initWebSocket: () => void
  removeSocket: () => void
}

const useWebSocketStore = create<ISocketStoreType>((setState, getState) => ({
  websocket: null,
  initWebSocket: () => {
    const WEB_SOCKET_URL = `${SOCKET_URL}/ws?authToken=${useAuthStore.getState().authToken}`

    if (WEB_SOCKET_URL && !getState().websocket) {
      const _websocket = new WebSocket(WEB_SOCKET_URL)

      _websocket.onopen = () => {
        console.log('main ws connected')
      }

      _websocket.onclose = () => {
        console.log('main ws disconnected')
      }

      setState({ websocket: _websocket })
    }
  },
  removeSocket: () => {
    if (getState().websocket) {
      getState().websocket!.close()
      setState({ websocket: null })
    }
  }
}))

export default useWebSocketStore
