import axios from 'axios';

const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8081';

const axiosInstance = axios.create({
  baseURL: API_URL,
  withCredentials: true,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Добавляем перехватчик запросов для добавления токена
axiosInstance.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Добавляем перехватчик ответов для обработки ошибок
axiosInstance.interceptors.response.use(
  (response) => response,
  async (error) => {
    if (error.response) {
      // Если ошибка 401 (Unauthorized), пытаемся обновить токен
      if (error.response.status === 401) {
        const refreshToken = localStorage.getItem('refreshToken');
        if (refreshToken) {
          try {
            const response = await axios.post(`${API_URL}/auth/refresh`, {
              refresh_token: refreshToken
            });
            
            const { access_token, refresh_token } = response.data;
            localStorage.setItem('token', access_token);
            localStorage.setItem('refreshToken', refresh_token);

            // Повторяем оригинальный запрос с новым токеном
            error.config.headers.Authorization = `Bearer ${access_token}`;
            return axiosInstance(error.config);
          } catch (refreshError) {
            // Если не удалось обновить токен, очищаем хранилище
            localStorage.removeItem('token');
            localStorage.removeItem('refreshToken');
            window.location.href = '/login';
          }
        }
      }
    }
    return Promise.reject(error);
  }
);

export default axiosInstance; 