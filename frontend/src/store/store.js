import { configureStore } from '@reduxjs/toolkit';
import authReducer from './authSlice';
import chatReducer from './chatSlice';
import forumReducer from './forumSlice';
import adminReducer from './adminSlice'; 

export const store = configureStore({
  reducer: {
    auth: authReducer,
    chat: chatReducer,
    forum: forumReducer,
    admin: adminReducer,
  },
});