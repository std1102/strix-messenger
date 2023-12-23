import axiosInstance from '../libs/axios'
import IChatSessionInit from '../interfaces/IChatSessionInit'

const chatRepository = {
  initChatSession: (body: IChatSessionInit) => axiosInstance.post('/chatSession/init', body),
  getPendingChatSession: () => axiosInstance.get('/chatSession'),
  getPendingChatMessage: (chatSessionId: string) =>
    axiosInstance.get(`/message?chatSessionId=${chatSessionId}`),
  completeChatSession: (chatSessionId: string) =>
    axiosInstance.put(`chatSession/complete?chatSessionId=${chatSessionId}`)
}

export default chatRepository
