import axios from 'axios';

const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:5000/api';

export const fetchPosts = async () => {
  return await axios.get(`${API_URL}/posts`);
};

export const createPost = async (postData) => {
  return await axios.post(`${API_URL}/posts`, postData, {
    withCredentials: true,
  });
};

export const deletePost = async (postId) => {
  return await axios.delete(`${API_URL}/posts/${postId}`, {
    withCredentials: true,
  });
};