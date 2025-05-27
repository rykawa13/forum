// src/store/adminSlice.js
import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import { fetchUsers as fetchUsersApi, deleteUser as deleteUserApi } from '../api/admin';// Предполагается, что у вас есть API-функции

// Асинхронные действия (thunks)
export const fetchUsers = createAsyncThunk(
  'admin/fetchUsers',
  async (_, { rejectWithValue }) => {
    try {
      const response = await fetchUsersApi();
      return response.data;
    } catch (error) {
      return rejectWithValue(error.response.data);
    }
  }
);

export const deleteUser = createAsyncThunk(
  'admin/deleteUser',
  async (userId, { rejectWithValue }) => {
    try {
      await deleteUserApi(userId);
      return userId; // Возвращаем ID для удаления из состояния
    } catch (error) {
      return rejectWithValue(error.response.data);
    }
  }
);

const adminSlice = createSlice({
  name: 'admin',
  initialState: {
    users: [],
    loading: false,
    error: null,
  },
  reducers: {},
  extraReducers: (builder) => {
    builder
      // Загрузка пользователей
      .addCase(fetchUsers.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchUsers.fulfilled, (state, action) => {
        state.loading = false;
        state.users = action.payload;
      })
      .addCase(fetchUsers.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload;
      })
      // Удаление пользователя
      .addCase(deleteUser.fulfilled, (state, action) => {
        state.users = state.users.filter(user => user.id !== action.payload);
      })
      .addCase(deleteUser.rejected, (state, action) => {
        state.error = action.payload;
      });
  },
});

export default adminSlice.reducer;