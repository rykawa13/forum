import axios from 'axios';
import { getAuthToken, setAuthToken, getRefreshToken, setRefreshToken, clearTokens } from '../utils/auth';

const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8081';

// Создаем обработчик для 401 ошибки
let unauthorizedCallback = () => {};

export const setUnauthorizedCallback = (callback) => {
  unauthorizedCallback = callback;
};

const axiosInstance = axios.create({
  baseURL: API_URL,
});

// Флаг, указывающий, что идет обновление токена
let isRefreshing = false;
// Очередь запросов, ожидающих обновления токена
let failedQueue = [];

const processQueue = (error, token = null) => {
  failedQueue.forEach(prom => {
    if (error) {
      prom.reject(error);
    } else {
      prom.resolve(token);
    }
  });
  failedQueue = [];
};

// Добавляем токен к каждому запросу
axiosInstance.interceptors.request.use(
  (config) => {
    const token = getAuthToken();
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Обрабатываем ответы и ошибки
axiosInstance.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;

    // Если ошибка 401 и это не запрос на обновление токена
    if (error.response?.status === 401 && !originalRequest._retry) {
      if (isRefreshing) {
        // Если уже идет обновление токена, добавляем запрос в очередь
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject });
        })
          .then(token => {
            originalRequest.headers['Authorization'] = 'Bearer ' + token;
            return axiosInstance(originalRequest);
          })
          .catch(err => Promise.reject(err));
      }

      originalRequest._retry = true;
      isRefreshing = true;

      const refreshToken = getRefreshToken();
      if (!refreshToken) {
        clearTokens();
        unauthorizedCallback();
        return Promise.reject(error);
      }

      try {
        const response = await axios.post(`${API_URL}/auth/refresh`, {
          refresh_token: refreshToken
        });

        const { access_token, refresh_token } = response.data;
        setAuthToken(access_token);
        setRefreshToken(refresh_token);

        // Обновляем заголовок для текущего запроса
        originalRequest.headers['Authorization'] = 'Bearer ' + access_token;
        
        // Обрабатываем очередь запросов
        processQueue(null, access_token);
        
        return axiosInstance(originalRequest);
      } catch (refreshError) {
        processQueue(refreshError, null);
        clearTokens();
        unauthorizedCallback();
        return Promise.reject(refreshError);
      } finally {
        isRefreshing = false;
      }
    }
    return Promise.reject(error);
  }
);

export default axiosInstance; 