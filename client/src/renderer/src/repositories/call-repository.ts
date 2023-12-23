import axiosInstance from '../libs/axios'
import { CALL_VIDEO_EVENT, CHAT_AUDIO_EVENT } from '../configs/consts'

const callRepository = {
  initVOIP: (
    username: string,
    callType: typeof CALL_VIDEO_EVENT | typeof CHAT_AUDIO_EVENT,
    ephemeralKey: string
  ) =>
    axiosInstance.put(
      `/voip/init?userId=${username}&callType=${callType}&ephemeralKey=${ephemeralKey}`
    )
}

export default callRepository
