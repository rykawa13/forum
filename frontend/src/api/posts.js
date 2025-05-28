import axios from 'axios';
import { enrichWithUserInfo } from './users';

const API_URL = process.env.REACT_APP_FORUM_URL || 'http://localhost:8082';

// Создаем переменную для хранения таймера дебаунсинга
let debounceTimer;

export const fetchPosts = async () => {
  try {
    // Очищаем предыдущий таймер, если он есть
    if (debounceTimer) {
      clearTimeout(debounceTimer);
    }

    // Возвращаем Promise, который резолвится после задержки
    return new Promise((resolve, reject) => {
      debounceTimer = setTimeout(async () => {
        try {
          console.log('Fetching posts from:', `${API_URL}/posts`);
          const response = await axios.get(`${API_URL}/posts`, {
            withCredentials: true
          });
          
          // Подробное логирование структуры ответа
          console.log('Raw response:', {
            status: response.status,
            headers: response.headers,
            data: response.data
          });

          // Проверяем и преобразуем данные
          const responseData = response.data;
          let posts = Array.isArray(responseData) ? responseData :
                     Array.isArray(responseData.posts) ? responseData.posts :
                     [];
          
          // Обогащаем посты информацией о пользователях
          posts = await enrichWithUserInfo(posts);
          
          console.log('Processed posts:', {
            postsCount: posts.length,
            firstPost: posts[0],
            isArray: Array.isArray(posts)
          });

          // Возвращаем данные в нужном формате
          resolve({
            data: {
              posts: posts,
              total: posts.length
            }
          });
        } catch (error) {
          console.error('Error fetching posts:', error.response || error);
          reject(error);
        }
      }, 300); // Добавляем задержку в 300мс
    });
  } catch (error) {
    console.error('Error fetching posts:', error.response || error);
    throw error;
  }
};

export const createPost = async (postData) => {
  try {
    const response = await axios.post(`${API_URL}/posts`, postData, {
      withCredentials: true,
    });
    console.log('Create post response:', response);
    return response;
  } catch (error) {
    console.error('Error creating post:', error.response || error);
    throw error;
  }
};

export const deletePost = async (postId) => {
  try {
    const response = await axios.delete(`${API_URL}/posts/${postId}`, {
      withCredentials: true,
    });
    console.log('Delete post response:', response);
    return response;
  } catch (error) {
    console.error('Error deleting post:', error.response || error);
    throw error;
  }
};