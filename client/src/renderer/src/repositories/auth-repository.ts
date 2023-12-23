import axiosInstance from '../libs/axios'
import IAuthData from '../interfaces/IAuthData'
import ILoginResponse from '../interfaces/ILoginResponse'
import IResponse from '../interfaces/IResponse'

const authRepository = {
  getUserInfo: () => axiosInstance.get('/user'),
  uploadExternalKey: (data) => axiosInstance.post('/user/uploadKey', data),
  login: (data: IAuthData): Promise<IResponse<ILoginResponse>> =>
    axiosInstance.post('/auth/login', { ...data, loginType: 'password' }),
  register: (data: IAuthData) => axiosInstance.post('/auth/register', data),
  getAuthToken: (): Promise<IResponse<{ authToken: string }>> => axiosInstance.get('/ws/init')
}

export default authRepository
