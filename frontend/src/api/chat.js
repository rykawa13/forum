import axios from 'axios';

const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:5000/api';

export const fetchMessagesApi = async () => {
  return await axios.get(`${API_URL}/chat/messages`);
};

export const sendMessageApi = async (message) => {
  return await axios.post(`${API_URL}/chat/send`, { message }, {
    withCredentials: true,
  });
};