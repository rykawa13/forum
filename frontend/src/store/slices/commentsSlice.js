import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import axios from 'axios';
import { enrichWithUserInfo } from '../../api/users';

const API_URL = process.env.REACT_APP_FORUM_URL || 'http://localhost:8082';

export const fetchComments = createAsyncThunk(
  'comments/fetchComments',
  async (postId, { rejectWithValue }) => {
    try {
      const response = await axios.get(`${API_URL}/posts/${postId}/replies`, {
        withCredentials: true
      });
      console.log('Fetched comments:', response.data);
      
      // Проверяем структуру ответа и получаем массив комментариев
      let comments;
      if (Array.isArray(response.data)) {
        comments = response.data;
      } else if (response.data && Array.isArray(response.data.replies)) {
        comments = response.data.replies;
      } else {
        comments = [];
      }

      // Обогащаем комментарии информацией о пользователях
      comments = await enrichWithUserInfo(comments);
      
      return comments;
    } catch (error) {
      console.error('Error fetching comments:', error.response || error);
      return rejectWithValue(
        error.response?.data?.error || 
        error.response?.data?.message || 
        error.message || 
        'Ошибка при загрузке комментариев'
      );
    }
  }
);

export const addComment = createAsyncThunk(
  'comments/addComment',
  async ({ postId, content, userId }, { rejectWithValue, dispatch }) => {
    try {
      const response = await axios.post(`${API_URL}/posts/${postId}/replies`, {
        content,
        author_id: userId,
      }, {
        withCredentials: true,
      });
      console.log('Added comment:', response.data);
      // После успешного добавления комментария обновляем список
      dispatch(fetchComments(postId));
      return response.data;
    } catch (error) {
      console.error('Error adding comment:', error.response || error);
      return rejectWithValue(
        error.response?.data?.error || 
        error.response?.data?.message || 
        error.message || 
        'Ошибка при добавлении комментария'
      );
    }
  }
);

const commentsSlice = createSlice({
  name: 'comments',
  initialState: {
    items: {},
    loading: false,
    error: null,
  },
  reducers: {
    clearComments: (state) => {
      state.items = {};
      state.error = null;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(fetchComments.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchComments.fulfilled, (state, action) => {
        state.loading = false;
        const postId = action.meta.arg;
        // Теперь action.payload всегда будет массивом
        state.items[postId] = action.payload;
        state.error = null;
      })
      .addCase(fetchComments.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload || 'Ошибка при загрузке комментариев';
      })
      .addCase(addComment.pending, (state) => {
        state.error = null;
      })
      .addCase(addComment.fulfilled, (state, action) => {
        state.error = null;
      })
      .addCase(addComment.rejected, (state, action) => {
        state.error = action.payload || 'Ошибка при добавлении комментария';
      });
  },
});

export const { clearComments } = commentsSlice.actions;
export default commentsSlice.reducer; 