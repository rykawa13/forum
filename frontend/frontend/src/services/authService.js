export const login = async (credentials) => {
    const response = await fetch('/api/auth/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(credentials),
    });
    return handleResponse(response);
  };
  
  export const register = async (userData) => {
    const response = await fetch('/api/auth/register', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(userData),
    });
    return handleResponse(response);
  };
  
  const handleResponse = async (response) => {
    if (!response.ok) {
      const error = await response.text();
      throw new Error(error);
    }
    return response.json();
  };