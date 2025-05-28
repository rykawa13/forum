import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import { fetchMessagesApi } from '../api/chat';
import { sendMessage as sendWebSocketMessage, isSocketConnected } from '../utils/socket';

export const fetchMessages = createAsyncThunk(
  'chat/fetchMessages',
  async (_, { rejectWithValue }) => {
    try {
      const response = await fetchMessagesApi();
      return response.data;
    } catch (error) {
      return rejectWithValue(error.response?.data || error.message);
    }
  }
);

export const sendMessage = createAsyncThunk(
  'chat/sendMessage',
  async (message, { rejectWithValue, getState }) => {
    try {
      if (!isSocketConnected()) {
        throw new Error('Chat connection is not available. Please try again.');
      }
      
      const state = getState();
      const { user } = state.auth;
      
      const tempId = Date.now().toString();
      const newMessage = {
        id: tempId,
        content: message,
        user_id: user.id,
        username: user.username,
        created_at: new Date().toISOString(),
        type: 'message',
        tempId: tempId
      };

      sendWebSocketMessage({
        type: 'message',
        content: message,
        tempId: tempId
      });
      
      return newMessage;
    } catch (error) {
      return rejectWithValue(error.message);
    }
  }
);

const chatSlice = createSlice({
  name: 'chat',
  initialState: {
    messages: [],
    loading: false,
    error: null,
    connected: false,
    connecting: false
  },
  reducers: {
    receiveMessage: (state, action) => {
      console.log('Processing received message:', action.payload);
      if (action.payload.type === 'message') {
        // Находим сообщение по tempId
        const existingIndex = state.messages.findIndex(
          msg => msg.tempId === action.payload.tempId
        );

        if (existingIndex !== -1) {
          // Обновляем существующее сообщение
          state.messages[existingIndex] = action.payload;
        } else {
          // Добавляем новое сообщение
          state.messages.push(action.payload);
        }
      }
    },
    clearError: (state) => {
      state.error = null;
    },
    setConnected: (state, action) => {
      state.connected = action.payload;
      if (action.payload) {
        state.connecting = false;
      }
    },
    setConnecting: (state, action) => {
      state.connecting = action.payload;
    }
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
      })
      .addCase(sendMessage.pending, (state) => {
        state.error = null;
      })
      .addCase(sendMessage.fulfilled, (state, action) => {
        state.messages = [...state.messages, action.payload];
      })
      .addCase(sendMessage.rejected, (state, action) => {
        state.error = action.payload;
      });
  },
});

export const { receiveMessage, clearError, setConnected, setConnecting } = chatSlice.actions;
export default chatSlice.reducer;