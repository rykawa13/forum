import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import { fetchPosts as apiFetchPosts, createPost as apiCreatePost } from '../api/posts';

export const fetchPosts = createAsyncThunk(
  'forum/fetchPosts',
  async (_, { rejectWithValue }) => {
    try {
      const response = await apiFetchPosts(); // Используем переименованную функцию
      return response.data;
    } catch (error) {
      return rejectWithValue(error.response.data);
    }
  }
);

export const createNewPost = createAsyncThunk(
  'forum/createPost',
  async (postData, { rejectWithValue }) => {
    try {
      const response = await apiCreatePost(postData); // Используем переименованную функцию
      return response.data;
    } catch (error) {
      return rejectWithValue(error.response.data);
    }
  }
);

const forumSlice = createSlice({
  name: 'forum',
  initialState: {
    posts: [],
    loading: false,
    error: null,
  },
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(fetchPosts.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchPosts.fulfilled, (state, action) => {
        state.loading = false;
        state.posts = action.payload;
      })
      .addCase(fetchPosts.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload;
      })
      .addCase(createNewPost.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(createNewPost.fulfilled, (state, action) => {
        state.loading = false;
        state.posts.unshift(action.payload);
      })
      .addCase(createNewPost.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload;
      });
  },
});

export default forumSlice.reducer;