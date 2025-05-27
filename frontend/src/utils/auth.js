export const setAuthToken = (token) => {
  if (token) {
    localStorage.setItem('token', token);
  } else {
    localStorage.removeItem('token');
  }
};

export const getAuthToken = () => {
  return localStorage.getItem('token');
};

export const setRefreshToken = (token) => {
  if (token) {
    localStorage.setItem('refreshToken', token);
  } else {
    localStorage.removeItem('refreshToken');
  }
};

export const getRefreshToken = () => {
  return localStorage.getItem('refreshToken');
};

export const clearTokens = () => {
  localStorage.removeItem('token');
  localStorage.removeItem('refreshToken');
};

export const isAdmin = (user) => {
  return user?.isAdmin;
};