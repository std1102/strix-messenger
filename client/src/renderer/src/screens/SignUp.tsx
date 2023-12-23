import Input from '../components/Input'
import { FormEvent, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { HOME_PAGE, SIGN_IN_PAGE } from '../configs/routes'
import { toast } from 'react-toastify'
import authRepository from '../repositories/auth-repository'
import PinInput from 'react-pin-input'
import { ACCESS_TOKEN_KEY, AVATAR_DEFAULT, IMAGE_URL } from '../configs/consts'
import axiosInstance from '../libs/axios'
import useAuthStore from '../stores/useAuthStore'
import { hashSync } from 'bcryptjs'

const SignUpScreen = () => {
  const navigate = useNavigate()

  const { setIsAuth, setUserInfo, setAuthToken } = useAuthStore()

  const [formValue, setFormValue] = useState({
    username: '',
    password: '',
    confirmPassword: ''
  })
  const [isRegisterSuccess, setIsRegisterSuccess] = useState(false)
  const [pinValue, setPinValue] = useState('')

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    if (formValue.password !== formValue.confirmPassword) {
      toast.error('Mật khẩu không khớp')
      return
    }
    try {
      await authRepository.register({
        username: formValue.username,
        password: formValue.password
      })
      toast.success('Đăng ký thành công')
      setIsRegisterSuccess(true)
    } catch (error) {
      toast.error('Đăng ký thất bại')
      console.error('ERROR: ', error)
    }
  }

  const handleComplete = async (value: string) => {
    setPinValue(value)
  }

  const handlePinSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    try {
      await window.startUp(pinValue)
      await window.generateInternalKeyBundle()
      const keyJSON = await window.saveInternalKey()

      keyJSON.pin = hashSync(pinValue, 10)

      await window.api.writeAuthFile(keyJSON)

      const internalKey = await window.api.getInternalKey()
      if (internalKey) {
        await window.loadInternalKey(internalKey)
      }

      toast.success('Tạo mã pin thành công')

      const res = await authRepository.login(formValue)
      localStorage.setItem(ACCESS_TOKEN_KEY, res.data.accessToken)
      axiosInstance.defaults.headers['Authorization'] = `Bearer ${res.data.accessToken}`

      const externalKeyBundle = await window.populateExternalKeyBundle()

      const [userInfoRes, authToken] = await Promise.all([
        authRepository.getUserInfo(),
        authRepository.getAuthToken(),
        authRepository.uploadExternalKey(externalKeyBundle)
      ])

      userInfoRes.data.avatar = userInfoRes.data.avatar
        ? IMAGE_URL + userInfoRes.data.avatar
        : AVATAR_DEFAULT
      setAuthToken(authToken.data.authToken)
      setUserInfo(userInfoRes.data)
      setIsAuth(true)

      navigate(HOME_PAGE)
    } catch (error) {
      console.error('ERROR', error)
    }
  }

  return isRegisterSuccess ? (
    <form
      className="fixed left-1/2 text-center top-1/4 -translate-x-1/2"
      onSubmit={handlePinSubmit}
    >
      <h2 className="mb-8 text-3xl font-bold">Tạo mã pin bảo mật</h2>

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

      <button className="mt-8 bg-prim-100 text-white py-2 px-6 rounded" type="submit">
        Xác nhận
      </button>
    </form>
  ) : (
    <div className="fixed top-20 left-1/2 min-w-[500px] -translate-x-1/2">
      <div className="text-right">
        <span className="mr-4">Bạn đã có tài khoản?</span>
        <button
          className="border py-2 px-5 rounded border-gray-200"
          onClick={() => navigate(SIGN_IN_PAGE)}
        >
          Đăng nhập
        </button>
      </div>

      <h1 className="mb-3 mt-10 text-3xl font-bold">Chào mừng!</h1>
      <p className="mb-8 text-sm">
        <strong>Đăng ký </strong>
        <span className="text-gray-500">tài khoản ngay hôm nay</span>
      </p>
      <form onSubmit={handleSubmit}>
        <Input
          placeholder="Username"
          value={formValue.username}
          onChange={(e) => setFormValue({ ...formValue, username: e.target.value })}
        />
        <Input
          type="password"
          placeholder="Password"
          value={formValue.password}
          onChange={(e) => setFormValue({ ...formValue, password: e.target.value })}
        />
        <Input
          type="password"
          placeholder="Confirm password"
          value={formValue.confirmPassword}
          onChange={(e) => setFormValue({ ...formValue, confirmPassword: e.target.value })}
        />
        <button
          className="bg-prim-100 text-white font-semibold w-full rounded py-3 mt-2 hover:bg-prim-100/90 duration-200"
          type="submit"
        >
          Đăng ký
        </button>
      </form>
    </div>
  )
}

export default SignUpScreen
