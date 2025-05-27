import axios from 'axios';

const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:5000/api';

export const fetchUsers = async () => {
  return await axios.get(`${API_URL}/admin/users`, {
    withCredentials: true,
  });
};

export const deleteUser = async (userId) => {
  return await axios.delete(`${API_URL}/admin/users/${userId}`, {
    withCredentials: true,
  });
};

export const getStats = async () => {
  return await axios.get(`${API_URL}/admin/stats`, {
    withCredentials: true,
  });
};