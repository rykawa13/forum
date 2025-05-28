import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import { fetchPosts as apiFetchPosts, createPost as apiCreatePost } from '../api/posts';

export const fetchPosts = createAsyncThunk(
  'forum/fetchPosts',
  async (_, { rejectWithValue }) => {
    try {
      const response = await apiFetchPosts();
      console.log('API Response in thunk:', response);

      // Проверяем структуру данных
      const data = response.data;
      if (!data || !Array.isArray(data.posts)) {
        console.error('Invalid response format:', data);
        throw new Error('Invalid response format');
      }

      console.log('Thunk returning:', {
        posts: data.posts,
        total: data.total
      });

      // Возвращаем данные в нужном формате
      return {
        posts: data.posts,
        total: data.total
      };
    } catch (error) {
      console.error('Error in fetchPosts thunk:', error);
      return rejectWithValue(
        error.response?.data?.error || 
        error.response?.data?.message || 
        error.message || 
        'Ошибка при загрузке постов'
      );
    }
  }
);

export const createNewPost = createAsyncThunk(
  'forum/createPost',
  async (postData, { rejectWithValue, dispatch }) => {
    try {
      const response = await apiCreatePost(postData);
      if (!response.data) {
        throw new Error('No data received from server');
      }
      // После успешного создания поста, обновляем список постов
      dispatch(fetchPosts());
      return response.data;
    } catch (error) {
      console.error('Error in createNewPost thunk:', error);
      return rejectWithValue(
        error.response?.data?.error || 
        error.response?.data?.message || 
        error.message || 
        'Ошибка при создании поста'
      );
    }
  }
);

const initialState = {
  posts: [],
  total: 0,
  loading: false,
  error: null,
};

const forumSlice = createSlice({
  name: 'forum',
  initialState,
  reducers: {
    clearError: (state) => {
      state.error = null;
    }
  },
  extraReducers: (builder) => {
    builder
      .addCase(fetchPosts.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchPosts.fulfilled, (state, action) => {
        console.log('Reducer received payload:', action.payload);
        
        if (!action.payload || !Array.isArray(action.payload.posts)) {
          console.error('Invalid payload in reducer:', action.payload);
          return;
        }

        state.posts = action.payload.posts;
        state.total = action.payload.total;
        state.loading = false;
        state.error = null;

        console.log('State updated:', {
          postsCount: state.posts.length,
          total: state.total,
          firstPost: state.posts[0]
        });
      })
      .addCase(fetchPosts.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload || 'Ошибка при загрузке постов';
        console.error('Posts fetch failed:', action.payload);
      })
      .addCase(createNewPost.pending, (state) => {
        state.error = null;
      })
      .addCase(createNewPost.fulfilled, (state) => {
        state.error = null;
      })
      .addCase(createNewPost.rejected, (state, action) => {
        state.error = action.payload || 'Ошибка при создании поста';
        console.error('Post creation failed:', action.payload);
      });
  },
});

export const { clearError } = forumSlice.actions;
export default forumSlice.reducer;