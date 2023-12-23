import { ReactNode, useEffect } from 'react'
import useWebSocketStore from '../stores/useWebSocketStore'
import useAuthStore from '../stores/useAuthStore'

const SocketProvider = ({ children }: { children: ReactNode }) => {
  const { authToken } = useAuthStore()
  const { websocket, initWebSocket } = useWebSocketStore()

  useEffect(() => {
    if ((websocket && websocket.readyState !== websocket.CLOSED) || !authToken) return
    else initWebSocket()
  }, [websocket, initWebSocket, authToken])

  return children
}

export default SocketProvider
