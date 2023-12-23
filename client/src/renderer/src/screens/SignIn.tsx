import Input from '../components/Input'
import { FormEvent, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { HOME_PAGE, SIGN_UP_PAGE } from '../configs/routes'
import { toast } from 'react-toastify'
import authRepository from '../repositories/auth-repository'
import { ACCESS_TOKEN_KEY, AVATAR_DEFAULT, IMAGE_URL } from '../configs/consts'
import useAuthStore from '../stores/useAuthStore'
import axiosInstance from '../libs/axios'

const SignInScreen = () => {
  const navigate = useNavigate()

  const { setUserInfo, setAuthToken } = useAuthStore()

  const [formValue, setFormValue] = useState<{ username: string; password: string }>({
    username: '',
    password: ''
  })

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    try {
      const res = await authRepository.login(formValue)
      localStorage.setItem(ACCESS_TOKEN_KEY, res.data.accessToken)
      axiosInstance.defaults.headers['Authorization'] = `Bearer ${res.data.accessToken}`

      const [userInfoRes, authToken] = await Promise.all([
        authRepository.getUserInfo(),
        authRepository.getAuthToken()
      ])

      userInfoRes.data.avatar = userInfoRes.data.avatar
        ? IMAGE_URL + userInfoRes.data.avatar
        : AVATAR_DEFAULT
      setAuthToken(authToken.data.authToken)
      setUserInfo(userInfoRes.data)
      toast.success('Đăng nhập thành công')
      navigate(HOME_PAGE)
      localStorage.setItem(ACCESS_TOKEN_KEY, res.data.accessToken)
    } catch (error) {
      toast.error('Đăng nhập thất bại')
      console.error('ERROR', error)
    }
  }

  return (
    <div className="fixed top-20 left-1/2 min-w-[500px] -translate-x-1/2">
      <div className="text-right">
        <span className="mr-4">Bạn chưa có tài khoản?</span>
        <button
          className="border py-2 px-5 rounded border-gray-200"
          onClick={() => navigate(SIGN_UP_PAGE)}
        >
          Đăng ký
        </button>
      </div>

      <h1 className="mb-3 mt-10 text-3xl font-bold">Chào mừng!</h1>
      <p className="mb-8 text-sm">
        <strong>Đăng nhập </strong>
        <span className="text-gray-500">để trải nghiệm những nội dung thú vị</span>
      </p>
      <form onSubmit={handleSubmit}>
        <Input
          value={formValue.username}
          placeholder="Username"
          onChange={(e) => setFormValue({ ...formValue, username: e.target.value })}
        />
        <Input
          type="password"
          value={formValue.password}
          placeholder="Password"
          onChange={(e) => setFormValue({ ...formValue, password: e.target.value })}
        />
        <button
          className="bg-prim-100 text-white font-semibold w-full rounded py-3 mt-2 hover:bg-prim-100/90 duration-200"
          type="submit"
        >
          Đăng nhập
        </button>
      </form>
    </div>
  )
}

export default SignInScreen
