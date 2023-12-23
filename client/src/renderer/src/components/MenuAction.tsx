import clsx from 'clsx'
import { Image, LogOut } from 'lucide-react'
import { useNavigate } from 'react-router-dom'
import { SIGN_IN_PAGE } from '../configs/routes'
import { ACCESS_TOKEN_KEY } from '../configs/consts'
import { useState } from 'react'
import ChangeAvatarModal from './ChangeAvatarModal'
import useWebSocketStore from '../stores/useWebSocketStore'
import useAuthStore from '../stores/useAuthStore'
import useConversationStore from '../stores/useConversationStore'

type MenuActionProps = {
  isMenuOpen: boolean
  setIsMenuOpen: (value: boolean) => void
}

const MenuAction = (props: MenuActionProps) => {
  const { isMenuOpen, setIsMenuOpen } = props

  const navigate = useNavigate()
  const { removeSocket } = useWebSocketStore()
  const { setCurrentConversation, setCurrentRatchetId } = useConversationStore()
  const { userInfo, setAuthToken } = useAuthStore()

  const [isOpenChangeAvatar, setIsOpenChangeAvatar] = useState(false)

  const handleLogout = () => {
    localStorage.removeItem(ACCESS_TOKEN_KEY)
    setAuthToken('')
    setCurrentConversation(null)
    setCurrentRatchetId(null)
    removeSocket()
    navigate(SIGN_IN_PAGE)
  }

  return (
    <div
      className={clsx(
        'fixed duration-300 inset-0 z-20',
        isMenuOpen ? 'visible opacity-100' : 'opacity-0 invisible'
      )}
    >
      {isOpenChangeAvatar && <ChangeAvatarModal setIsOpen={setIsOpenChangeAvatar} />}

      <div
        className="fixed duration-300 inset-0 bg-black/20"
        onClick={() => setIsMenuOpen(false)}
      ></div>

      <div
        className={clsx(
          'absolute top-0 duration-300 bottom-0 left-0 w-1/4 px-4 py-8 bg-yellow-50 shadow-xl',
          isMenuOpen ? 'translate-x-0' : '-translate-x-full'
        )}
      >
        <div className="flex gap-x-3 items-center">
          <img src={userInfo?.avatar} alt="image" className="w-10 h-10 rounded-full" />
          <div className="flex flex-col">
            <h4 className="font-semibold text-sm">{userInfo?.userName}</h4>
          </div>
        </div>

        <ul className="mt-8">
          <li
            className="flex p-2 gap-x-3 items-center cursor-pointer hover:bg-yellow-500/20 transition-colors"
            onClick={() => setIsOpenChangeAvatar(true)}
          >
            <Image />
            <span>Change Avatar</span>
          </li>
          <li
            className="flex p-2 gap-x-3 items-center cursor-pointer hover:bg-yellow-500/20 transition-colors"
            onClick={handleLogout}
          >
            <LogOut />
            <span>Logout</span>
          </li>
        </ul>
      </div>
    </div>
  )
}

export default MenuAction
