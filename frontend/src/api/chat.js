import axios from 'axios';

const API_URL = process.env.REACT_APP_CHAT_API_URL;

export const fetchMessagesApi = async () => {
  try {
    console.log('Fetching messages from:', `${API_URL}/chat/messages`);
    const response = await axios.get(`${API_URL}/chat/messages`, {
      withCredentials: true
    });
    console.log('Messages response:', response);
    return response;
  } catch (error) {
    console.error('Error fetching messages:', error);
    throw error;
  }
};

export const sendMessageApi = async (message) => {
  try {
    console.log('Sending message to:', `${API_URL}/chat/send`);
    const response = await axios.post(`${API_URL}/chat/send`, { message }, {
      withCredentials: true,
    });
    console.log('Send message response:', response);
    return response;
  } catch (error) {
    console.error('Error sending message:', error);
    throw error;
  }
};