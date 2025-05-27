import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import { loginUser, registerUser, checkAuth as checkAuthApi, logoutUser, refreshTokens } from '../api/auth';
import { setAuthToken, getAuthToken, setRefreshToken, getRefreshToken, clearTokens } from '../utils/auth';

// Создаем начальное состояние
const initialState = {
  user: null,
  isAuthenticated: false,
  loading: false,
  error: null,
  // Добавляем флаг инициализации
  isInitialized: false,
};

export const login = createAsyncThunk(
  'auth/login',
  async ({ email, password }, { rejectWithValue }) => {
    try {
      const response = await loginUser(email, password);
      const { access_token, refresh_token } = response.data;
      setAuthToken(access_token);
      setRefreshToken(refresh_token);
      return response.data;
    } catch (error) {
      console.error('Login error:', error.response?.data || error.message);
      return rejectWithValue(
        error.response?.data?.error || 
        error.response?.data || 
        error.message || 
        'Ошибка входа'
      );
    }
  }
);

export const register = createAsyncThunk(
  'auth/register',
  async ({ username, email, password }, { rejectWithValue }) => {
    try {
      const response = await registerUser(username, email, password);
      const { access_token, refresh_token } = response.data;
      setAuthToken(access_token);
      setRefreshToken(refresh_token);
      return response.data;
    } catch (error) {
      console.error('Registration error:', error.response?.data || error.message);
      return rejectWithValue(
        error.response?.data?.error || 
        error.response?.data || 
        error.message || 
        'Ошибка регистрации'
      );
    }
  }
);

export const refresh = createAsyncThunk(
  'auth/refresh',
  async (_, { rejectWithValue }) => {
    try {
      const refreshToken = getRefreshToken();
      if (!refreshToken) {
        throw new Error('No refresh token found');
      }

      const response = await refreshTokens();
      const { access_token, refresh_token } = response.data;
      setAuthToken(access_token);
      setRefreshToken(refresh_token);
      return response.data;
    } catch (error) {
      clearTokens();
      return rejectWithValue(error.response?.data || error.message);
    }
  }
);

export const checkAuth = createAsyncThunk(
  'auth/checkAuth',
  async (_, { dispatch, rejectWithValue }) => {
    try {
      // Сначала проверяем наличие refresh token
      const refreshToken = getRefreshToken();
      if (!refreshToken) {
        throw new Error('No refresh token found');
      }

      // Проверяем access token
      const accessToken = getAuthToken();
      if (!accessToken) {
        // Если нет access token, но есть refresh token - пробуем обновить токены
        await dispatch(refresh()).unwrap();
      }
      
      const response = await checkAuthApi();
      return response.data;
    } catch (error) {
      if (error.response?.status === 401) {
        // При 401 ошибке пробуем обновить токены
        try {
          await dispatch(refresh()).unwrap();
          // После успешного обновления токенов повторяем запрос
          const response = await checkAuthApi();
          return response.data;
        } catch (refreshError) {
          clearTokens();
          return rejectWithValue(refreshError.response?.data || refreshError.message);
        }
      }
      return rejectWithValue(error.response?.data || error.message);
    }
  }
);

export const logout = createAsyncThunk(
  'auth/logout',
  async (_, { rejectWithValue }) => {
    try {
      await logoutUser();
      clearTokens();
    } catch (error) {
      console.error('Logout error:', error);
      // Даже при ошибке очищаем токены
      clearTokens();
    }
  }
);

const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    clearError: (state) => {
      state.error = null;
    },
    // Добавляем действие для сброса состояния
    resetState: () => initialState,
  },
  extraReducers: (builder) => {
    builder
      // Login cases
      .addCase(login.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(login.fulfilled, (state, action) => {
        state.loading = false;
        state.isAuthenticated = true;
        state.user = action.payload.user;
        state.error = null;
        state.isInitialized = true;
      })
      .addCase(login.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload;
        state.isAuthenticated = false;
        state.user = null;
        state.isInitialized = true;
      })
      // Register cases
      .addCase(register.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(register.fulfilled, (state, action) => {
        state.loading = false;
        state.isAuthenticated = true;
        state.user = action.payload.user;
        state.error = null;
        state.isInitialized = true;
      })
      .addCase(register.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload;
        state.isAuthenticated = false;
        state.user = null;
        state.isInitialized = true;
      })
      // Refresh cases
      .addCase(refresh.fulfilled, (state, action) => {
        state.isAuthenticated = true;
        state.user = action.payload.user;
        state.error = null;
      })
      .addCase(refresh.rejected, (state) => {
        state.isAuthenticated = false;
        state.user = null;
      })
      // Check auth cases
      .addCase(checkAuth.pending, (state) => {
        state.loading = true;
      })
      .addCase(checkAuth.fulfilled, (state, action) => {
        state.loading = false;
        state.isAuthenticated = true;
        state.user = action.payload;
        state.error = null;
        state.isInitialized = true;
      })
      .addCase(checkAuth.rejected, (state, action) => {
        state.loading = false;
        state.isAuthenticated = false;
        state.user = null;
        state.error = action.payload;
        state.isInitialized = true;
      })
      // Logout cases
      .addCase(logout.pending, (state) => {
        state.loading = true;
      })
      .addCase(logout.fulfilled, (state) => {
        return { ...initialState, isInitialized: true };
      })
      .addCase(logout.rejected, (state) => {
        return { ...initialState, isInitialized: true };
      });
  },
});

export const { clearError, resetState } = authSlice.actions;
export default authSlice.reducer;