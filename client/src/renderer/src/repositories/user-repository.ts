import axiosInstance from '../libs/axios'
import IResponse from '../interfaces/IResponse'

const userRepository = {
  searchUser: (search: string) => axiosInstance.get(`/user/search?keyWord=${search}`),
  getExternalUserKey: (userId: string) => axiosInstance.get(`/user/${userId}/externalKey`),
  uploadAvatar: (formData: FormData): Promise<IResponse<{ filePath: string }>> =>
    axiosInstance.post('/file/avatar', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    }),
  downloadAvatar: (fileId: string) => axiosInstance.get(`/file?fileId=${fileId}`)
}

export default userRepository
