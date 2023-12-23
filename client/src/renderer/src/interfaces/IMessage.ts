import { FILE_TYPE, IMAGE_TYPE, TEXT_TYPE, VIDEO_TYPE } from '../configs/consts'

interface IMessage {
  content: string
  index?: number
  type?: typeof TEXT_TYPE | typeof IMAGE_TYPE | typeof FILE_TYPE | typeof VIDEO_TYPE
  filePath?: string
  sender: string
  isDeleted: boolean
}

export default IMessage
