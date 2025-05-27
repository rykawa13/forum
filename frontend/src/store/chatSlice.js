import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import { fetchMessagesApi, sendMessageApi } from '../api/chat';
import { getSocket } from '../utils/socket';

export const fetchMessages = createAsyncThunk(
  'chat/fetchMessages',
  async (_, { rejectWithValue }) => {
    try {
      const response = await fetchMessagesApi();
      return response.data;
    } catch (error) {
      return rejectWithValue(error.response.data);
    }
  }
);

export const sendMessage = createAsyncThunk(
  'chat/sendMessage',
  async (message, { rejectWithValue }) => {
    try {
      const socket = getSocket();
      socket.emit('newMessage', message);
      return message;
    } catch (error) {
      return rejectWithValue(error.response.data);
    }
  }
);

const chatSlice = createSlice({
  name: 'chat',
  initialState: {
    messages: [],
    onlineUsers: 0,
    loading: false,
    error: null,
  },
  reducers: {
    receiveMessage: (state, action) => {
      state.messages.push(action.payload);
    },
    setOnlineUsers: (state, action) => {
      state.onlineUsers = action.payload;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(fetchMessages.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchMessages.fulfilled, (state, action) => {
        state.loading = false;
        state.messages = action.payload;
      })
      .addCase(fetchMessages.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload;
      });
  },
});

export const { receiveMessage, setOnlineUsers } = chatSlice.actions;
export default chatSlice.reducer;