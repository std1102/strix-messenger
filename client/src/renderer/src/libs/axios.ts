import axios from 'axios'
import { API_URL } from '../configs/consts'

const axiosInstance = axios.create({
  baseURL: `${API_URL}/api/v1`,
  headers: {
    'Content-Type': 'application/json'
    // Authorization: localStorage.getItem(ACCESS_TOKEN_KEY)
    //   ? `Bearer ${localStorage.getItem(ACCESS_TOKEN_KEY)}`
    //   : undefined
  }
})

export default axiosInstance
