import PinInput from 'react-pin-input'
import { FormEvent, ReactNode, useCallback, useEffect, useState } from 'react'
import useAuthStore from '../stores/useAuthStore'
import SocketProvider from '../providers/SocketProvider'
import { useNavigate } from 'react-router-dom'
import { SIGN_IN_PAGE } from '../configs/routes'
import { ACCESS_TOKEN_KEY, AVATAR_DEFAULT, IMAGE_URL } from '../configs/consts'
import authRepository from '../repositories/auth-repository'
import { Loader2 } from 'lucide-react'
import IAuthFile from '../interfaces/IAuthFile'
import axiosInstance from '../libs/axios'
import { compareSync } from 'bcryptjs'

const PinAuthentication = ({ children }: { children: ReactNode }) => {
  const navigate = useNavigate()

  const { isAuth, setAuthToken, setUserInfo, setIsAuth } = useAuthStore()

  const [pinValue, setPinValue] = useState('')
  const [errorMessage, setErrorMessage] = useState('')
  const [isInit, setIsInit] = useState(true)

  const handleComplete = (value: string) => {
    setPinValue(value)
  }

  const handleCheckPin = async () => {
    const keyStrSaved = await window.api.readAuthFile()
    if (!keyStrSaved) navigate(SIGN_IN_PAGE)

    const keySaved = JSON.parse(keyStrSaved) as IAuthFile

    const isMatch = compareSync(pinValue, keySaved.pin)

    if (isMatch) {
      try {
        await window.startUp(pinValue)

        const internalKey = await window.api.getInternalKey()
        if (internalKey) {
          await window.loadInternalKey(internalKey)
        }

        const [userInfoRes, authToken] = await Promise.all([
          authRepository.getUserInfo(),
          authRepository.getAuthToken()
        ])

        setAuthToken(authToken.data.authToken)
        userInfoRes.data.avatar = userInfoRes.data.avatar
          ? IMAGE_URL + userInfoRes.data.avatar
          : AVATAR_DEFAULT
        console.log('userInfoRes', userInfoRes)
        setUserInfo(userInfoRes.data)
        setIsAuth(true)
        setErrorMessage('')
      } catch (error) {
        navigate(SIGN_IN_PAGE)
      }
    } else {
      setErrorMessage('Mã pin không đúng')
    }
  }

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    handleCheckPin().then()
  }

  const autoSignIn = useCallback(async () => {
    const isExistAuthFile = await window.api.existAuthFile()
    if (!isExistAuthFile || !localStorage.getItem(ACCESS_TOKEN_KEY)) {
      localStorage.removeItem(ACCESS_TOKEN_KEY)
      navigate(SIGN_IN_PAGE)
    } else {
      axiosInstance.defaults.headers['Authorization'] = `Bearer ${localStorage.getItem(
        ACCESS_TOKEN_KEY
      )}`
    }

    setIsInit(false)
  }, [])

  useEffect(() => {
    autoSignIn()
  }, [autoSignIn])

  return isInit ? (
    <div className="w-screen h-screen flex items-center justify-center">
      <Loader2 className=" animate-spin w-24 h-24" />
    </div>
  ) : (
    <SocketProvider>
      {isAuth ? (
        children
      ) : (
        <form
          className="fixed left-1/2 text-center top-1/4 -translate-x-1/2"
          onSubmit={handleSubmit}
        >
          <h2 className="mb-8 text-3xl font-bold">Vui lòng nhập mã pin</h2>

          <PinInput
            length={4}
            initialValue=""
            focus
            inputStyle={{
              borderRadius: 4,
              margin: '0 10px'
            }}
            onComplete={handleComplete}
          />
          {errorMessage && <p className="text-red-500 mt-6">{errorMessage}</p>}

          <button className="mt-8 bg-prim-100 text-white py-2 px-6 rounded" type="submit">
            Xác nhận
          </button>
        </form>
      )}
    </SocketProvider>
  )
}

export default PinAuthentication
