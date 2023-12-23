import axiosInstance from '../libs/axios'
import IResponse from '../interfaces/IResponse'

const uploadRepository = {
  uploadFile: (formData: FormData): Promise<IResponse<{ filePath: string }>> =>
    axiosInstance.post('/file', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    }),
  downloadFile: (filePath: string) =>
    axiosInstance.get(`/file/get?fileId=${filePath}`, { responseType: 'blob' })
}

export default uploadRepository
