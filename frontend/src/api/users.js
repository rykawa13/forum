import axiosInstance from './axios';

const API_URL = process.env.REACT_APP_FORUM_URL || 'http://localhost:8082';
const AUTH_URL = process.env.REACT_APP_API_URL || 'http://localhost:8081';

// Кэш для хранения информации о пользователях
const userCache = new Map();

export const fetchUserInfo = async (userId) => {
  try {
    // Проверяем кэш
    if (userCache.has(userId)) {
      return userCache.get(userId);
    }

    // Пытаемся получить информацию о пользователе
    try {
      const response = await axiosInstance.get(`/api/users/${userId}`);
      const userData = response.data;
      
      // Сохраняем в кэш
      userCache.set(userId, userData);
      return userData;
    } catch (error) {
      console.error(`Error fetching user ${userId} info:`, error);
      // В случае ошибки возвращаем базовую информацию
      const basicUserInfo = {
        id: userId,
        username: `User ${userId}`,
      };
      userCache.set(userId, basicUserInfo);
      return basicUserInfo;
    }
  } catch (error) {
    console.error(`Error in fetchUserInfo for user ${userId}:`, error);
    return {
      id: userId,
      username: `User ${userId}`,
    };
  }
};

export const enrichWithUserInfo = async (items) => {
  if (!Array.isArray(items)) return items;

  const enrichedItems = await Promise.all(
    items.map(async (item) => {
      if (!item.author_id) return item;

      const userInfo = await fetchUserInfo(item.author_id);
      return {
        ...item,
        author: userInfo || { username: `User ${item.author_id}` }
      };
    })
  );

  return enrichedItems;
}; 