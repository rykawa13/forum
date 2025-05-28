import React, { useEffect, useCallback } from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { connectSocket, disconnectSocket } from '../utils/socket';
import { fetchMessages, sendMessage, receiveMessage } from '../store/chatSlice';
import ChatWindow from '../components/chat/ChatWindow';
import { Box, Typography } from '@mui/material';

const ChatPage = () => {
  const dispatch = useDispatch();
  const { isAuthenticated } = useSelector(state => state.auth);
  const { messages } = useSelector(state => state.chat);

  const handleNewMessage = useCallback((message) => {
    console.log('Handling new message in ChatPage:', message);
    dispatch(receiveMessage(message));
  }, [dispatch]);

  useEffect(() => {
    console.log('ChatPage: Initializing chat connection');
    const token = localStorage.getItem('token');
    
    if (!token) {
      console.error('No authentication token found');
      return;
    }
    
    console.log('ChatPage: Got token, connecting to WebSocket');
    connectSocket(token, handleNewMessage);

    return () => {
      console.log('ChatPage: Cleaning up chat connection');
      disconnectSocket();
    };
  }, [handleNewMessage]);

  const handleSendMessage = (messageText) => {
    if (messageText.trim()) {
      console.log('ChatPage: Sending message:', messageText);
      dispatch(sendMessage(messageText));
    }
  };

  return (
    <Box sx={{ p: 3 }}>
      <Typography variant="h4" gutterBottom>
        Чат форума
      </Typography>
      <ChatWindow messages={messages} onSendMessage={handleSendMessage} />
    </Box>
  );
};

export default ChatPage;