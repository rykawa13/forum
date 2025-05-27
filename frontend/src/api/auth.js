import axiosInstance from './axios';
import { getRefreshToken } from '../utils/auth';

const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8081';

export const loginUser = async (email, password) => {
  console.log('Attempting to login with:', { email, password });
  try {
    const response = await axiosInstance.post('/auth/sign-in', { email, password });
    console.log('Login response:', response.data);
    return response;
  } catch (error) {
    console.error('Login request error:', error.response?.data || error);
    throw error;
  }
};

export const registerUser = async (username, email, password) => {
  return await axiosInstance.post('/auth/sign-up', { 
    username, 
    email, 
    password 
  });
};

export const checkAuth = async () => {
  return await axiosInstance.get('/api/me');
};

export const refreshTokens = async () => {
  const refreshToken = getRefreshToken();
  return await axiosInstance.post('/auth/refresh', { refresh_token: refreshToken });
};

export const logoutUser = async () => {
  const refreshToken = getRefreshToken();
  return await axiosInstance.post('/auth/logout', { refresh_token: refreshToken });
};