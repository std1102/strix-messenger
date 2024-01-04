import { useRef, useState } from 'react'
import { PenLine } from 'lucide-react'
import userRepository from '../repositories/user-repository'
import { toast } from 'react-toastify'
import { IMAGE_URL } from '../configs/consts'
import useAuthStore from '../stores/useAuthStore'

type ChangeAvatarModalProps = {
  setIsOpen: (value: boolean) => void
}

const ChangeAvatarModal = (props: ChangeAvatarModalProps) => {
  const { setIsOpen } = props

  const { userInfo, setUserInfo } = useAuthStore()

  const inputRef = useRef<HTMLInputElement>(null)
  const [avatar, setAvatar] = useState('')

  const handlePreviewAvatar = (e) => {
    setAvatar(URL.createObjectURL(e.target.files[0]))
  }

  const handleUploadAvatar = async () => {
    try {
      if (!inputRef.current?.files?.[0]) return

      const formData = new FormData()
      formData.append('upload', inputRef.current!.files[0])
      const res = await userRepository.uploadAvatar(formData)
      setIsOpen(false)
      // eslint-disable-next-line @typescript-eslint/ban-ts-comment
      setUserInfo({ ...userInfo, avatar: IMAGE_URL + res.data.filePath } as never)
      toast.success('Cập nhật ảnh đại diện thành công')
    } catch (error) {
      console.error('ERROR', error)
    }
  }

  return (
    <div className="fixed inset-0 z-20">
      <div className="absolute inset-0 bg-black/30" onClick={() => setIsOpen(false)}></div>
      <div
        className="w-[500px] absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 p-8 h-[400px] bg-white rounded-lg flex items-center flex-col">
        <div className="relative group">
          <input ref={inputRef} type="file" hidden onChange={handlePreviewAvatar} />
          <span
            className="font-bold absolute inset-0 z-20 bg-black/40 opacity-0 group-hover:opacity-100 transition-colors rounded-full flex items-center justify-center text-white cursor-pointer"
            onClick={() => inputRef.current?.click()}
          >
            <PenLine />
          </span>
          <img
            alt="avatar"
            src={avatar}
            className="w-56 h-56 bg-blue-500 rounded-full object-cover"
          />
        </div>

        <button
          className="bg-prim-100 text-white font-semibold w-24 rounded py-3 mt-8 hover:bg-prim-100/90 duration-200"
          onClick={handleUploadAvatar}
        >
          Lưu
        </button>
      </div>
    </div>
  )
}

export default ChangeAvatarModal
